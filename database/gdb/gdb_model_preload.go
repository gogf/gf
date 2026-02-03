// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// relationFieldInfo holds information about a relation field in a struct.
// It caches parsed information to avoid repeated reflection and parsing.
type relationFieldInfo struct {
	Field       gstructs.Field                  // Field information from gstructs
	ParsedTag   parseWithTagInFieldStructOutput // Parsed ORM tag information
	isSliceType bool                            // Cached: whether the field is a slice type
	sourceField string                          // Cached: source field name in the related table (DB column)
	targetField string                          // Cached: target field name in the current struct
}

// isSlice returns whether this relation field is a slice type (one-to-many).
func (r *relationFieldInfo) isSlice() bool {
	return r.isSliceType
}

// newRelationFieldInfo creates a new relationFieldInfo with pre-computed cache values.
func newRelationFieldInfo(field gstructs.Field, parsedTag parseWithTagInFieldStructOutput) *relationFieldInfo {
	info := &relationFieldInfo{
		Field:     field,
		ParsedTag: parsedTag,
	}

	// Pre-compute slice type
	kind := field.Type().Kind()
	if kind == reflect.Pointer {
		elem := field.Type().Elem()
		if elem != nil {
			kind = elem.Kind()
		}
	}
	info.isSliceType = kind == reflect.Slice || kind == reflect.Array

	// Pre-parse relation field pair from "with" tag
	// Format: "source_field=target_field" or just "field_name" (same name in both tables)
	parts := gstr.SplitAndTrim(parsedTag.With, "=")
	if len(parts) == 1 {
		// Same field name in both tables
		info.sourceField = parts[0]
		info.targetField = parts[0]
	} else if len(parts) >= 2 {
		// Different field names
		info.sourceField = parts[0]
		info.targetField = parts[1]
	}

	return info
}

// preloadContext holds the context for recursive preload operations.
// It tracks visited types to detect circular references using backtracking algorithm.
type preloadContext struct {
	model        *Model               // Parent Model (reuses its db, hook, cache configurations)
	visitedTypes map[string]bool      // Visited types for circular reference detection (backtracking)
	allRelations []*relationFieldInfo // All relation fields (for chunkName group lookup)
}

// batchQueryResult holds the result of a batch query for a relation field.
type batchQueryResult struct {
	FieldName string            // Name of the relation field
	DataMap   map[string]Result // Map from relation key (as string) to query results
	Error     error             // Query error if any
}

// doPreloadScan is the entry point for preload mode scanning.
// It performs batch recursive scanning for association operations to solve the N+1 problem.
func (m *Model) doPreloadScan(pointer any) error {
	// Create preload context
	ctx := &preloadContext{
		model:        m,
		visitedTypes: make(map[string]bool),
	}
	return ctx.recursivePreload(pointer)
}

// recursivePreload performs recursive preload operations on the given pointer.
// It collects all relation fields, executes batch queries, maps results, and recursively processes nested relations.
// Circular references are detected using a backtracking algorithm with visitedTypes map.
func (p *preloadContext) recursivePreload(pointer interface{}) error {
	// 1. Get element type and check for circular references
	sliceValue := reflect.ValueOf(pointer)
	if sliceValue.Kind() != reflect.Ptr {
		return gerror.NewCode(gcode.CodeInvalidParameter, "pointer must be a pointer to slice")
	}
	sliceValue = sliceValue.Elem()
	if sliceValue.Kind() != reflect.Slice {
		return gerror.NewCode(gcode.CodeInvalidParameter, "pointer must be a pointer to slice")
	}

	// Empty slice, nothing to do
	if sliceValue.Len() == 0 {
		return nil
	}

	// Get element type: []*Struct -> Struct
	elemType := sliceValue.Type().Elem()
	if elemType.Kind() == reflect.Pointer {
		elemType = elemType.Elem()
	}

	// Check for circular reference using backtracking
	// This allows A->B->C->A structure as long as they're not in the same path
	typeName := elemType.String()
	if p.visitedTypes[typeName] {
		// Already visiting this type in the current path, skip to avoid infinite loop
		return nil
	}
	p.visitedTypes[typeName] = true
	defer delete(p.visitedTypes, typeName) // Backtrack: remove from visited when returning

	// 2. Collect relation fields
	relations, err := p.collectRelations(pointer)
	if err != nil {
		return err
	}
	if len(relations) == 0 {
		return nil // No relations to preload
	}

	// Store all relations for chunkName group lookup
	p.allRelations = relations

	// 3. Batch query all relation fields (sequential execution, no goroutines)
	batchResults := make(map[string]*batchQueryResult)
	for _, relation := range relations {
		result := p.queryRelation(pointer, relation)
		batchResults[relation.Field.Name()] = result
		if result.Error != nil {
			return result.Error
		}
	}

	// 4. Map results to struct fields
	if err := p.mapResults(pointer, relations, batchResults); err != nil {
		return err
	}

	// 5. Recursively process next level
	for _, relation := range relations {
		if err := p.recursivePreloadNext(pointer, relation); err != nil {
			return err
		}
	}

	return nil
}

