//go:build !unix

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file preserves session file writes on platforms without atomic rename.

package gsession

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

// setSessionFile retains direct writes where os.Rename is not guaranteed atomic.
func (*StorageFile) setSessionFile(ctx context.Context, path string, content []byte) error {
	file, err := gfile.OpenWithFlagPerm(
		path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm,
	)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			intlog.Errorf(ctx, `close session file "%s" failed: %+v`, path, closeErr)
		}
	}()
	if _, err = file.Write(gbinary.EncodeInt64(gtime.TimestampMilli())); err != nil {
		return gerror.Wrapf(err, `write data failed to file "%s"`, path)
	}
	if _, err = file.Write(content); err != nil {
		return gerror.Wrapf(err, `write data failed to file "%s"`, path)
	}
	return nil
}
