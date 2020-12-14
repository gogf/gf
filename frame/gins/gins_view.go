// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"fmt"
	"github.com/gogf/gf/os/gview"
	"github.com/gogf/gf/util/gutil"
)

const (
	gFRAME_CORE_COMPONENT_NAME_VIEWER = "gf.core.component.viewer"
	gVIEWER_NODE_NAME                 = "viewer"
)

// View returns an instance of View with default settings.
// The parameter <name> is the name for the instance.
func View(name ...string) *gview.View {
	instanceName := gview.DefaultName
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_VIEWER, instanceName)
	return instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		return getViewInstance(instanceName)
	}).(*gview.View)
}

func getViewInstance(name ...string) *gview.View {
	instanceName := gview.DefaultName
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	view := gview.Instance(instanceName)
	// To avoid file no found error while it's not necessary.
	if Config().Available() {
		var m map[string]interface{}
		nodeKey, _ := gutil.MapPossibleItemByKey(Config().GetMap("."), gVIEWER_NODE_NAME)
		if nodeKey == "" {
			nodeKey = gVIEWER_NODE_NAME
		}
		m = Config().GetMap(fmt.Sprintf(`%s.%s`, nodeKey, instanceName))
		if len(m) == 0 {
			m = Config().GetMap(nodeKey)
		}
		if len(m) > 0 {
			if err := view.SetConfigWithMap(m); err != nil {
				panic(err)
			}
		}
	}
	return view
}
