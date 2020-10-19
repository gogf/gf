// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/internal/json"
	"os"
	"time"

	"github.com/gogf/gf/crypto/gaes"

	"github.com/gogf/gf/os/gtimer"

	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/encoding/gbinary"

	"github.com/gogf/gf/os/gtime"

	"github.com/gogf/gf/os/gfile"
)

// StorageFile implements the Session Storage interface with file system.
type StorageFile struct {
	path          string
	cryptoKey     []byte
	cryptoEnabled bool
	updatingIdSet *gset.StrSet
}

var (
	DefaultStorageFilePath          = gfile.TempDir("gsessions")
	DefaultStorageFileCryptoKey     = []byte("Session storage file crypto key!")
	DefaultStorageFileCryptoEnabled = false
	DefaultStorageFileLoopInterval  = 10 * time.Second
)

// NewStorageFile creates and returns a file storage object for session.
func NewStorageFile(path ...string) *StorageFile {
	storagePath := DefaultStorageFilePath
	if len(path) > 0 && path[0] != "" {
		storagePath, _ = gfile.Search(path[0])
		if storagePath == "" {
			panic(fmt.Sprintf("'%s' does not exist", path[0]))
		}
		if !gfile.IsWritable(storagePath) {
			panic(fmt.Sprintf("'%s' is not writable", path[0]))
		}
	}
	if storagePath != "" {
		if err := gfile.Mkdir(storagePath); err != nil {
			panic(fmt.Sprintf("mkdir '%s' failed: %v", path[0], err))
		}
	}
	s := &StorageFile{
		path:          storagePath,
		cryptoKey:     DefaultStorageFileCryptoKey,
		cryptoEnabled: DefaultStorageFileCryptoEnabled,
		updatingIdSet: gset.NewStrSet(true),
	}
	// Batch updates the TTL for session ids timely.
	gtimer.AddSingleton(DefaultStorageFileLoopInterval, func() {
		//intlog.Print("StorageFile.timer start")
		var id string
		var err error
		for {
			if id = s.updatingIdSet.Pop(); id == "" {
				break
			}
			if err = s.doUpdateTTL(id); err != nil {
				intlog.Error(err)
			}
		}
		//intlog.Print("StorageFile.timer end")
	})
	return s
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
func (s *StorageFile) sessionFilePath(id string) string {
	return gfile.Join(s.path, id)
}

// New creates a session id.
// This function can be used for custom session creation.
func (s *StorageFile) New(ttl time.Duration) (id string) {
	return ""
}

// Get retrieves session value with given key.
// It returns nil if the key does not exist in the session.
func (s *StorageFile) Get(id string, key string) interface{} {
	return nil
}

// GetMap retrieves all key-value pairs as map from storage.
func (s *StorageFile) GetMap(id string) map[string]interface{} {
	return nil
}

// GetSize retrieves the size of key-value pairs from storage.
func (s *StorageFile) GetSize(id string) int {
	return -1
}

// Set sets key-value session pair to the storage.
// The parameter <ttl> specifies the TTL for the session id (not for the key-value pair).
func (s *StorageFile) Set(id string, key string, value interface{}, ttl time.Duration) error {
	return ErrorDisabled
}

// SetMap batch sets key-value session pairs with map to the storage.
// The parameter <ttl> specifies the TTL for the session id(not for the key-value pair).
func (s *StorageFile) SetMap(id string, data map[string]interface{}, ttl time.Duration) error {
	return ErrorDisabled
}

// Remove deletes key with its value from storage.
func (s *StorageFile) Remove(id string, key string) error {
	return ErrorDisabled
}

// RemoveAll deletes all key-value pairs from storage.
func (s *StorageFile) RemoveAll(id string) error {
	return ErrorDisabled
}

// GetSession returns the session data as *gmap.StrAnyMap for given session id from storage.
//
// The parameter <ttl> specifies the TTL for this session, and it returns nil if the TTL is exceeded.
// The parameter <data> is the current old session data stored in memory,
// and for some storage it might be nil if memory storage is disabled.
//
// This function is called ever when session starts.
func (s *StorageFile) GetSession(id string, ttl time.Duration, data *gmap.StrAnyMap) (*gmap.StrAnyMap, error) {
	if data != nil {
		return data, nil
	}
	//intlog.Printf("StorageFile.GetSession: %s, %v", id, ttl)
	path := s.sessionFilePath(id)
	content := gfile.GetBytes(path)
	if len(content) > 8 {
		timestampMilli := gbinary.DecodeToInt64(content[:8])
		if timestampMilli+ttl.Nanoseconds()/1e6 < gtime.TimestampMilli() {
			return nil, nil
		}
		var err error
		content = content[8:]
		// Decrypt with AES.
		if s.cryptoEnabled {
			content, err = gaes.Decrypt(content, DefaultStorageFileCryptoKey)
			if err != nil {
				return nil, err
			}
		}
		var m map[string]interface{}
		if err = json.Unmarshal(content, &m); err != nil {
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
func (s *StorageFile) SetSession(id string, data *gmap.StrAnyMap, ttl time.Duration) error {
	intlog.Printf("StorageFile.SetSession: %s, %v, %v", id, data, ttl)
	path := s.sessionFilePath(id)
	content, err := json.Marshal(data)
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
		return err
	}
	if _, err = file.Write(content); err != nil {
		return err
	}
	return nil
}

// UpdateTTL updates the TTL for specified session id.
// This function is called ever after session, which is not dirty, is closed.
// It just adds the session id to the async handling queue.
func (s *StorageFile) UpdateTTL(id string, ttl time.Duration) error {
	intlog.Printf("StorageFile.UpdateTTL: %s, %v", id, ttl)
	if ttl >= DefaultStorageRedisLoopInterval {
		s.updatingIdSet.Add(id)
	}
	return nil
}

// doUpdateTTL updates the TTL for session id.
func (s *StorageFile) doUpdateTTL(id string) error {
	intlog.Printf("StorageFile.doUpdateTTL: %s", id)
	path := s.sessionFilePath(id)
	file, err := gfile.OpenWithFlag(path, os.O_WRONLY)
	if err != nil {
		return err
	}
	if _, err = file.WriteAt(gbinary.EncodeInt64(gtime.TimestampMilli()), 0); err != nil {
		return err
	}
	return file.Close()
}
