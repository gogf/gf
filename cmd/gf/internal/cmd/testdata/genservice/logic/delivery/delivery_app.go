// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package delivery

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genservice/service"
)

type sDeliveryApp struct{}

func NewDeliveryApp() *sDeliveryApp {
	return &sDeliveryApp{}
}

func (s *sDeliveryApp) Create(ctx context.Context) (i service.IDeliveryCluster, err error) {
	return
}

func (s *sDeliveryApp) GetList(ctx context.Context, i service.IDeliveryCluster) (err error) {
	service.Article().Get(ctx, 1)
	return
}

func (s *sDeliveryApp) GetOne(ctx context.Context) (err error) {
	return
}

func (s *sDeliveryApp) Delete(ctx context.Context) (err error) {
	return
}

func (s *sDeliveryApp) AA(ctx context.Context) (err error) { return }
