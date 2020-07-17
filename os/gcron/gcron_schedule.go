// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gcron

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jin502437344/gf/text/gregex"
)

// 运行时间管理对象
type cronSchedule struct {
	create  int64  // 创建时间戳(秒)
	every   int64  // 运行时间间隔(秒)
	pattern string // 原始注册字符串
	second  map[int]struct{}
	minute  map[int]struct{}
	hour    map[int]struct{}
	day     map[int]struct{}
	week    map[int]struct{}
	month   map[int]struct{}
}

const (
	gREGEX_FOR_CRON = `^([\-/\d\*\?,]+)\s+([\-/\d\*\?,]+)\s+([\-/\d\*\?,]+)\s+([\-/\d\*\?,]+)\s+([\-/\d\*\?,]+)\s+([\-/\d\*\?,]+)$`
)

var (
	// 预定义的定时格式
	predefinedPatternMap = map[string]string{
		"@yearly":   "0 0 0 1 1 *",
		"@annually": "0 0 0 1 1 *",
		"@monthly":  "0 0 0 1 * *",
		"@weekly":   "0 0 0 * * 0",
		"@daily":    "0 0 0 * * *",
		"@midnight": "0 0 0 * * *",
		"@hourly":   "0 0 * * * *",
	}
	// 月份与数字对应表
	monthMap = map[string]int{
		"jan": 1,
		"feb": 2,
		"mar": 3,
		"apr": 4,
		"may": 5,
		"jun": 6,
		"jul": 7,
		"aug": 8,
		"sep": 9,
		"oct": 10,
		"nov": 11,
		"dec": 12,
	}
	// 星期与数字对应表
	weekMap = map[string]int{
		"sun": 0,
		"mon": 1,
		"tue": 2,
		"wed": 3,
		"thu": 4,
		"fri": 5,
		"sat": 6,
	}
)

// 解析定时格式为cronSchedule对象
func newSchedule(pattern string) (*cronSchedule, error) {
	// 处理预定义的定时格式
	if match, _ := gregex.MatchString(`(@\w+)\s*(\w*)\s*`, pattern); len(match) > 0 {
		key := strings.ToLower(match[1])
		if v, ok := predefinedPatternMap[key]; ok {
			pattern = v
		} else if strings.Compare(key, "@every") == 0 {
			if d, err := time.ParseDuration(match[2]); err != nil {
				return nil, err
			} else {
				return &cronSchedule{
					create:  time.Now().Unix(),
					every:   int64(d.Seconds()),
					pattern: pattern,
				}, nil
			}
		} else {
			return nil, errors.New(fmt.Sprintf(`invalid pattern: "%s"`, pattern))
		}
	}
	// 处理通用的定时格式定义
	if match, _ := gregex.MatchString(gREGEX_FOR_CRON, pattern); len(match) == 7 {
		schedule := &cronSchedule{
			create:  time.Now().Unix(),
			every:   0,
			pattern: pattern,
		}
		// 秒
		if m, err := parseItem(match[1], 0, 59, false); err != nil {
			return nil, err
		} else {
			schedule.second = m
		}
		// 分
		if m, err := parseItem(match[2], 0, 59, false); err != nil {
			return nil, err
		} else {
			schedule.minute = m
		}
		// 时
		if m, err := parseItem(match[3], 0, 23, false); err != nil {
			return nil, err
		} else {
			schedule.hour = m
		}
		// 天
		if m, err := parseItem(match[4], 1, 31, true); err != nil {
			return nil, err
		} else {
			schedule.day = m
		}
		// 月
		if m, err := parseItem(match[5], 1, 12, false); err != nil {
			return nil, err
		} else {
			schedule.month = m
		}
		// 周
		if m, err := parseItem(match[6], 0, 6, true); err != nil {
			return nil, err
		} else {
			schedule.week = m
		}
		return schedule, nil
	} else {
		return nil, errors.New(fmt.Sprintf(`invalid pattern: "%s"`, pattern))
	}
}

// 解析定时格式中的每一项定时配置
func parseItem(item string, min int, max int, allowQuestionMark bool) (map[int]struct{}, error) {
	m := make(map[int]struct{}, max-min+1)
	if item == "*" || (allowQuestionMark && item == "?") {
		for i := min; i <= max; i++ {
			m[i] = struct{}{}
		}
	} else {
		for _, item := range strings.Split(item, ",") {
			interval := 1
			intervalArray := strings.Split(item, "/")
			if len(intervalArray) == 2 {
				if i, err := strconv.Atoi(intervalArray[1]); err != nil {
					return nil, errors.New(fmt.Sprintf(`invalid pattern item: "%s"`, item))
				} else {
					interval = i
				}
			}
			rangeMin := min
			rangeMax := max
			rangeArray := strings.Split(intervalArray[0], "-")
			valueType := byte(0)
			switch max {
			case 6:
				valueType = 'w'
			case 11:
				valueType = 'm'
			}
			// 例如: */5
			if rangeArray[0] != "*" {
				if i, err := parseItemValue(rangeArray[0], valueType); err != nil {
					return nil, errors.New(fmt.Sprintf(`invalid pattern item: "%s"`, item))
				} else {
					rangeMin = i
					rangeMax = i
				}
			}
			if len(rangeArray) == 2 {
				if i, err := parseItemValue(rangeArray[1], valueType); err != nil {
					return nil, errors.New(fmt.Sprintf(`invalid pattern item: "%s"`, item))
				} else {
					rangeMax = i
				}
			}
			for i := rangeMin; i <= rangeMax; i += interval {
				m[i] = struct{}{}
			}
		}
	}
	return m, nil
}

// 将配置项值转换为数字
func parseItemValue(value string, valueType byte) (int, error) {
	if gregex.IsMatchString(`^\d+$`, value) {
		// 纯数字
		if i, err := strconv.Atoi(value); err == nil {
			return i, nil
		}
	} else {
		// 英文字母
		switch valueType {
		case 'm':
			if i, ok := monthMap[strings.ToLower(value)]; ok {
				return int(i), nil
			}
		case 'w':
			if i, ok := weekMap[strings.ToLower(value)]; ok {
				return int(i), nil
			}
		}
	}
	return 0, errors.New(fmt.Sprintf(`invalid pattern value: "%s"`, value))
}

// 判断给定的时间是否满足schedule
func (s *cronSchedule) meet(t time.Time) bool {
	if s.every != 0 {
		diff := t.Unix() - s.create
		if diff > 0 {
			return diff%s.every == 0
		}
		return false
	} else {
		if _, ok := s.second[t.Second()]; !ok {
			return false
		}
		if _, ok := s.minute[t.Minute()]; !ok {
			return false
		}
		if _, ok := s.hour[t.Hour()]; !ok {
			return false
		}
		if _, ok := s.day[t.Day()]; !ok {
			return false
		}
		if _, ok := s.month[int(t.Month())]; !ok {
			return false
		}
		if _, ok := s.week[int(t.Weekday())]; !ok {
			return false
		}
		return true
	}
}
