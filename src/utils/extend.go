package utils

import (
	"fmt"
	"strings"
)

type Predicate[T any] func(T) bool

func Apply[T, R any](value T, fn func(T) R) func() R {
	return func() R {
		return fn(value)
	}
}
func Extract(prefix string, text string, suffix string) string {
	size := len(text)

	if size < len(prefix)+len(suffix) {
		return ""
	}

	for i := range size {
		hasPrefix := prefix != "" && strings.HasPrefix(text[i:], prefix)
		if prefix == "" || hasPrefix {
			if hasPrefix {
				i += 1
			}

			for j := size; j >= 0; j-- {
				hasSuffix := suffix != "" && strings.HasSuffix(text[i:j], suffix)
				if suffix == "" || hasSuffix {
					if hasSuffix {
						j -= 1
					}

					return text[i:j]
				}
			}
		}
	}

	return ""
}

func Filter[T any](slice []T, predicate Predicate[T]) []T {
	result := make([]T, 0, len(slice))

	for _, value := range slice {
		if predicate(value) {
			result = append(result, value)
		}
	}
	return result
}
func Map[T, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, 0, len(slice))

	for _, value := range slice {
		result = append(result, mapper(value))
	}
	return result
}
func Reduce[T, U any](slice []T, reducer func(U, T) U, initial U) U {
	result := initial

	for _, value := range slice {
		result = reducer(result, value)
	}
	return result
}
func AsStrings(slice []any) []string {
	return Map(slice, func(e any) string { return fmt.Sprint(e) })
}
func Podate(slice []string, cutset string) []string {
	return Map(slice, func(e string) string {
		return strings.Trim(e, cutset)
	})
}

func Join(v []any, sep string) string {
	values := AsStrings(v)
	joinedArgs := strings.Join(values, sep)

	return strings.TrimSpace(joinedArgs)
}
