// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"fmt"
	"github.com/gogf/gf/os/gview"
)

const (
	gFRAME_CORE_COMPONENT_NAME_VIEWER = "gf.core.component.viewer"
	gVIEWER_NODE_NAME                 = "viewer"
)

// View returns an instance of View with default settings.
// The parameter <name> is the name for the instance.
func View(name ...string) *gview.View {
	instanceName := gview.DEFAULT_NAME
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_VIEWER, instanceName)
	return instances.GetOrSetFunc(instanceKey, func() interface{} {
		view := gview.Instance(instanceName)
		// To avoid file no found error while it's not necessary.
		if Config().Available() {
			var m map[string]interface{}
			// It firstly searches the configuration of the instance name.
			if m = Config().GetMap(fmt.Sprintf(`%s.%s`, gVIEWER_NODE_NAME, instanceName)); m == nil {
				// If the configuration for the instance does not exist,
				// it uses the default view configuration.
				m = Config().GetMap(gVIEWER_NODE_NAME)
			}
			if m != nil {
				if err := view.SetConfigWithMap(m); err != nil {
					panic(err)
				}
			}
		}
		return view
	}).(*gview.View)
}
