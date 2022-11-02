package util

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

type LocalTime struct {
	time.Time
}

func (localTime *LocalTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		t2, err2 := time.ParseInLocation(time.RFC3339, s, time.Local)
		if err2 != nil {
			return errors.New(fmt.Sprint("Invalid Time Format (expected in YYYY-MM-DD format or RFC3339), ", err.Error(), ", ", err2.Error()))
		}
		t = t2
	}
	localTime.Time = t
	return nil
}

func (localTime *LocalTime) Value() (driver.Value, error) {
	if localTime == nil {
		return nil, nil
	}
	return localTime.Time, nil
}

func (localTime *LocalTime) Scan(src interface{}) (err error) {
	var ok bool
	localTime.Time, ok = src.(time.Time)
	if !ok {
		return errors.New("incompatible type src to time.Time")
	}

	return nil
}

func (localTime *LocalTime) ScanString(src string) (err error) {
	localTime.Time, err = time.Parse(time.RFC3339, src)
	if err != nil {
		return err
	}
	return nil
}
