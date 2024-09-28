// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

type memoryData struct {
	mu   sync.RWMutex                   // dataMu ensures the concurrent safety of underlying data map.
	data map[interface{}]memoryDataItem // data is the underlying cache data which is stored in a hash table.
}

// memoryDataItem holds the internal cache item data.
type memoryDataItem struct {
	v interface{} // Value.
	e int64       // Expire timestamp in milliseconds.
}

func newMemoryData() *memoryData {
	return &memoryData{
		data: make(map[interface{}]memoryDataItem),
	}
}

// Update updates the value of `key` without changing its expiration and returns the old value.
// The returned value `exist` is false if the `key` does not exist in the cache.
//
// It deletes the `key` if given `value` is nil.
// It does nothing if `key` does not exist in the cache.
func (d *memoryData) Update(key interface{}, value interface{}) (oldValue interface{}, exist bool, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if item, ok := d.data[key]; ok {
		d.data[key] = memoryDataItem{
			v: value,
			e: item.e,
		}
		return item.v, true, nil
	}
	return nil, false, nil
}

// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
//
// It returns -1 and does nothing if the `key` does not exist in the cache.
// It deletes the `key` if `duration` < 0.
func (d *memoryData) UpdateExpire(key interface{}, expireTime int64) (oldDuration time.Duration, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if item, ok := d.data[key]; ok {
		d.data[key] = memoryDataItem{
			v: item.v,
			e: expireTime,
		}
		return time.Duration(item.e-gtime.TimestampMilli()) * time.Millisecond, nil
	}
	return -1, nil
}

// Remove deletes the one or more keys from cache, and returns its value.
// If multiple keys are given, it returns the value of the deleted last item.
func (d *memoryData) Remove(keys ...interface{}) (removedKeys []interface{}, value interface{}, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	removedKeys = make([]interface{}, 0)
	for _, key := range keys {
		item, ok := d.data[key]
		if ok {
			value = item.v
			delete(d.data, key)
			removedKeys = append(removedKeys, key)
		}
	}
	return removedKeys, value, nil
}

// Data returns a copy of all key-value pairs in the cache as map type.
func (d *memoryData) Data() (map[interface{}]interface{}, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var (
		data     = make(map[interface{}]interface{}, len(d.data))
		nowMilli = gtime.TimestampMilli()
	)
	for k, v := range d.data {
		if v.e > nowMilli {
			data[k] = v.v
		}
	}
	return data, nil
}

// Keys returns all keys in the cache as slice.
func (d *memoryData) Keys() ([]interface{}, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var (
		keys     = make([]interface{}, 0, len(d.data))
		nowMilli = gtime.TimestampMilli()
	)
	for k, v := range d.data {
		if v.e > nowMilli {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

// Values returns all values in the cache as slice.
func (d *memoryData) Values() ([]interface{}, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var (
		values   = make([]interface{}, 0, len(d.data))
		nowMilli = gtime.TimestampMilli()
	)
	for _, v := range d.data {
		if v.e > nowMilli {
			values = append(values, v.v)
		}
	}
	return values, nil
}

// Size returns the size of the cache that not expired.
func (d *memoryData) Size() (size int, err error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var nowMilli = gtime.TimestampMilli()
	for _, v := range d.data {
		if v.e > nowMilli {
			size++
		}
	}
	return size, nil
}

// Clear clears all data of the cache.
// Note that this function is sensitive and should be carefully used.
func (d *memoryData) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data = make(map[interface{}]memoryDataItem)
}

func (d *memoryData) Get(key interface{}) (item memoryDataItem, ok bool) {
	d.mu.RLock()
	item, ok = d.data[key]
	d.mu.RUnlock()
	return
}

func (d *memoryData) Set(key interface{}, value memoryDataItem) {
	d.mu.Lock()
	d.data[key] = value
	d.mu.Unlock()
}

// SetMap batch sets cache with key-value pairs by `data`, which is expired after `duration`.
//
// It does not expire if `duration` == 0.
// It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
func (d *memoryData) SetMap(data map[interface{}]interface{}, expireTime int64) error {
	d.mu.Lock()
	for k, v := range data {
		d.data[k] = memoryDataItem{
			v: v,
			e: expireTime,
		}
	}
	d.mu.Unlock()
	return nil
}

func (d *memoryData) SetWithLock(ctx context.Context, key interface{}, value interface{}, expireTimestamp int64) (interface{}, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var (
		err error
	)
	if v, ok := d.data[key]; ok && !v.IsExpired() {
		return v.v, nil
	}
	f, ok := value.(Func)
	if !ok {
		// Compatible with raw function value.
		f, ok = value.(func(ctx context.Context) (value interface{}, err error))
	}
	if ok {
		if value, err = f(ctx); err != nil {
			return nil, err
		}
		if value == nil {
			return nil, nil
		}
	}
	d.data[key] = memoryDataItem{v: value, e: expireTimestamp}
	return value, nil
}

func (d *memoryData) Delete(key interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.data, key)
}
