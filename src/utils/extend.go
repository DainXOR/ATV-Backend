package utils

import "strings"

type Predicate[T any] func(T) bool

func Filter[T any](slice []T, predicate Predicate[T]) []T {
	var result []T
	for _, value := range slice {
		if predicate(value) {
			result = append(result, value)
		}
	}
	return result
}
func Map[T, U any](slice []T, mapper func(T) U) []U {
	var result []U
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

func Apply[T, R any](fn func(T) R, value T) func() R {
	return func() R {
		return fn(value)
	}
}

func Capture(prefix string, text string, suffix string) string {
	size := len(text)

	if size < len(prefix)+len(suffix) {
		return ""
	}

	for i := range size {
		if strings.HasPrefix(text[i:], prefix) {
			for j := size; j >= 0; j-- {
				if strings.HasSuffix(text[i:j], suffix) {
					return text[i+1 : j-1]
				}
			}
		}
	}

	return ""
}
