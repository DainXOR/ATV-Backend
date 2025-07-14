package types

type Pair[T, U any] struct {
	First  T
	Second U
}
type SPair[T any] struct {
	First  T
	Second T
}
