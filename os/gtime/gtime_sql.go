package gtime

import (
	"database/sql/driver"
)

// Scan implements interface used by Scan in package database/sql for Scanning value
// from database to local golang variable.
func (t *Time) Scan(value interface{}) error {
	if t == nil {
		return nil
	}
	newTime := New(value)
	*t = *newTime
	return nil
}

// Value is the interface providing the Value method for package database/sql/driver
// for retrieving value from golang variable to database.
func (t *Time) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	if t.IsZero() {
		return nil, nil
	}

	if t.Year() == 0 {
		// Only time.
		return t.Format("15:04:05"), nil
	}

	return t.Time, nil
}