// collectRelations collects all relation fields from the struct that should be preloaded.
// It uses struct cache to avoid repeated reflection operations.
func (p *preloadContext) collectRelations(pointer interface{}) ([]*relationFieldInfo, error) {
	// Get slice value
	sliceValue := reflect.ValueOf(pointer).Elem()
	if sliceValue.Len() == 0 {
		return nil, nil
	}

	// Get element type
	elemType := sliceValue.Type().Elem()
	if elemType.Kind() == reflect.Pointer {
		elemType = elemType.Elem()
	}

	// Get cached struct info (only caches gstructs reflection results)
	cached, err := getCachedStructInfo(elemType)
	if err != nil {
		return nil, err
	}

	// Iterate fields and parse tags (tag parsing is done every time, not cached)
	var relations []*relationFieldInfo
	for _, field := range cached.fields {
		// Parse tag every time (low cost, maintains flexibility)
		parsedTag := p.model.parseWithTagInFieldStruct(field)
		if parsedTag.With == "" {
			continue // No "with" tag, skip
		}

		// Check if this field should be preloaded
		if !p.shouldPreload(field) {
			continue
		}

		// Create relationFieldInfo (not cached)
		relation := newRelationFieldInfo(field, parsedTag)
		relations = append(relations, relation)
	}

	return relations, nil
}

// shouldPreload checks if a field should be preloaded based on Model configuration.
func (p *preloadContext) shouldPreload(field gstructs.Field) bool {
	// WithAll mode: all fields with "with" tag are allowed
	if p.model.withAll {
		return true
	}

	// With mode: check if field type is in withArray
	fieldTypeStr := gstr.TrimAll(field.Type().String(), "*[]")
	for _, withItem := range p.model.withArray {
		withItemType, err := gstructs.StructType(withItem)
		if err != nil {
			continue
		}
		withItemTypeStr := gstr.TrimAll(withItemType.String(), "*[]")
		if gstr.Compare(fieldTypeStr, withItemTypeStr) == 0 {
			return true
		}
	}

	return false
}

