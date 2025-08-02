package types

type Optional[T any] struct {
	value   T
	present bool
}

func OptionalOf[T any](value T, condition ...bool) Optional[T] {
	for _, c := range condition {
		if c {
			return Optional[T]{value: value, present: true}
		}
	}

	return OptionalEmpty[T]()
}

func OptionalEmpty[T any]() Optional[T] {
	var zeroValue T
	return Optional[T]{value: zeroValue, present: false}
}

func (o Optional[T]) IsPresent() bool {
	return o.present
}
func (o Optional[T]) IsEmpty() bool {
	return !o.IsPresent()
}

func (o Optional[T]) Get() T {
	return o.value
}

func (o Optional[T]) GetOr(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

func (o Optional[T]) IfPresent(fn func(T) Optional[T]) Optional[T] {
	if o.present {
		if res := fn(o.value); res.IsPresent() {
			return res
		} else {
			return OptionalEmpty[T]()
		}
	}
	return o
}
func (o Optional[T]) IfEmpty(fn func() Optional[T]) Optional[T] {
	if !o.present {
		if res := fn(); res.IsPresent() {
			return res
		} else {
			return OptionalEmpty[T]()
		}
	}
	return o
}
