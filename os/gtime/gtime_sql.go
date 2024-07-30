package gtime

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
