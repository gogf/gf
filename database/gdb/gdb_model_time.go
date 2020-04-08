// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
)

const (
	gSOFT_FIELD_NAME_CREATE = "create_at"
	gSOFT_FIELD_NAME_UPDATE = "update_at"
	gSOFT_FIELD_NAME_DELETE = "delete_at"
)

// getSoftFieldNameCreate checks and returns the field name for record creating time.
// If there's no field name for storing creating time, it returns an empty string.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func (m *Model) getSoftFieldNameCreate() (field string) {
	fieldsMap, _ := m.db.TableFields(m.tables)
	if len(fieldsMap) > 0 {
		field, _ = gutil.MapPossibleItemByKey(
			gconv.Map(fieldsMap), gSOFT_FIELD_NAME_CREATE,
		)
	}
	return
}

// getSoftFieldNameUpdate checks and returns the field name for record updating time.
// If there's no field name for storing updating time, it returns an empty string.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func (m *Model) getSoftFieldNameUpdate() (field string) {
	fieldsMap, _ := m.db.TableFields(m.tables)
	if len(fieldsMap) > 0 {
		field, _ = gutil.MapPossibleItemByKey(
			gconv.Map(fieldsMap), gSOFT_FIELD_NAME_UPDATE,
		)
	}
	return
}

// getSoftFieldNameDelete checks and returns the field name for record deleting time.
// If there's no field name for storing deleting time, it returns an empty string.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func (m *Model) getSoftFieldNameDelete() (field string) {
	fieldsMap, _ := m.db.TableFields(m.tables)
	if len(fieldsMap) > 0 {
		field, _ = gutil.MapPossibleItemByKey(
			gconv.Map(fieldsMap), gSOFT_FIELD_NAME_DELETE,
		)
	}
	return
}
