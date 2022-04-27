// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/crypto/gaes"
	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
)

// StorageFile implements the Session Storage interface with file system.
type StorageFile struct {
	StorageBase
	path          string        // Session file storage folder path.
	ttl           time.Duration // Session TTL.
	cryptoKey     []byte        // Used when enable crypto feature.
	cryptoEnabled bool          // Used when enable crypto feature.
	updatingIdSet *gset.StrSet  // To be batched updated session id set.
}

const (
	DefaultStorageFileCryptoEnabled        = false
	DefaultStorageFileUpdateTTLInterval    = 10 * time.Second
	DefaultStorageFileClearExpiredInterval = time.Hour
)

var (
	DefaultStorageFilePath      = gfile.Temp("gsessions")
	DefaultStorageFileCryptoKey = []byte("Session storage file crypto key!")
)

// NewStorageFile creates and returns a file storage object for session.
func NewStorageFile(path string, ttl time.Duration) *StorageFile {
	var (
		ctx         = context.TODO()
		storagePath = DefaultStorageFilePath
	)
	if path != "" {
		storagePath, _ = gfile.Search(path)
		if storagePath == "" {
			panic(gerror.NewCodef(gcode.CodeInvalidParameter, `"%s" does not exist`, path))
		}
		if !gfile.IsWritable(storagePath) {
			panic(gerror.NewCodef(gcode.CodeInvalidParameter, `"%s" is not writable`, path))
		}
	}
	if storagePath != "" {
		if err := gfile.Mkdir(storagePath); err != nil {
			panic(gerror.Wrapf(err, `Mkdir "%s" failed in PWD "%s"`, path, gfile.Pwd()))
		}
	}
	s := &StorageFile{
		path:          storagePath,
		ttl:           ttl,
		cryptoKey:     DefaultStorageFileCryptoKey,
		cryptoEnabled: DefaultStorageFileCryptoEnabled,
		updatingIdSet: gset.NewStrSet(true),
	}

	gtimer.AddSingleton(ctx, DefaultStorageFileUpdateTTLInterval, s.timelyUpdateSessionTTL)
	gtimer.AddSingleton(ctx, DefaultStorageFileClearExpiredInterval, s.timelyClearExpiredSessionFile)
	return s
}

// timelyUpdateSessionTTL batch updates the TTL for sessions timely.
func (s *StorageFile) timelyUpdateSessionTTL(ctx context.Context) {
	var (
		sessionId string
		err       error
	)
	// Batch updating sessions.
	for {
		if sessionId = s.updatingIdSet.Pop(); sessionId == "" {
			break
		}
		if err = s.updateSessionTTl(context.TODO(), sessionId); err != nil {
			intlog.Errorf(context.TODO(), `%+v`, err)
		}
	}
}

// timelyClearExpiredSessionFile deletes all expired files timely.
func (s *StorageFile) timelyClearExpiredSessionFile(ctx context.Context) {
	files, err := gfile.ScanDirFile(s.path, "*.session", false)
	if err != nil {
		intlog.Errorf(ctx, `%+v`, err)
		return
	}
	for _, file := range files {
		if err = s.checkAndClearSessionFile(ctx, file); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
		}
	}
}

// SetCryptoKey sets the crypto key for session storage.
// The crypto key is used when crypto feature is enabled.
func (s *StorageFile) SetCryptoKey(key []byte) {
	s.cryptoKey = key
}

// SetCryptoEnabled enables/disables the crypto feature for session storage.
func (s *StorageFile) SetCryptoEnabled(enabled bool) {
	s.cryptoEnabled = enabled
}

// sessionFilePath returns the storage file path for given session id.
func (s *StorageFile) sessionFilePath(sessionId string) string {
	return gfile.Join(s.path, sessionId) + ".session"
}

// RemoveAll deletes all key-value pairs from storage.
func (s *StorageFile) RemoveAll(ctx context.Context, sessionId string) error {
	return gfile.Remove(s.sessionFilePath(sessionId))
}

