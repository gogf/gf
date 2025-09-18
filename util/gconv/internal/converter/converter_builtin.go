// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"reflect"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

func (c *Converter) builtInAnyConvertFuncForInt64(from any, to reflect.Value) error {
	v, err := c.Int64(from)
	if err != nil {
		return err
	}
	to.SetInt(v)
	return nil
}

func (c *Converter) builtInAnyConvertFuncForUint64(from any, to reflect.Value) error {
	v, err := c.Uint64(from)
	if err != nil {
		return err
	}
	to.SetUint(v)
	return nil
}

func (c *Converter) builtInAnyConvertFuncForString(from any, to reflect.Value) error {
	v, err := c.String(from)
	if err != nil {
		return err
	}
	to.SetString(v)
	return nil
}

func (c *Converter) builtInAnyConvertFuncForFloat64(from any, to reflect.Value) error {
	v, err := c.Float64(from)
	if err != nil {
		return err
	}
	to.SetFloat(v)
	return nil
}

func (c *Converter) builtInAnyConvertFuncForBool(from any, to reflect.Value) error {
	v, err := c.Bool(from)
	if err != nil {
		return err
	}
	to.SetBool(v)
	return nil
}

func (c *Converter) builtInAnyConvertFuncForBytes(from any, to reflect.Value) error {
	v, err := c.Bytes(from)
	if err != nil {
		return err
	}
	to.SetBytes(v)
	return nil
}

func (c *Converter) builtInAnyConvertFuncForTime(from any, to reflect.Value) error {
	t, err := c.Time(from)
	if err != nil {
		return err
	}
	*to.Addr().Interface().(*time.Time) = t
	return nil
}

// builtInAnyConvertFuncForGTime converts any type to *gtime.Time.
//
// THEORETICAL BASIS AND PRINCIPLES:
//
// This function implements a type-specific conversion strategy based on the principle
// that different input types require different handling approaches to preserve semantic
// meaning, particularly timezone information in temporal data.
//
// CORE PRINCIPLES:
//
//  1. DIRECT TYPE PRESERVATION PRINCIPLE
//     When the source and target types are semantically equivalent (gtime.Time variants),
//     use direct assignment to preserve all metadata including timezone, precision,
//     and calendar information without any intermediate transformations.
//
//  2. STRUCTURED DATA EXTRACTION PRINCIPLE
//     When the source is a structured container (map) containing temporal data,
//     extract the actual temporal value and convert it directly rather than
//     serializing the entire container, which would lose semantic context.
//
//  3. MINIMAL TRANSFORMATION PRINCIPLE
//     Apply the least amount of transformation necessary to achieve type compatibility,
//     reducing opportunities for information loss during conversion.
//
//  4. FALLBACK WITH PRESERVATION PRINCIPLE
//     For unknown types, use enhanced general conversion that attempts to preserve
//     timezone information through improved string representations (RFC3339).
//
// CONVERSION PATHS AND RATIONALE:
//
// Path 1: gtime.Time -> gtime.Time (Direct Assignment)
//   - Rationale: Same semantic type, zero transformation needed
//   - Preserves: Timezone, precision, all temporal metadata
//   - Performance: O(1) memory copy operation
//
// Path 2: *gtime.Time -> gtime.Time (Pointer Dereferencing)
//   - Rationale: Pointer wrapper around same semantic type
//   - Preserves: All temporal data after nil safety check
//   - Performance: O(1) with nil check overhead
//
// Path 3: map[string]interface{} -> gtime.Time (Value Extraction)
//   - Rationale: ORM results typically contain temporal data in map structures
//   - Problem Solved: Prevents lossy map->string->time conversion chain
//   - Preserves: Timezone by extracting and converting actual gtime value
//   - Performance: O(1) for single-entry maps (common case)
//
// Path 4: Other Types -> gtime.Time (Enhanced General Conversion)
//   - Rationale: Fallback for unknown types with best-effort preservation
//   - Uses: Enhanced c.GTime() with RFC3339 timezone support
//   - Preserves: Timezone where possible through improved string handling
func (c *Converter) builtInAnyConvertFuncForGTime(from any, to reflect.Value) error {
	// Helper function to efficiently set gtime.Time value to reflect.Value
	// Avoids repeated CanAddr() checks and optimizes assignment operations
	setGTimeValue := func(gtimeVal gtime.Time) {
		if to.CanAddr() {
			*to.Addr().Interface().(*gtime.Time) = gtimeVal
		} else {
			to.Set(reflect.ValueOf(gtimeVal))
		}
	}

	// Cached zero value to avoid repeated gtime.New() allocations
	var zeroGTime = *gtime.New()

	// CONVERSION PATH 1: Direct gtime.Time Assignment (Fast Path)
	// Most common cases handled first for optimal performance
	switch v := from.(type) {
	case *gtime.Time:
		if v != nil {
			// Direct memory copy preserves timezone, precision, and all metadata
			setGTimeValue(*v)
		} else {
			// Nil pointer safety: Use cached zero value
			setGTimeValue(zeroGTime)
		}
		return nil

	case gtime.Time:
		// Direct value assignment for non-pointer gtime types
		// Preserves all temporal information without transformation
		setGTimeValue(v)
		return nil

	// CONVERSION PATH 2: Structured Data Value Extraction (Optimized)
	// Theoretical basis: Extract semantic content from containers rather than
	// serializing containers themselves, which loses semantic context
	case map[string]any:
		// Common in ORM scenarios: {"column_name": gtime_value}
		// Instead of converting entire map to string (lossy), extract the gtime value
		if len(v) > 0 {
			for _, value := range v {
				// Fast path for direct gtime types in map to avoid recursive c.GTime() call
				switch gtimeVal := value.(type) {
				case *gtime.Time:
					if gtimeVal != nil {
						setGTimeValue(*gtimeVal)
					} else {
						setGTimeValue(zeroGTime)
					}
					return nil
				case gtime.Time:
					setGTimeValue(gtimeVal)
					return nil
				default:
					// Only fall back to c.GTime() for non-gtime types
					gtimeResult, err := c.GTime(value)
					if err != nil {
						return err
					}
					if gtimeResult != nil {
						setGTimeValue(*gtimeResult)
					} else {
						setGTimeValue(zeroGTime)
					}
					return nil
				}
			}
		}
		// Empty map case: Use cached zero value
		setGTimeValue(zeroGTime)
		return nil
	}

	// CONVERSION PATH 3: Enhanced General Conversion
	// Theoretical basis: For unknown types, use enhanced converter that attempts
	// timezone preservation through improved string representations and parsing
	gtimeResult, err := c.GTime(from)
	if err != nil {
		return err
	}
	if gtimeResult != nil {
		setGTimeValue(*gtimeResult)
	} else {
		setGTimeValue(zeroGTime)
	}
	return nil
}
