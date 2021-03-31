package gtime

import (
	"database/sql/driver"
)

//Scanner is an interface used by Scan.
//Scan value from database
//database/sql
func (t *Time) Scan(value interface{}) error {
	if t == nil {
		return nil
	}
	newTime := New(value)
	t.Time = newTime.Time
	return nil
}

// Valuer is the interface providing the Value method. database/sql/driver
// Value insert into mysql need this function.
func (t *Time) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	if t.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}