// getChunkConfig returns the chunk configuration for a relation field.
// It follows the priority: API config (by chunkName) > Tag config > chunkName group config.
//
// Design principle: Chunking is opt-in, not default.
// - If no chunk config is provided, use single batch query (no chunking)
// - Only when explicitly configured (Chunked=true or API override), enable chunking
//
// Value semantics:
// - In Tag: Only when both chunkSize and chunkMinRows are configured and > 0, enable chunking
// - In API: -1 means use tag/group config, 0 means disable chunking, >0 means enable with that value
func (p *preloadContext) getChunkConfig(relation *relationFieldInfo) (chunkSize, chunkMinRows int) {
	chunkName := relation.ParsedTag.ChunkName
	chunkSize = -1
	chunkMinRows = -1

	// Priority 1: API configuration (matched by chunkName)
	if chunkName != "" && p.model.preloadOptions != nil {
		if config, ok := p.model.preloadOptions[chunkName]; ok {
			// API config found, use it
			chunkSize = config.ChunkSize
			chunkMinRows = config.ChunkMinRows

			// If both are explicitly set (not -1), use them directly
			if chunkSize != -1 && chunkMinRows != -1 {
				// chunkSize=0 means disable chunking
				if chunkSize == 0 {
					return 0, 0
				}
				return chunkSize, chunkMinRows
			}

			// If only one is set, need to get the other from tag/group
			// Continue to priority 2
		}
	}

	// Priority 2: Tag configuration (only if Chunked=true, meaning both are configured)
	if relation.ParsedTag.Chunked {
		// Both chunkSize and chunkMinRows are configured in tag
		if chunkSize == -1 {
			chunkSize = relation.ParsedTag.ChunkSize
		}
		if chunkMinRows == -1 {
			chunkMinRows = relation.ParsedTag.ChunkMinRows
		}
		return chunkSize, chunkMinRows
	}

	// Priority 3: ChunkName group config (look for other fields with same chunkName)
	if chunkName != "" {
		for _, rel := range p.allRelations {
			if rel.ParsedTag.ChunkName == chunkName && rel != relation && rel.ParsedTag.Chunked {
				// Found a field with same chunkName that has chunk config
				if chunkSize == -1 {
					chunkSize = rel.ParsedTag.ChunkSize
				}
				if chunkMinRows == -1 {
					chunkMinRows = rel.ParsedTag.ChunkMinRows
				}
				if chunkSize > 0 && chunkMinRows > 0 {
					return chunkSize, chunkMinRows
				}
			}
		}
	}

	// No chunk config found, disable chunking (use single batch query)
	return 0, 0
}

// queryRelation executes a batch query for a single relation field.
// It collects all unique relation keys and performs a single WHERE IN query.
func (p *preloadContext) queryRelation(pointer interface{}, relation *relationFieldInfo) *batchQueryResult {
	result := &batchQueryResult{
		FieldName: relation.Field.Name(),
		DataMap:   make(map[string]Result),
	}

	// 1. Collect unique relation key values
	// We need to find the actual struct field name that matches the target field (case-insensitive)
	sliceValue := reflect.ValueOf(pointer).Elem()
	if sliceValue.Len() == 0 {
		return result
	}

	// Get the first item to find field names
	firstItem := sliceValue.Index(0)
	if firstItem.Kind() == reflect.Pointer {
		firstItem = firstItem.Elem()
	}

	// Use cached struct info to find field name
	cached, err := getCachedStructInfo(firstItem.Type())
	if err != nil {
		result.Error = err
		return result
	}

	// Find the actual field name that matches relation.targetField (case-insensitive)
	var actualFieldName string
	for _, field := range cached.fields {
		if utils.EqualFoldWithoutChars(field.Name(), relation.targetField) {
			actualFieldName = field.Name()
			break
		}
	}

	if actualFieldName == "" {
		return result
	}

	targetValues := ListItemValuesUnique(pointer, actualFieldName)
	if len(targetValues) == 0 {
		return result // No values to query
	}

	// 2. Build query model
	// Use the field value to get the correct table name from ORM metadata
	fieldValue := relation.Field.Value
	if fieldValue.Kind() == reflect.Pointer {
		// For pointer types, create a new instance to get metadata
		elemType := fieldValue.Type().Elem()
		fieldValue = reflect.New(elemType)
	} else if fieldValue.Kind() == reflect.Slice {
		// For slice types, get the element type
		elemType := fieldValue.Type().Elem()
		if elemType.Kind() == reflect.Pointer {
			elemType = elemType.Elem()
		}
		fieldValue = reflect.New(elemType)
	}

	model := p.model.db.Model(fieldValue.Interface())
	model = model.Hook(p.model.hookHandler)

	// Apply tag conditions
	if relation.ParsedTag.Where != "" {
		model = model.Where(relation.ParsedTag.Where)
	}
	if relation.ParsedTag.Order != "" {
		model = model.Order(relation.ParsedTag.Order)
	}
	if relation.ParsedTag.Unscoped == "true" {
		model = model.Unscoped()
	}

	// Apply cache if enabled
	if p.model.cacheEnabled && p.model.cacheOption.Name == "" {
		model = model.Cache(p.model.cacheOption)
	}

	// 3. Get chunk configuration (API > Tag > chunkName group > Global default)
	chunkSize, chunkMinRows := p.getChunkConfig(relation)

	// Determine if chunking is needed
	shouldChunk := chunkSize > 0 && len(targetValues) >= chunkMinRows

	// 4. Execute batch query with WHERE IN (with optional chunking)
	var records Result
	if shouldChunk {
		// Execute chunked queries
		for i := 0; i < len(targetValues); i += chunkSize {
			end := i + chunkSize
			if end > len(targetValues) {
				end = len(targetValues)
			}
			chunkValues := targetValues[i:end]

			// IMPORTANT: Clone the model for each chunk to avoid accumulating WHERE conditions
			chunkModel := model.Clone()
			chunkRecords, err := chunkModel.Where(relation.sourceField, chunkValues).All()
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				result.Error = err
				return result
			}
			records = append(records, chunkRecords...)
		}
	} else {
		// Execute single query
		records, result.Error = model.Where(relation.sourceField, targetValues).All()
		if result.Error != nil && !errors.Is(result.Error, sql.ErrNoRows) {
			return result
		}
	}
	result.Error = nil // Ignore ErrNoRows (relation data may not exist)

	// 5. Build map grouped by relation key (use string as key to avoid type mismatch)
	for _, record := range records {
		key := gconv.String(record[relation.sourceField].Interface())
		result.DataMap[key] = append(result.DataMap[key], record)
	}

	return result
}

