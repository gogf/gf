// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package user

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genservice/service"
)

func init() {
	service.RegisterUser(New())
}

type sUser struct {
}

func New() *sUser {
	return &sUser{}
}

// Create creates a new user.
func (s *sUser) Create(ctx context.Context, name string) (id int, err error) {
	return 0, nil
}

// GetOne retrieves user by id.
func (s *sUser) GetOne(ctx context.Context, id int) (name string, err error) {
	return "", nil
}

// GetList retrieves user list.
func (s *sUser) GetList(ctx context.Context) (names []string, err error) {
	return nil, nil
}

// Update updates user by id.
func (s *sUser) Update(ctx context.Context, id int) (name string, err error) {
	return "", nil
}

// Delete deletes user by id.
func (s *sUser) Delete(ctx context.Context, id int) (err error) {
	return nil
}
