// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

// DriverWrapper is a driver wrapper for extending features with embedded driver.
type DriverWrapper struct {
	driver Driver
}

// New creates and returns a database object for mysql.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *DriverWrapper) New(core *Core, node *ConfigNode) (DB, error) {
	db, err := d.driver.New(core, node)
	if err != nil {
		return nil, err
	}
	return &DriverWrapperDB{
		DB: db,
	}, nil
}

// newDriverWrapper creates and returns a driver wrapper.
func newDriverWrapper(driver Driver) Driver {
	return &DriverWrapper{
		driver: driver,
	}
}
