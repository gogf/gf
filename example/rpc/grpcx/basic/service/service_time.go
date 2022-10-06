// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package service

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/net/context"

	"github.com/gogf/gf/example/rpc/grpcx/basic/protobuf"
)

// Time is the service for time.
type Time struct{}

// Now implements the protobuf.TimeServer interface.
func (s *Time) Now(ctx context.Context, r *protobuf.NowReq) (*protobuf.NowRes, error) {
	text := fmt.Sprintf(`%s: %s`, gcmd.GetOpt("node", "default"), gtime.Now().String())
	return &protobuf.NowRes{Time: text}, nil
}
