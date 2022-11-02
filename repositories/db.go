package repositories

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"edufund/model"
	"edufund/util"

	"github.com/jinzhu/gorm"
)

type DBHandler interface {
	Close() error
	DB() *sql.DB
	New() DBHandler
	NewScope(value interface{}) *gorm.Scope
	CommonDB() gorm.SQLCommon
	Callback() *gorm.Callback
	SetLogger(l gorm.Logger)
	SetQueryLogger(ql *QueryLogger)
	LogMode(enable bool) DBHandler
	SingularTable(enable bool)
	Where(query interface{}, args ...interface{}) DBHandler
	Or(query interface{}, args ...interface{}) DBHandler
	Limit(value int) DBHandler
	Offset(value int) DBHandler
	Order(value string, reorder ...bool) DBHandler
	Select(query interface{}, args ...interface{}) DBHandler
	Omit(columns ...string) DBHandler
	Group(query string) DBHandler
	Having(query string, values ...interface{}) DBHandler
	Joins(query string, args ...interface{}) DBHandler
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) DBHandler
	Unscoped() DBHandler
	Attrs(attrs ...interface{}) DBHandler
	Assign(attrs ...interface{}) DBHandler
	First(out interface{}, where ...interface{}) DBHandler
	Last(out interface{}, where ...interface{}) DBHandler
	Find(out interface{}, where ...interface{}) DBHandler
	Scan(dest interface{}) DBHandler
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	ScanRows(rows *sql.Rows, result interface{}) error
	Pluck(column string, value interface{}) DBHandler
	Count(value interface{}) DBHandler
	Related(value interface{}, foreignKeys ...string) DBHandler
	FirstOrInit(out interface{}, where ...interface{}) DBHandler
	FirstOrCreate(out interface{}, where ...interface{}) DBHandler
	Update(attrs ...interface{}) DBHandler
	Updates(values interface{}, ignoreProtectedAttrs ...bool) DBHandler
	UpdateColumn(attrs ...interface{}) DBHandler
	UpdateColumns(values interface{}) DBHandler
	SaveDefault(value interface{}) DBHandler
	Save(value interface{}) DBHandler
	CreateDefault(value interface{}) DBHandler
	Create(value interface{}) DBHandler
	Delete(value interface{}, where ...interface{}) DBHandler
	Raw(sql string, values ...interface{}) DBHandler
	Exec(sql string, values ...interface{}) DBHandler
	Model(value interface{}) DBHandler
	Table(name string) DBHandler
	Debug() DBHandler
	Begin() DBHandler
	Commit() DBHandler
	Rollback() DBHandler
	NewRecord(value interface{}) bool
	RecordNotFound() bool
	CreateTable(values ...interface{}) DBHandler
	DropTable(values ...interface{}) DBHandler
	DropTableIfExists(values ...interface{}) DBHandler
	HasTable(value interface{}) bool
	AutoMigrate(values ...interface{}) DBHandler
	ModifyColumn(column string, typ string) DBHandler
	DropColumn(column string) DBHandler
	AddIndex(indexName string, column ...string) DBHandler
	AddUniqueIndex(indexName string, column ...string) DBHandler
	RemoveIndex(indexName string) DBHandler
	AddForeignKey(field string, dest string, onDelete string, onUpdate string) DBHandler
	Association(column string) *gorm.Association
	Preload(column string, conditions ...interface{}) DBHandler
	Set(name string, value interface{}) DBHandler
	InstantSet(name string, value interface{}) DBHandler
	Get(name string) (value interface{}, ok bool)
	SetJoinTableHandler(source interface{}, column string, handler gorm.JoinTableHandlerInterface)
	AddError(err error) error
	GetErrors() (errors []error)
	Error() error
	RowsAffected() int64
	RowsToJSON() ([]byte, error)
	RowsToJSONArray() ([][]byte, error)
	SaveJSONByTable(table string, dataJSON []byte) error
}

type db struct {
	gormDB      *gorm.DB
	queryLogger *QueryLogger
}

func OpenDB(dialect string, args ...interface{}) (db DBHandler, err error) {
	gormDB, err := gorm.Open(dialect, args...)
	return Wrap(gormDB, nil), err
}

