package utils

import (
	"dainxor/atv/types"
	"fmt"
	"strings"
)

type Predicate[T any] func(T) bool

func Curry[T, R any](value T, fn func(T) R) func() R {
	return func() R {
		return fn(value)
	}
}
func Partial[T, U, R any](value T, fn func(T, U) R) func(arg U) R {
	return func(arg U) R {
		return fn(value, arg)
	}
}
func Transform[T, R any](value T, fn func(T) R) R {
	return fn(value)
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
func ForEach[T any](slice []T, action func(int, T)) {
	for i, value := range slice {
		action(i, value)
	}
}
func Reduce[T, U any](slice []T, reducer func(U, T) U, initial U) U {
	result := initial

	for _, value := range slice {
		result = reducer(result, value)
	}
	return result
}

func DForEach[K comparable, V any](m map[K]V, action func(K, V)) {
	for k, v := range m {
		action(k, v)
	}
}
func DMap[K, NK comparable, V, NV any](m map[K]V, mapper func(K, V) (NK, NV)) map[NK]NV {
	result := make(map[NK]NV, len(m))

	for k, v := range m {
		newKey, newVal := mapper(k, v)
		result[newKey] = newVal
	}
	return result
}
func DFlatten[K comparable, V, S any](m map[K]V, flattener func(K, V) S) []S {
	result := make([]S, 0, len(m))

	for k, v := range m {
		result = append(result, flattener(k, v))
	}

	return result
}
func DZip[K comparable, V1, V2 any](m1 map[K]V1, m2 map[K]V2) map[K]types.Pair[V1, V2] {
	result := make(map[K]types.Pair[V1, V2], len(m1))

	for k, v1 := range m1 {
		if v2, exists := m2[k]; exists {
			result[k] = types.Pair[V1, V2]{First: v1, Second: v2}
		}
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
