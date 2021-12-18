// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
)

// Session struct for storing single session data, which is bound to a single request.
// The Session struct is the interface with user, but the Storage is the underlying adapter designed interface
// for functionality implements.
type Session struct {
	id      string          // Session id.
	ctx     context.Context // Context for current session. Please note that, session live along with context.
	data    *gmap.StrAnyMap // Session data.
	dirty   bool            // Used to mark session is modified.
	start   bool            // Used to mark session is started.
	manager *Manager        // Parent manager.

	// idFunc is a callback function used for creating custom session id.
	// This is called if session id is empty ever when session starts.
	idFunc func(ttl time.Duration) (id string)
}

// init does the lazy initialization for session.
// It here initializes real session if necessary.
func (s *Session) init() error {
	if s.start {
		return nil
	}
	var err error
	if s.id != "" {
		// Retrieve memory session data from manager.
		r, err := s.manager.sessionData.Get(s.ctx, s.id)
		if err != nil && err != ErrorDisabled {
			return err
		}
		if r != nil {
			s.data = r.Val().(*gmap.StrAnyMap)
			intlog.Print(s.ctx, "session init data:", s.data)
		}
		// Retrieve stored session data from storage.
		if s.manager.storage != nil {
			if s.data, err = s.manager.storage.GetSession(s.ctx, s.id, s.manager.ttl, s.data); err != nil && err != ErrorDisabled {
				intlog.Errorf(s.ctx, "session restoring failed for id '%s': %v", s.id, err)
				return err
			}
		}
	}
	// Use custom session id creating function.
	if s.id == "" && s.idFunc != nil {
		s.id = s.idFunc(s.manager.ttl)
	}
	// Use default session id creating function of storage.
	if s.id == "" {
		s.id, err = s.manager.storage.New(s.ctx, s.manager.ttl)
		if err != nil && err != ErrorDisabled {
			intlog.Errorf(s.ctx, "create session id failed: %v", err)
			return err
		}
	}
	// Use default session id creating function.
	if s.id == "" {
		s.id = NewSessionId()
	}
	if s.data == nil {
		s.data = gmap.NewStrAnyMap(true)
	}
	s.start = true
	return nil
}

// Close closes current session and updates its ttl in the session manager.
// If this session is dirty, it also exports it to storage.
//
// NOTE that this function must be called ever after a session request done.
func (s *Session) Close() error {
	if s.start && s.id != "" {
		size := s.data.Size()
		if s.manager.storage != nil {
			if s.dirty {
				if err := s.manager.storage.SetSession(s.ctx, s.id, s.data, s.manager.ttl); err != nil && err != ErrorDisabled {
					return err
				}
			} else if size > 0 {
				if err := s.manager.storage.UpdateTTL(s.ctx, s.id, s.manager.ttl); err != nil && err != ErrorDisabled {
					return err
				}
			}
		}
		if s.dirty || size > 0 {
			s.manager.UpdateSessionTTL(s.id, s.data)
		}
	}
	return nil
}

// Set sets key-value pair to this session.
func (s *Session) Set(key string, value interface{}) error {
	if err := s.init(); err != nil {
		return err
	}
	if err := s.manager.storage.Set(s.ctx, s.id, key, value, s.manager.ttl); err != nil {
		if err == ErrorDisabled {
			s.data.Set(key, value)
		} else {
			return err
		}
	}
	s.dirty = true
	return nil
}

// SetMap batch sets the session using map.
func (s *Session) SetMap(data map[string]interface{}) error {
	if err := s.init(); err != nil {
		return err
	}
	if err := s.manager.storage.SetMap(s.ctx, s.id, data, s.manager.ttl); err != nil {
		if err == ErrorDisabled {
			s.data.Sets(data)
		} else {
			return err
		}
	}
	s.dirty = true
	return nil
}

// Remove removes key along with its value from this session.
func (s *Session) Remove(keys ...string) error {
	if s.id == "" {
		return nil
	}
	if err := s.init(); err != nil {
		return err
	}
	for _, key := range keys {
		if err := s.manager.storage.Remove(s.ctx, s.id, key); err != nil {
			if err == ErrorDisabled {
				s.data.Remove(key)
			} else {
				return err
			}
		}
	}
	s.dirty = true
	return nil
}

// RemoveAll deletes all key-value pairs from this session.
func (s *Session) RemoveAll() error {
	if s.id == "" {
		return nil
	}
	if err := s.init(); err != nil {
		return err
	}
	if err := s.manager.storage.RemoveAll(s.ctx, s.id); err != nil {
		if err == ErrorDisabled {
			s.data.Clear()
		} else {
			return err
		}
	}
	s.dirty = true
	return nil
}