// mapResults maps batch query results to struct fields.
func (p *preloadContext) mapResults(
	pointer interface{},
	relations []*relationFieldInfo,
	batchResults map[string]*batchQueryResult,
) error {
	sliceValue := reflect.ValueOf(pointer).Elem()
	if sliceValue.Len() == 0 {
		return nil
	}

	firstItem := sliceValue.Index(0)
	if firstItem.Kind() == reflect.Pointer {
		firstItem = firstItem.Elem()
	}
	cached, err := getCachedStructInfo(firstItem.Type())
	if err != nil {
		return err
	}

	for i := 0; i < sliceValue.Len(); i++ {
		item := sliceValue.Index(i)
		if item.Kind() == reflect.Pointer {
			item = item.Elem()
		}

		for _, relation := range relations {
			// Get relation key value from current item - need to use actual field name
			var actualTargetFieldName string
			for _, field := range cached.fields {
				if utils.EqualFoldWithoutChars(field.Name(), relation.targetField) {
					actualTargetFieldName = field.Name()
					break
				}
			}

			if actualTargetFieldName == "" {
				continue
			}

			targetField := item.FieldByName(actualTargetFieldName)
			if !targetField.IsValid() {
				continue
			}
			targetValue := targetField.Interface()
			targetValueStr := gconv.String(targetValue)

			// Get corresponding query results (use string key to avoid type mismatch)
			records := batchResults[relation.Field.Name()].DataMap[targetValueStr]
			if len(records) == 0 {
				continue // No related data
			}

			// Map to field
			fieldValue := item.FieldByName(relation.Field.Name())
			if !fieldValue.IsValid() || !fieldValue.CanSet() {
				continue
			}

			if relation.isSlice() {
				// Slice type: map all records (one-to-many)
				if err := gconv.Scan(records, fieldValue.Addr().Interface()); err != nil {
					return err
				}
			} else {
				// Single type: map only first record (one-to-one)
				// For pointer fields, we need to create a new instance first
				if fieldValue.Kind() == reflect.Pointer {
					// Create new instance of the pointer's element type
					elemType := fieldValue.Type().Elem()
					newElem := reflect.New(elemType)
					if err := gconv.Scan(records[0], newElem.Interface()); err != nil {
						return err
					}
					fieldValue.Set(newElem)
				} else {
					// For non-pointer fields, scan directly
					if err := gconv.Scan(records[0], fieldValue.Addr().Interface()); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// recursivePreloadNext recursively processes the next level of relations.
func (p *preloadContext) recursivePreloadNext(pointer interface{}, relation *relationFieldInfo) error {
	sliceValue := reflect.ValueOf(pointer).Elem()

	if relation.isSlice() {
		// For slice type relations, collect all child slices and merge them into one big slice
		// This allows batch processing of all nested records together

		// Get the element type of the slice
		var sliceElemType reflect.Type
		for i := 0; i < sliceValue.Len(); i++ {
			item := sliceValue.Index(i)
			if item.Kind() == reflect.Pointer {
				item = item.Elem()
			}

			fieldValue := item.FieldByName(relation.Field.Name())
			if !fieldValue.IsValid() || fieldValue.IsZero() {
				continue
			}

			if fieldValue.Kind() == reflect.Pointer {
				fieldValue = fieldValue.Elem()
			}

			if fieldValue.Kind() == reflect.Slice && fieldValue.Len() > 0 {
				sliceElemType = fieldValue.Type().Elem()
				break
			}
		}

		if sliceElemType == nil {
			return nil // No valid slice found
		}

		// Create a merged slice to hold all child records
		// IMPORTANT: We need to keep references to the original slices so that
		// modifications to the merged slice will be reflected in the original data
		mergedSliceType := reflect.SliceOf(sliceElemType)
		mergedSlice := reflect.MakeSlice(mergedSliceType, 0, sliceValue.Len()*10) // Pre-allocate with estimated capacity

		// Collect all child records from all parent records
		// We append the actual slice elements (which are pointers), so modifications will be reflected
		for i := 0; i < sliceValue.Len(); i++ {
			item := sliceValue.Index(i)
			if item.Kind() == reflect.Pointer {
				item = item.Elem()
			}

			fieldValue := item.FieldByName(relation.Field.Name())
			if !fieldValue.IsValid() || fieldValue.IsZero() {
				continue
			}

			if fieldValue.Kind() == reflect.Pointer {
				fieldValue = fieldValue.Elem()
			}

			if fieldValue.Kind() == reflect.Slice && fieldValue.Len() > 0 {
				// Append all elements from this child slice to the merged slice
				// Since the elements are pointers, modifications to them will be reflected in the original slice
				for j := 0; j < fieldValue.Len(); j++ {
					mergedSlice = reflect.Append(mergedSlice, fieldValue.Index(j))
				}
			}
		}

		// If we have collected records, recursively process them as one batch
		if mergedSlice.Len() > 0 {
			// Create a pointer to the merged slice
			mergedSlicePtr := reflect.New(mergedSliceType)
			mergedSlicePtr.Elem().Set(mergedSlice)

			if err := p.recursivePreload(mergedSlicePtr.Interface()); err != nil {
				return err
			}

			// IMPORTANT: Since we're working with pointers, the modifications made by recursivePreload
			// are automatically reflected in the original parent slices. No need to copy back.
		}
	} else {
		// For single type relations, collect all non-nil values into a temporary slice
		// Get the element type of the pointer field
		fieldType := relation.Field.Type().Type
		if fieldType.Kind() != reflect.Pointer {
			return nil // Not a pointer field, skip
		}

		// Create a slice to hold all non-nil pointer values
		sliceType := reflect.SliceOf(fieldType)
		tempSlice := reflect.MakeSlice(sliceType, 0, sliceValue.Len())

		// Collect all non-nil pointer values
		for i := 0; i < sliceValue.Len(); i++ {
			item := sliceValue.Index(i)
			if item.Kind() == reflect.Pointer {
				item = item.Elem()
			}

			fieldValue := item.FieldByName(relation.Field.Name())
			if !fieldValue.IsValid() || fieldValue.IsZero() {
				continue
			}

			// Append the pointer value directly (it's already a pointer)
			if fieldValue.Kind() == reflect.Pointer && !fieldValue.IsNil() {
				tempSlice = reflect.Append(tempSlice, fieldValue)
			}
		}

		// If we have collected values, recursively process them
		if tempSlice.Len() > 0 {
			// Create a pointer to the temporary slice
			tempSlicePtr := reflect.New(sliceType)
			tempSlicePtr.Elem().Set(tempSlice)

			if err := p.recursivePreload(tempSlicePtr.Interface()); err != nil {
				return err
			}

			// IMPORTANT: Since we're working with pointers, the modifications are automatically reflected
		}
	}

	return nil
}
