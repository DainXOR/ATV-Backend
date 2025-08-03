package types

type Pair[T, U any] struct {
	First  T
	Second U
}
type SPair[T any] struct {
	First  T
	Second T
}

func NewPair[T, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{First: first, Second: second}
}
func NewSPair[T any](first, second T) SPair[T] {
	return SPair[T]{First: first, Second: second}
}
