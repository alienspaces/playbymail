package jsonschema

import (
	"reflect"
	"strings"
)

type JSONType struct {
	ElemType string // string | number | boolean
	IsArray  bool
}

// CreateJSONTypeMap takes a struct with JSON tagged fields to creates a map of JSON field names to their JSON element type classification (string, number, or boolean).
func CreateJSONTypeMap(v any) map[string]JSONType {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return addFieldsToMap(t, make(map[string]JSONType))
}

func addFieldsToMap(t reflect.Type, result map[string]JSONType) map[string]JSONType {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		if field.Anonymous {
			result = addFieldsToMap(fieldType, result)
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		fieldName := strings.Split(jsonTag, ",")[0] // remove omitempty tag value if present

		switch fieldType.Kind() {
		case reflect.Slice, reflect.Array:
			jsontype := classifyKind(fieldType.Elem().Kind())
			jsontype.IsArray = true
			result[fieldName] = jsontype
		default:
			result[fieldName] = classifyKind(fieldType.Kind())
		}
	}

	return result
}

func classifyKind(k reflect.Kind) JSONType {
	switch k {
	case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return JSONType{ElemType: "number"}
	case reflect.Bool:
		return JSONType{ElemType: "boolean"}
	default:
		return JSONType{ElemType: "string"}
	}
}
