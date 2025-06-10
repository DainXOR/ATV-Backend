package utils

import (
	"fmt"
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

// Function to create an instance of the element type of a slice
// It returns a pointer to the new instance or an error if the input is not a slice or pointer to slice
func SliceElemInstance(slice any) (any, error) {
	if slice == nil {
		return nil, fmt.Errorf("Input is nil, expected a slice or pointer to slice")
	}

	t := reflect.TypeOf(slice)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Expected a slice or pointer to slice, got: %s", t.Kind())
	}

	elemType := t.Elem()

	var instance reflect.Value
	if elemType.Kind() == reflect.Ptr {
		instance = reflect.New(elemType.Elem())
	} else {
		instance = reflect.New(elemType)
	}

	return instance.Interface(), nil
}
