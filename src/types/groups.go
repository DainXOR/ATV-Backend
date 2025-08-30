package types

type Pair[T, U any] struct {
	First  T
	Second U
}
type SPair[T any] struct {
	Pair[T, T]
}

func NewPair[T, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{First: first, Second: second}
}
func NewSPair[T any](first, second T) SPair[T] {
	return SPair[T]{
		Pair[T, T]{
			First:  first,
			Second: second,
		},
	}
}

type Triplet[T, U, V any] struct {
	First  T
	Second U
	Third  V
}
type STriplet[T any] struct {
	Triplet[T, T, T]
}

func NewTriplet[T, U, V any](first T, second U, third V) Triplet[T, U, V] {
	return Triplet[T, U, V]{First: first, Second: second, Third: third}
}
func NewSTriplet[T any](first, second, third T) STriplet[T] {
	return STriplet[T]{
		Triplet: Triplet[T, T, T]{
			First:  first,
			Second: second,
			Third:  third,
		},
	}
}
