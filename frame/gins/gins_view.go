// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	frameCoreComponentNameViewer = "gf.core.component.viewer"
	configNodeNameViewer         = "viewer"
)

// View returns an instance of View with default settings.
// The parameter `name` is the name for the instance.
// Note that it panics if any error occurs duration instance creating.
func View(name ...string) *gview.View {
	instanceName := gview.DefaultName
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameViewer, instanceName)
	return localInstances.GetOrSetFuncLock(instanceKey, func() interface{} {
		return getViewInstance(instanceName)
	}).(*gview.View)
}

func getViewInstance(name ...string) *gview.View {
	var (
		ctx          = context.Background()
		instanceName = gview.DefaultName
	)
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	view := gview.Instance(instanceName)
	// To avoid file no found error while it's not necessary.
	var (
		configMap      map[string]interface{}
		configNodeName = configNodeNameViewer
	)
	if configData, _ := Config().Data(ctx); len(configData) > 0 {
		if v, _ := gutil.MapPossibleItemByKey(configData, configNodeNameViewer); v != "" {
			configNodeName = v
		}
	}
	if v, _ := Config().Get(ctx, fmt.Sprintf(`%s.%s`, configNodeName, instanceName)); !v.IsEmpty() {
		configMap = v.Map()
	}
	if len(configMap) == 0 {
		if v, _ := Config().Get(ctx, configNodeName); !v.IsEmpty() {
			configMap = v.Map()
		}
	}
	if len(configMap) > 0 {
		if err := view.SetConfigWithMap(configMap); err != nil {
			panic(err)
		}
	}
	return view
}
