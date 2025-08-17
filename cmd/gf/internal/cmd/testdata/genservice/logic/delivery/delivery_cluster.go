// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package delivery

import (
	"context"

	gdbas "github.com/gogf/gf/v2/database/gdb"
)

type sDeliveryCluster struct{}

func NewDeliveryCluster() *sDeliveryCluster {
	return &sDeliveryCluster{}
}

// Create 自动创建Cluster及Project.
func (s *sDeliveryCluster) Create(ctx context.Context) (err error, gdb gdbas.Model) {
	return
}

func (s *sDeliveryCluster) Delete(ctx context.Context) (err error) {
	return
}

func (s *sDeliveryCluster) GetList(ctx context.Context) (err error) {
	return
}