func Wrap(gormDB *gorm.DB, queryLogger *QueryLogger) DBHandler {
	return &db{gormDB, queryLogger}
}

func (it *db) Close() error {
	return it.gormDB.Close()
}

func (it *db) DB() *sql.DB {
	return it.gormDB.DB()
}

func (it *db) New() DBHandler {
	return Wrap(it.gormDB.New(), it.queryLogger)
}

func (it *db) NewScope(value interface{}) *gorm.Scope {
	return it.gormDB.NewScope(value)
}

func (it *db) CommonDB() gorm.SQLCommon {
	return it.gormDB.CommonDB()
}

func (it *db) Callback() *gorm.Callback {
	return it.gormDB.Callback()
}

func (it *db) SetLogger(log gorm.Logger) {
	it.gormDB.SetLogger(log)
}

func (it *db) SetQueryLogger(log *QueryLogger) {
	it.queryLogger = log
	it.gormDB.SetLogger(log)
}

func (it *db) LogMode(enable bool) DBHandler {
	return Wrap(it.gormDB.LogMode(enable), it.queryLogger)
}

func (it *db) SingularTable(enable bool) {
	it.gormDB.SingularTable(enable)
}

func (it *db) Where(query interface{}, args ...interface{}) DBHandler {
	return Wrap(it.gormDB.Where(query, args...), it.queryLogger)
}

func (it *db) Or(query interface{}, args ...interface{}) DBHandler {
	return Wrap(it.gormDB.Or(query, args...), it.queryLogger)
}

func (it *db) Not(query interface{}, args ...interface{}) DBHandler {
	return Wrap(it.gormDB.Not(query, args...), it.queryLogger)
}

func (it *db) Limit(value int) DBHandler {
	return Wrap(it.gormDB.Limit(value), it.queryLogger)
}

func (it *db) Offset(value int) DBHandler {
	return Wrap(it.gormDB.Offset(value), it.queryLogger)
}

func (it *db) Order(value string, reorder ...bool) DBHandler {
	return Wrap(it.gormDB.Order(value, reorder...), it.queryLogger)
}

func (it *db) Select(query interface{}, args ...interface{}) DBHandler {
	return Wrap(it.gormDB.Select(query, args...), it.queryLogger)
}

func (it *db) Omit(columns ...string) DBHandler {
	return Wrap(it.gormDB.Omit(columns...), it.queryLogger)
}

func (it *db) Group(query string) DBHandler {
	return Wrap(it.gormDB.Group(query), it.queryLogger)
}

func (it *db) Having(query string, values ...interface{}) DBHandler {
	return Wrap(it.gormDB.Having(query, values...), it.queryLogger)
}

func (it *db) Joins(query string, args ...interface{}) DBHandler {
	return Wrap(it.gormDB.Joins(query, args...), it.queryLogger)
}

func (it *db) Scopes(funcs ...func(*gorm.DB) *gorm.DB) DBHandler {
	return Wrap(it.gormDB.Scopes(funcs...), it.queryLogger)
}

func (it *db) Unscoped() DBHandler {
	return Wrap(it.gormDB.Unscoped(), it.queryLogger)
}

func (it *db) Attrs(attrs ...interface{}) DBHandler {
	return Wrap(it.gormDB.Attrs(attrs...), it.queryLogger)
}

func (it *db) Assign(attrs ...interface{}) DBHandler {
	return Wrap(it.gormDB.Assign(attrs...), it.queryLogger)
}

func (it *db) First(out interface{}, where ...interface{}) DBHandler {
	return Wrap(it.gormDB.First(out, where...), it.queryLogger)
}

func (it *db) Last(out interface{}, where ...interface{}) DBHandler {
	return Wrap(it.gormDB.Last(out, where...), it.queryLogger)
}

func (it *db) Find(out interface{}, where ...interface{}) DBHandler {
	return Wrap(it.gormDB.Find(out, where...), it.queryLogger)
}

func (it *db) Scan(dest interface{}) DBHandler {
	return Wrap(it.gormDB.Scan(dest), it.queryLogger)
}

