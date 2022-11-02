package util

import (
	"net/mail"
	"strings"
)

type StringUtil struct{}

var String StringUtil = StringUtil{}

func (su StringUtil) IsNull(str string) bool {
	return strings.ToLower(str) == "null"
}

func (su StringUtil) HaveValueArray(strArray []string) bool {
	return strArray != nil || len(strArray) > 0
}

func (su StringUtil) IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

func GetExcelColumnName(columnNumber int) string {
	dividend := columnNumber
	var columnName string
	var modulo int

	for dividend > 0 {
		modulo = (dividend - 1) % 26
		columnName = toCharStr(modulo) + columnName
		dividend = int((dividend - modulo) / 26)
	}

	return columnName
}

func toCharStr(i int) string {
	return string(rune('A' + i))
}

func ValidMailAddress(address string) (string, error) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "", err
	}
	return addr.Address, err
}
