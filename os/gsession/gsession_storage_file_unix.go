//go:build unix

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file publishes complete session file snapshots on Unix.

package gsession

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

// setSessionFile publishes the timestamp and payload using same-directory replacement.
func (*StorageFile) setSessionFile(ctx context.Context, path string, content []byte) error {
	// Each writer prepares a private file on the same filesystem as the session.
	dir, _ := filepath.Split(path)
	file, err := os.CreateTemp(dir, "gsession-*.tmp")
	if err != nil {
		return gerror.Wrapf(err, `create temporary session file for "%s" failed`, path)
	}
	tempPath := file.Name()
	defer func() {
		closeErr := file.Close()
		if errors.Is(closeErr, os.ErrClosed) {
			closeErr = nil
		}
		removeErr := os.Remove(tempPath)
		if os.IsNotExist(removeErr) {
			removeErr = nil
		}
		if closeErr != nil || removeErr != nil {
			intlog.Errorf(ctx, `cleanup temporary session file "%s" failed: %+v`, tempPath, errors.Join(closeErr, removeErr))
		}
	}()
	// Keep existing session permissions; newly created sessions remain private.
	if info, statErr := os.Stat(path); statErr == nil {
		if err = file.Chmod(info.Mode().Perm()); err != nil {
			return gerror.Wrapf(err, `set temporary session file permissions for "%s" failed`, path)
		}
	} else if !os.IsNotExist(statErr) {
		return gerror.Wrapf(statErr, `stat session file "%s" failed`, path)
	}
	if _, err = file.Write(gbinary.EncodeInt64(gtime.TimestampMilli())); err != nil {
		return gerror.Wrapf(err, `write data failed to file "%s"`, path)
	}
	if _, err = file.Write(content); err != nil {
		return gerror.Wrapf(err, `write data failed to file "%s"`, path)
	}
	if err = file.Close(); err != nil {
		return gerror.Wrapf(err, `close temporary session file "%s" failed`, tempPath)
	}
	// Publish only complete snapshots after closing the temporary file.
	return gfile.Rename(tempPath, path)
}
