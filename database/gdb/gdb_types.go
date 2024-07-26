package gdb

import "database/sql/driver"

type ScannerForFieldType interface {
	ScanForFieldType(fieldType string, src any) error
}

type ValuerForFieldType interface {
	ValueForFieldType(fieldType string) (driver.Value, error)
}
