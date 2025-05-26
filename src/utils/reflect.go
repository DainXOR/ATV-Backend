package utils

import (
	"reflect"
	"strings"
)

// Function to get the structure of a struct as a string
func StructToString(obj any) string {
	t := reflect.TypeOf(obj)
	//v := reflect.ValueOf(obj)

	if t.Kind() != reflect.Struct {
		return "{ }"
	}

	var sb strings.Builder
	sb.WriteString("{")

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		sb.WriteString(string(field.Tag.Get("json")))
		sb.WriteString(": ")
		sb.WriteString(field.Type.String())

		if i < t.NumField()-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString("}")
	return sb.String()
}

// Function to convert a struct to a map[string]any
// It takes a filter function to determine which fields to include in the map
func StructToMap(obj any, filter func(reflect.StructField, reflect.Value) bool) map[string]any {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	if t.Kind() != reflect.Struct {
		return nil
	}

	result := make(map[string]any)

	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i)

		if filter != nil && filter(field, value) {
			result[string(field.Tag.Get("json"))] = value.Interface()
		}

	}

	return result
}