// GetSession returns the session data as *gmap.StrAnyMap for given session id from storage.
//
// The parameter `ttl` specifies the TTL for this session, and it returns nil if the TTL is exceeded.
// The parameter `data` is the current old session data stored in memory,
// and for some storage it might be nil if memory storage is disabled.
//
// This function is called ever when session starts.
func (s *StorageFile) GetSession(ctx context.Context, sessionId string, ttl time.Duration) (sessionData *gmap.StrAnyMap, err error) {
	var (
		path    = s.sessionFilePath(sessionId)
		content = gfile.GetBytes(path)
	)
	// It updates the TTL only if the session file already exists.
	if len(content) > 8 {
		timestampMilli := gbinary.DecodeToInt64(content[:8])
		if timestampMilli+ttl.Nanoseconds()/1e6 < gtime.TimestampMilli() {
			return nil, nil
		}
		content = content[8:]
		// Decrypt with AES.
		if s.cryptoEnabled {
			content, err = gaes.Decrypt(content, DefaultStorageFileCryptoKey)
			if err != nil {
				return nil, err
			}
		}
		var m map[string]interface{}
		if err = json.UnmarshalUseNumber(content, &m); err != nil {
			return nil, err
		}
		if m == nil {
			return nil, nil
		}
		return gmap.NewStrAnyMapFrom(m, true), nil
	}
	return nil, nil
}

// SetSession updates the data map for specified session id.
// This function is called ever after session, which is changed dirty, is closed.
// This copy all session data map from memory to storage.
func (s *StorageFile) SetSession(ctx context.Context, sessionId string, sessionData *gmap.StrAnyMap, ttl time.Duration) error {
	intlog.Printf(ctx, "StorageFile.SetSession: %s, %v, %v", sessionId, sessionData, ttl)
	path := s.sessionFilePath(sessionId)
	content, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}
	// Encrypt with AES.
	if s.cryptoEnabled {
		content, err = gaes.Encrypt(content, DefaultStorageFileCryptoKey)
		if err != nil {
			return err
		}
	}
	file, err := gfile.OpenWithFlagPerm(
		path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm,
	)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(gbinary.EncodeInt64(gtime.TimestampMilli())); err != nil {
		err = gerror.Wrapf(err, `write data failed to file "%s"`, path)
		return err
	}
	if _, err = file.Write(content); err != nil {
		err = gerror.Wrapf(err, `write data failed to file "%s"`, path)
		return err
	}
	return nil
}

// UpdateTTL updates the TTL for specified session id.
// This function is called ever after session, which is not dirty, is closed.
// It just adds the session id to the async handling queue.
func (s *StorageFile) UpdateTTL(ctx context.Context, sessionId string, ttl time.Duration) error {
	intlog.Printf(ctx, "StorageFile.UpdateTTL: %s, %v", sessionId, ttl)
	if ttl >= DefaultStorageFileUpdateTTLInterval {
		s.updatingIdSet.Add(sessionId)
	}
	return nil
}

// updateSessionTTL updates the TTL for specified session id.
func (s *StorageFile) updateSessionTTl(ctx context.Context, sessionId string) error {
	intlog.Printf(ctx, "StorageFile.updateSession: %s", sessionId)
	path := s.sessionFilePath(sessionId)
	file, err := gfile.OpenWithFlag(path, os.O_WRONLY)
	if err != nil {
		return err
	}
	if _, err = file.WriteAt(gbinary.EncodeInt64(gtime.TimestampMilli()), 0); err != nil {
		err = gerror.Wrapf(err, `write data failed to file "%s"`, path)
		return err
	}
	return file.Close()
}

func (s *StorageFile) checkAndClearSessionFile(ctx context.Context, path string) (err error) {
	var (
		file                *os.File
		readBytesCount      int
		timestampMilliBytes = make([]byte, 8)
	)
	file, err = gfile.OpenWithFlag(path, os.O_RDONLY)
	if err != nil {
		return err
	}
	defer file.Close()
	// Read the session file updated timestamp in milliseconds.
	readBytesCount, err = file.Read(timestampMilliBytes)
	if err != nil {
		return
	}
	if readBytesCount != 8 {
		return gerror.Newf(`invalid read bytes count "%d", expect "8"`, readBytesCount)
	}
	// Remove expired session file.
	var (
		ttlInMilliseconds     = s.ttl.Nanoseconds() / 1e6
		fileTimestampMilli    = gbinary.DecodeToInt64(timestampMilliBytes)
		currentTimestampMilli = gtime.TimestampMilli()
	)
	if fileTimestampMilli+ttlInMilliseconds < currentTimestampMilli {
		intlog.PrintFunc(ctx, func() string {
			return fmt.Sprintf(
				`clear expired session file "%s": updated datetime "%s", ttl "%s"`,
				path, gtime.NewFromTimeStamp(fileTimestampMilli), s.ttl,
			)
		})
		return gfile.Remove(path)
	}
	return nil
}
