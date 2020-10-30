// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package guid provides simple and high performance unique id generation functionality.
//
// Unique String ID:
// PLEASE VERY NOTE:
// This package only provides unique number generation for simple, convenient and most common
// usage purpose, but does not provide strict global unique number generation. Please refer
// to UUID algorithm for global unique number generation if necessary.
//
// Unique Number ID:
// An Improved SnowFlake ID is inspired by Twitter's Snowflake, which is composed of:
//     39 bits for time in units of 10 msec
//      8 bits for a sequence number
//     16 bits for a machine id
package guid
