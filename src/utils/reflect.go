package utils

import (
	"reflect"
	"strings"
)

// Function to get the structure of a struct as a string
func StructToString(obj interface{}) string {
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
