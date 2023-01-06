// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	setTimeZoneMu   sync.Mutex
	setTimeZoneName string
	zoneMap         = make(map[string]*time.Location)
	zoneMu          sync.RWMutex
)

// SetTimeZone sets the time zone for current whole process.
// The parameter `zone` is an area string specifying corresponding time zone,
// eg: Asia/Shanghai.
//
// PLEASE VERY NOTE THAT:
// 1. This should be called before package "time" import.
// 2. This function should be called once.
// 3. Please refer to issue: https://github.com/golang/go/issues/34814
func SetTimeZone(zone string) (err error) {
	setTimeZoneMu.Lock()
	defer setTimeZoneMu.Unlock()
	if setTimeZoneName != "" && !strings.EqualFold(zone, setTimeZoneName) {
		return gerror.NewCodef(
			gcode.CodeInvalidOperation,
			`process timezone already set using "%s"`,
			setTimeZoneName,
		)
	}
	defer func() {
		if err == nil {
			setTimeZoneName = zone
		}
	}()

	// Load zone info from specified name.
	location, err := time.LoadLocation(zone)
	if err != nil {
		err = gerror.WrapCodef(gcode.CodeInvalidParameter, err, `time.LoadLocation failed for zone "%s"`, zone)
		return err
	}

	// Update the time.Local for once.
	time.Local = location

	// Update the timezone environment for *nix systems.
	var (
		envKey   = "TZ"
		envValue = location.String()
	)
	if err = os.Setenv(envKey, envValue); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeUnknown,
			err,
			`set environment failed with key "%s", value "%s"`,
			envKey, envValue,
		)
	}
	return
}

// ToLocation converts current time to specified location.
func (t *Time) ToLocation(location *time.Location) *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.In(location)
	return newTime
}

// ToZone converts current time to specified zone like: Asia/Shanghai.
func (t *Time) ToZone(zone string) (*Time, error) {
	if location, err := t.getLocationByZoneName(zone); err == nil {
		return t.ToLocation(location), nil
	} else {
		return nil, err
	}
}

func (t *Time) getLocationByZoneName(name string) (location *time.Location, err error) {
	zoneMu.RLock()
	location = zoneMap[name]
	zoneMu.RUnlock()
	if location == nil {
		location, err = time.LoadLocation(name)
		if err != nil {
			err = gerror.Wrapf(err, `time.LoadLocation failed for name "%s"`, name)
		}
		if location != nil {
			zoneMu.Lock()
			zoneMap[name] = location
			zoneMu.Unlock()
		}
	}
	return
}

// Local converts the time to local timezone.
func (t *Time) Local() *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.Local()
	return newTime
}
