// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
)

// CacheOption is options for model cache control in query.
type CacheOption struct {
	// Duration is the TTL for the cache.
	// If the parameter `Duration` < 0, which means it clear the cache with given `Name`.
	// If the parameter `Duration` = 0, which means it never expires.
	// If the parameter `Duration` > 0, which means it expires after `Duration`.
	Duration time.Duration

	// Name is an optional unique name for the cache.
	// The Name is used to bind a name to the cache, which means you can later control the cache
	// like changing the `duration` or clearing the cache with specified Name.
	Name string

	// Force caches the query result whatever the result is nil or not.
	// It is used to avoid Cache Penetration.
	Force bool
}

// selectCacheItem is the cache item for SELECT statement result.
type selectCacheItem struct {
	Result            Result // Sql result of SELECT statement.
	FirstResultColumn string // The first column name of result, for Value/Count functions.
}

// Cache sets the cache feature for the model. It caches the result of the sql, which means
// if there's another same sql request, it just reads and returns the result from cache, it
// but not committed and executed into the database.
//
// Note that, the cache feature is disabled if the model is performing select statement
// on a transaction.
func (m *Model) Cache(option CacheOption) *Model {
	model := m.getModel()
	model.cacheOption = option
	model.cacheEnabled = true
	return model
}

// checkAndRemoveSelectCache checks and removes the cache in insert/update/delete statement if
// cache feature is enabled.
func (m *Model) checkAndRemoveSelectCache(ctx context.Context) {
	if m.cacheEnabled && m.cacheOption.Duration < 0 && len(m.cacheOption.Name) > 0 {
		var cacheKey = m.makeSelectCacheKey("")
		if _, err := m.db.GetCache().Remove(ctx, cacheKey); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
		}
	}
}

func (m *Model) getSelectResultFromCache(ctx context.Context, sql string, args ...interface{}) (result Result, err error) {
	if !m.cacheEnabled || m.tx != nil {
		return
	}
	var (
		ok        bool
		cacheItem *selectCacheItem
		cacheKey  = m.makeSelectCacheKey(sql, args...)
		cacheObj  = m.db.GetCache()
	)
	defer func() {
		if cacheItem != nil {
			if internalData := m.db.GetCore().GetInternalCtxDataFromCtx(ctx); internalData != nil {
				if cacheItem.FirstResultColumn != "" {
					internalData.FirstResultColumn = cacheItem.FirstResultColumn
				}
			}
		}
	}()
	if v, _ := cacheObj.Get(ctx, cacheKey); !v.IsNil() {
		if cacheItem, ok = v.Val().(*selectCacheItem); ok {
			// In-memory cache.
			return cacheItem.Result, nil
		} else {
			// Other cache, it needs conversion.
			if err = json.UnmarshalUseNumber(v.Bytes(), &cacheItem); err != nil {
				return nil, err
			}
			return cacheItem.Result, nil
		}
	}
	return
}

func (m *Model) saveSelectResultToCache(ctx context.Context, result Result, sql string, args ...interface{}) (err error) {
	if !m.cacheEnabled || m.tx != nil {
		return
	}
	var (
		cacheKey = m.makeSelectCacheKey(sql, args...)
		cacheObj = m.db.GetCache()
	)
	if m.cacheOption.Duration < 0 {
		if _, errCache := cacheObj.Remove(ctx, cacheKey); errCache != nil {
			intlog.Errorf(ctx, `%+v`, errCache)
		}
	} else {
		// In case of Cache Penetration.
		if result.IsEmpty() && m.cacheOption.Force {
			result = Result{}
		}
		var cacheItem = &selectCacheItem{
			Result: result,
		}
		if internalData := m.db.GetCore().GetInternalCtxDataFromCtx(ctx); internalData != nil {
			cacheItem.FirstResultColumn = internalData.FirstResultColumn
		}
		if errCache := cacheObj.Set(ctx, cacheKey, cacheItem, m.cacheOption.Duration); errCache != nil {
			intlog.Errorf(ctx, `%+v`, errCache)
		}
	}
	return nil
}

func (m *Model) makeSelectCacheKey(sql string, args ...interface{}) string {
	return m.db.GetCore().makeSelectCacheKey(
		m.cacheOption.Name,
		m.db.GetSchema(),
		m.db.GetCore().guessPrimaryTableName(m.tables),
		sql,
		args...,
	)
}
