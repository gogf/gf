// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package delivery

import (
	"context"
)

type sDeliveryApp struct{}

func NewDeliveryApp() *sDeliveryApp {
	return &sDeliveryApp{}
}

func (s *sDeliveryApp) Create(ctx context.Context) (err error) {
	return
}

func (s *sDeliveryApp) GetList(ctx context.Context) (err error) {
	return
}

func (s *sDeliveryApp) GetOne(ctx context.Context) (err error) {
	return
}

func (s *sDeliveryApp) Delete(ctx context.Context) (err error) {
	return
}

func (s *sDeliveryApp) AA(ctx context.Context) (err error) { return }