func (it *db) Row() *sql.Row {
	return it.gormDB.Row()
}

func (it *db) Rows() (*sql.Rows, error) {
	return it.gormDB.Rows()
}

func (it *db) ScanRows(rows *sql.Rows, result interface{}) error {
	return it.gormDB.ScanRows(rows, result)
}

func (it *db) Pluck(column string, value interface{}) DBHandler {
	return Wrap(it.gormDB.Pluck(column, value), it.queryLogger)
}

func (it *db) Count(value interface{}) DBHandler {
	return Wrap(it.gormDB.Count(value), it.queryLogger)
}

func (it *db) Related(value interface{}, foreignKeys ...string) DBHandler {
	return Wrap(it.gormDB.Related(value, foreignKeys...), it.queryLogger)
}

func (it *db) FirstOrInit(out interface{}, where ...interface{}) DBHandler {
	return Wrap(it.gormDB.FirstOrInit(out, where...), it.queryLogger)
}

func (it *db) FirstOrCreate(out interface{}, where ...interface{}) DBHandler {
	return Wrap(it.gormDB.FirstOrCreate(out, where...), it.queryLogger)
}

func (it *db) Update(attrs ...interface{}) DBHandler {
	return Wrap(it.gormDB.Update(attrs...), it.queryLogger)
}

func (it *db) Updates(values interface{}, ignoreProtectedAttrs ...bool) DBHandler {
	return Wrap(it.gormDB.Updates(values, ignoreProtectedAttrs...), it.queryLogger)
}

func (it *db) UpdateColumn(attrs ...interface{}) DBHandler {
	return Wrap(it.gormDB.UpdateColumn(attrs...), it.queryLogger)
}

func (it *db) UpdateColumns(values interface{}) DBHandler {
	return Wrap(it.gormDB.UpdateColumns(values), it.queryLogger)
}

func (it *db) SaveDefault(value interface{}) DBHandler {
	return Wrap(it.gormDB.Save(value), it.queryLogger)
}

func (it *db) Save(value interface{}) DBHandler {
	return Wrap(it.gormDB.Save(value), it.queryLogger)
}

func (it *db) CreateDefault(value interface{}) DBHandler {
	return Wrap(it.gormDB.Create(value), it.queryLogger)
}

func (it *db) Create(value interface{}) DBHandler {
	return Wrap(it.gormDB.Create(value), it.queryLogger)
}

func (it *db) Delete(value interface{}, where ...interface{}) DBHandler {
	return Wrap(it.gormDB.Delete(value, where...), it.queryLogger)
}

func (it *db) Raw(sql string, values ...interface{}) DBHandler {
	return Wrap(it.gormDB.Raw(sql, values...), it.queryLogger)
}

func (it *db) Exec(sql string, values ...interface{}) DBHandler {
	return Wrap(it.gormDB.Exec(sql, values...), it.queryLogger)
}

func (it *db) Model(value interface{}) DBHandler {
	return Wrap(it.gormDB.Model(value), it.queryLogger)
}

func (it *db) Table(name string) DBHandler {
	return Wrap(it.gormDB.Table(name), it.queryLogger)
}

func (it *db) Debug() DBHandler {
	return Wrap(it.gormDB.Debug(), it.queryLogger)
}

func (it *db) Begin() DBHandler {
	return Wrap(it.gormDB.Begin(), it.queryLogger)
}

func (it *db) Commit() DBHandler {
	return Wrap(it.gormDB.Commit(), it.queryLogger)
}

func (it *db) Rollback() DBHandler {
	return Wrap(it.gormDB.Rollback(), it.queryLogger)
}

func (it *db) NewRecord(value interface{}) bool {
	return it.gormDB.NewRecord(value)
}

func (it *db) RecordNotFound() bool {
	return it.gormDB.RecordNotFound()
}

func (it *db) CreateTable(values ...interface{}) DBHandler {
	return Wrap(it.gormDB.CreateTable(values...), it.queryLogger)
}

func (it *db) DropTable(values ...interface{}) DBHandler {
	return Wrap(it.gormDB.DropTable(values...), it.queryLogger)
}

