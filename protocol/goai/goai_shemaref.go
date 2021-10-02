package goai

import (
	"context"
	"fmt"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"reflect"
)

type SchemaRefs []SchemaRef

type SchemaRef struct {
	Ref   string
	Value *Schema
}

func (oai *OpenApiV3) newSchemaRefWithGolangType(golangType reflect.Type, tagMap map[string]string) SchemaRef {
	var (
		oaiType   = oai.golangTypeToOAIType(golangType)
		schemaRef = SchemaRef{}
		schema    = &Schema{
			Type: oaiType,
		}
	)
	if len(tagMap) > 0 {
		err := gconv.Struct(tagMap, schema)
		if err != nil {
			intlog.Error(context.TODO(), err)
		}
	}
	schemaRef.Value = schema
	switch oaiType {
	case
		TypeNumber,
		TypeString,
		TypeBoolean:
		// Nothing to do.

	case
		TypeArray:
		var (
			subSchemaRef = oai.newSchemaRefWithGolangType(golangType.Elem(), nil)
		)
		schema.Items = &subSchemaRef

	case
		TypeObject:
		var (
			structTypeName = gstr.SubStrFromREx(golangType.String(), ".")
		)
		// Specially for map type.
		if golangType.Kind() == reflect.Map {
			var (
				subSchemaRef = oai.newSchemaRefWithGolangType(golangType.Elem(), nil)
			)
			schema.AdditionalProperties = &subSchemaRef
			return schemaRef
		}
		// Normal struct object.
		if _, ok := oai.Components.Schemas[structTypeName]; !ok {
			oai.addSchema(reflect.New(golangType).Interface())
		}
		schemaRef.Ref = structTypeName
		schemaRef.Value = nil
	}
	return schemaRef
}

func (r SchemaRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return []byte(fmt.Sprintf(`{"$ref":"#/components/schemas/%s"}`, r.Ref)), nil
	}
	return json.Marshal(r.Value)
}
