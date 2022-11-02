package repositories

import (
	"database/sql/driver"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"

	"edufund/util"

	"github.com/go-stomp/stomp"
	"github.com/jinzhu/gorm"
)

type Pending struct {
	IsQuerySuccessCh chan bool
	Query            string
}
type QueryLogger struct {
	DB             DBHandler
	OnSuccess      func(conn *stomp.Conn, rawQuery string) error
	Conn           *stomp.Conn
	PendingSendMap map[string]*Pending
}

func (logger *QueryLogger) Print(v ...interface{}) {
	logType := v[0].(string)
	if !isSql(logType) {
		log.Println("[QueryLogger] Invalid logType, logType :", logType)
		log.Println("[QueryLogger] v :", v)
		return
	}
	query := v[3].(string)
	args := v[4].([]interface{})
	if !isSendable(query) {
		log.Println("[QueryLogger] Not Sendable, query :", query)
		return
	}
	pendingKey, err := getPendingKey(query, args)
	if err != nil {
		log.Println("[QueryLogger] Failed to get pendingKey, err :", err.Error())
		return
	}
	sendableQuery := getSendableQuery(v...)

	pendingSend, ok := logger.PendingSendMap[pendingKey]
	if !ok {
		log.Printf("[QueryLogger] Failed to get pendingSend with pendingKey : %s, possibly due to tx Rollback", pendingKey)
		return
	}

	pendingSend.Query = sendableQuery
	go logger.Send(pendingKey)
}

func (logger *QueryLogger) Println(v ...interface{}) {
	logger.Print(v...)
	log.Println(gorm.LogFormatter(v...)...)
}

func (logger *QueryLogger) Send(pendingKey string) {
	success := <-logger.PendingSendMap[pendingKey].IsQuerySuccessCh
	if success {
		if logger.OnSuccess == nil {
			return
		}

		err := logger.OnSuccess(logger.Conn, logger.PendingSendMap[pendingKey].Query)
		if err != nil {
			log.Println("[QueryLogger] Failed to execute OnSuccess function,", err.Error())
		}
	}

	delete(logger.PendingSendMap, pendingKey)
}

func (ql *QueryLogger) FormatPendingKey(id uint, tableName string) string {
	return fmt.Sprintf("%d - %s", id, tableName)
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func getSendableQuery(values ...interface{}) (validQuery string) {
	if len(values) > 1 {
		var (
			formattedValues []string
			level           = values[0].(string)
		)

		if isSql(level) {
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						if t.IsZero() {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", "0000-00-00 00:00:00"))
						} else {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05.000")))
						}
					} else if t, ok := value.(util.LocalTime); ok {
						if t.IsZero() {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", "0000-00-00 00:00:00"))
						} else {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05.000")))
						}
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						switch value.(type) {
						case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
							formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
						default:
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						}
					}
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			}

			numericPlaceHolderRegexp := regexp.MustCompile(`\$\d+`)
			sqlRegexp := regexp.MustCompile(`\?`)
			// differentiate between $n placeholders or else treat like ?
			if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
				validQuery = values[3].(string)
				for index, value := range formattedValues {
					placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
					validQuery = regexp.MustCompile(placeholder).ReplaceAllString(validQuery, value+"$1")
				}
			} else {
				formattedValuesLength := len(formattedValues)
				for index, value := range sqlRegexp.Split(values[3].(string), -1) {
					validQuery += value
					if index < formattedValuesLength {
						validQuery += formattedValues[index]
					}
				}
			}
		}
	}

	return
}

func isSql(logType string) bool {
	return logType == "sql"
}
func isInsert(query string) bool {
	return strings.HasPrefix(query, "INSERT")
}

func isInsertWithID(query string) bool {
	return isInsert(query) && strings.Contains(query, `("id"`)
}

func isUpdate(query string) bool {
	return strings.HasPrefix(query, "UPDATE")
}

func isSendable(query string) bool {
	return isInsertWithID(query) || isUpdate(query)
}

func getIDFromLogArgs(query string, args []interface{}) (id uint, err error) {
	if isInsertWithID(query) {
		argsID, ok := args[0].(uint)
		if ok {
			id = argsID
		} else {
			err = fmt.Errorf("[QueryLogger] Failed to get ID from INSERT query, fail to convert %v to uint", args[0])
		}
	} else if isUpdate(query) {
		argsID, ok := args[len(args)-1].(uint)
		if ok {
			id = argsID
		} else {
			err = fmt.Errorf("[QueryLogger] Failed to get ID from UPDATE query, fail to convert %v to uint", args[len(args)-1])
		}
	} else {
		err = fmt.Errorf("[QueryLogger] Failed to get ID, Invalid query, query is not INSERT with id nor UPDATE, query : %s", query)
	}

	return
}

func getTableName(query string) (tableName string, err error) {
	if isInsert(query) {
		splitQuery := strings.SplitN(query, " ", 4)
		tableName = strings.TrimSpace(strings.ReplaceAll(splitQuery[2], `"`, ""))
	} else if isUpdate(query) {
		splitQuery := strings.SplitN(query, " ", 3)
		tableName = strings.TrimSpace(strings.ReplaceAll(splitQuery[1], `"`, ""))
	} else {
		err = fmt.Errorf("[QueryLogger] Failed to get TableName, Unexpected query format, query : %s", query)
	}

	return
}

func FormatPendingKey(id uint, tableName string) string {
	return fmt.Sprintf("%d - %s", id, tableName)
}

func getPendingKey(query string, args []interface{}) (pendingKey string, err error) {
	id, err := getIDFromLogArgs(query, args)
	if err != nil {
		return
	}

	tableName, err := getTableName(query)
	if err != nil {
		return
	}

	pendingKey = FormatPendingKey(id, tableName)
	return
}
