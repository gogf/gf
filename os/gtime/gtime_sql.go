package gtime

import (
	"database/sql/driver"
)

//add Scanner
func (t *Time) Scan(value interface{}) error {
	newTime := New(value)
	t.Time = newTime.Time
	return nil
}

//add valuer
func (t *Time) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return t.Time, nil
}
