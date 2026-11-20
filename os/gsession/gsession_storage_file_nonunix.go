//go:build !unix

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file publishes session file snapshots on non-Unix platforms.
// Where os.Rename cannot be relied on to be atomic (e.g. Windows, where a
// reader holding the target open blocks the rename), a per-storage read/write
// lock serializes writers against readers so that the file is never observed
// half-written: the writer builds the complete payload in memory and writes it
// in a single call under the write lock, eliminating the two-write gap of the
// former O_TRUNC implementation.

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

// setSessionFile writes the timestamp and payload as a single snapshot under
// the storage write lock. It concatenates the 8-byte timestamp and the payload
// in memory and writes them in one call, so a concurrent reader (holding the
// read lock) never sees the file between the two writes.
func (s *StorageFile) setSessionFile(ctx context.Context, path string, content []byte) error {
	buf := make([]byte, 8+len(content))
	copy(buf[:8], gbinary.EncodeInt64(gtime.TimestampMilli()))
	copy(buf[8:], content)

	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := gfile.OpenWithFlagPerm(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			intlog.Errorf(ctx, `close session file "%s" failed: %+v`, path, closeErr)
		}
	}()
	if _, err = file.Write(buf); err != nil {
		return gerror.Wrapf(err, `write data failed to file "%s"`, path)
	}
	return nil
}

// readSessionBytes reads the session file under the read lock so that it never
// overlaps an in-progress write.
func (s *StorageFile) readSessionBytes(path string) []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return gfile.GetBytes(path)
}