// Id returns the session id for this session.
// It creates and returns a new session id if the session id is not passed in initialization.
func (s *Session) Id() (string, error) {
	if err := s.init(); err != nil {
		return "", err
	}
	return s.id, nil
}

// SetId sets custom session before session starts.
// It returns error if it is called after session starts.
func (s *Session) SetId(id string) error {
	if s.start {
		return gerror.NewCode(gcode.CodeInvalidOperation, "session already started")
	}
	s.id = id
	return nil
}

// SetIdFunc sets custom session id creating function before session starts.
// It returns error if it is called after session starts.
func (s *Session) SetIdFunc(f func(ttl time.Duration) string) error {
	if s.start {
		return gerror.NewCode(gcode.CodeInvalidOperation, "session already started")
	}
	s.idFunc = f
	return nil
}

// Data returns all data as map.
// Note that it's using value copy internally for concurrent-safe purpose.
func (s *Session) Data() (map[string]interface{}, error) {
	if s.id != "" {
		if err := s.init(); err != nil {
			return nil, err
		}
		data, err := s.manager.storage.Data(s.ctx, s.id)
		if err != nil && err != ErrorDisabled {
			intlog.Error(s.ctx, err)
		}
		if data != nil {
			return data, nil
		}
		return s.data.Map(), nil
	}
	return map[string]interface{}{}, nil
}

// Size returns the size of the session.
func (s *Session) Size() (int, error) {
	if s.id != "" {
		if err := s.init(); err != nil {
			return 0, err
		}
		size, err := s.manager.storage.GetSize(s.ctx, s.id)
		if err != nil && err != ErrorDisabled {
			intlog.Error(s.ctx, err)
		}
		if size >= 0 {
			return size, nil
		}
		return s.data.Size(), nil
	}
	return 0, nil
}

// Contains checks whether key exist in the session.
func (s *Session) Contains(key string) (bool, error) {
	if err := s.init(); err != nil {
		return false, err
	}
	v, err := s.Get(key)
	if err != nil {
		return false, err
	}
	return !v.IsNil(), nil
}

// IsDirty checks whether there's any data changes in the session.
func (s *Session) IsDirty() bool {
	return s.dirty
}

// Get retrieves session value with given key.
// It returns `def` if the key does not exist in the session if `def` is given,
// or else it returns nil.
func (s *Session) Get(key string, def ...interface{}) (*gvar.Var, error) {
	if s.id == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, `session id cannot be empty`)
	}
	if err := s.init(); err != nil {
		return nil, err
	}
	v, err := s.manager.storage.Get(s.ctx, s.id, key)
	if err != nil && err != ErrorDisabled {
		intlog.Error(s.ctx, err)
		return nil, err
	}
	if v != nil {
		return gvar.New(v), nil
	}
	if v := s.data.Get(key); v != nil {
		return gvar.New(v), nil
	}
	if len(def) > 0 {
		return gvar.New(def[0]), nil
	}
	return nil, nil
}

// MustId performs as function Id, but it panics if any error occurs.
func (s *Session) MustId() string {
	id, err := s.Id()
	if err != nil {
		panic(err)
	}
	return id
}

// MustGet performs as function Get, but it panics if any error occurs.
func (s *Session) MustGet(key string, def ...interface{}) *gvar.Var {
	v, err := s.Get(key, def...)
	if err != nil {
		panic(err)
	}
	return v
}

// MustSet performs as function Set, but it panics if any error occurs.
func (s *Session) MustSet(key string, value interface{}) {
	err := s.Set(key, value)
	if err != nil {
		panic(err)
	}
}

// MustSetMap performs as function SetMap, but it panics if any error occurs.
func (s *Session) MustSetMap(data map[string]interface{}) {
	err := s.SetMap(data)
	if err != nil {
		panic(err)
	}
}

// MustContains performs as function Contains, but it panics if any error occurs.
func (s *Session) MustContains(key string) bool {
	b, err := s.Contains(key)
	if err != nil {
		panic(err)
	}
	return b
}

// MustData performs as function Data, but it panics if any error occurs.
func (s *Session) MustData() map[string]interface{} {
	m, err := s.Data()
	if err != nil {
		panic(err)
	}
	return m
}

// MustSize performs as function Size, but it panics if any error occurs.
func (s *Session) MustSize() int {
	size, err := s.Size()
	if err != nil {
		panic(err)
	}
	return size
}

// MustRemove performs as function Remove, but it panics if any error occurs.
func (s *Session) MustRemove(keys ...string) {
	err := s.Remove(keys...)
	if err != nil {
		panic(err)
	}
}