func (it *db) DropTableIfExists(values ...interface{}) DBHandler {
	return Wrap(it.gormDB.DropTableIfExists(values...), it.queryLogger)
}

func (it *db) HasTable(value interface{}) bool {
	return it.gormDB.HasTable(value)
}

func (it *db) AutoMigrate(values ...interface{}) DBHandler {
	return Wrap(it.gormDB.AutoMigrate(values...), it.queryLogger)
}

func (it *db) ModifyColumn(column string, typ string) DBHandler {
	return Wrap(it.gormDB.ModifyColumn(column, typ), it.queryLogger)
}

func (it *db) DropColumn(column string) DBHandler {
	return Wrap(it.gormDB.DropColumn(column), it.queryLogger)
}

func (it *db) AddIndex(indexName string, columns ...string) DBHandler {
	return Wrap(it.gormDB.AddIndex(indexName, columns...), it.queryLogger)
}

func (it *db) AddUniqueIndex(indexName string, columns ...string) DBHandler {
	return Wrap(it.gormDB.AddUniqueIndex(indexName, columns...), it.queryLogger)
}

func (it *db) RemoveIndex(indexName string) DBHandler {
	return Wrap(it.gormDB.RemoveIndex(indexName), it.queryLogger)
}

func (it *db) Association(column string) *gorm.Association {
	return it.gormDB.Association(column)
}

func (it *db) Preload(column string, conditions ...interface{}) DBHandler {
	return Wrap(it.gormDB.Preload(column, conditions...), it.queryLogger)
}

func (it *db) Set(name string, value interface{}) DBHandler {
	return Wrap(it.gormDB.Set(name, value), it.queryLogger)
}

func (it *db) InstantSet(name string, value interface{}) DBHandler {
	return Wrap(it.gormDB.InstantSet(name, value), it.queryLogger)
}

func (it *db) Get(name string) (interface{}, bool) {
	return it.gormDB.Get(name)
}

func (it *db) SetJoinTableHandler(source interface{}, column string, handler gorm.JoinTableHandlerInterface) {
	it.gormDB.SetJoinTableHandler(source, column, handler)
}

func (it *db) AddForeignKey(field string, dest string, onDelete string, onUpdate string) DBHandler {
	return Wrap(it.gormDB.AddForeignKey(field, dest, onDelete, onUpdate), it.queryLogger)
}

func (it *db) AddError(err error) error {
	return it.gormDB.AddError(err)
}

func (it *db) GetErrors() (errors []error) {
	return it.gormDB.GetErrors()
}

func (it *db) RowsAffected() int64 {
	return it.gormDB.RowsAffected
}

func (it *db) Error() error {
	return it.gormDB.Error
}

func (it *db) RowsToJSON() ([]byte, error) {
	// returns rows *sql.Rows
	rows, err := it.gormDB.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				if util.IsFloat(b) {
					num, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
					if err != nil {
						return nil, err
					}
					v = num
				} else {
					v = string(b)
				}
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return json.Marshal(tableData)
}

func (it *db) RowsToJSONArray() ([][]byte, error) {
	// returns rows *sql.Rows
	rows, err := it.gormDB.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := [][]byte{}
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				if util.IsFloat(b) {
					num, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
					if err != nil {
						return nil, err
					}
					v = num
				} else {
					v = string(b)
				}
			} else {
				v = val
			}
			entry[col] = v
		}
		entryJSON, err := json.Marshal(entry)
		if err != nil {
			return nil, err
		}
		tableData = append(tableData, entryJSON)
	}
	return tableData, nil
}

func saveDataToTable(h DBHandler, data interface{}) error {
	tx := h.Begin()
	if err := tx.Save(data).Error(); err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
		return nil
	}
}

func (it *db) SaveJSONByTable(table string, dataJSON []byte) error {
	switch table {
	case "artemis_infos":
		data := &model.User{}
		err := json.Unmarshal(dataJSON, data)
		if err != nil {
			return err
		}

		h := it.Model(&model.User{})
		return saveDataToTable(h, data)
	}
	return nil
}

func AutoMigrate(db DBHandler) DBHandler {
	db.AutoMigrate(
		&model.User{},
	)
	return db
}
