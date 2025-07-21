package types

type Result[T any] struct {
	value Optional[T]
	err   error
}

func ResultErr[T any](err error) Result[T] {
	return Result[T]{value: OptionalEmpty[T](), err: err}
}
func ResultOk[T any](value T) Result[T] {
	res := Result[T]{value: OptionalOf(value), err: nil}
	return res
}
func ResultOf[T any](value T, err error, isError bool) Result[T] {
	if !isError {
		return ResultOk(value)
	}
	return ResultErr[T](err)
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r Result[T]) Value() T {
	return r.value.Get()
}
func (r Result[T]) ValueOr(value T) T {
	if r.IsOk() {
		return r.Value()
	}
	return value
}

func (r Result[T]) Error() error {
	return r.err
}
func (r Result[T]) ErrorOr(err error) error {
	if r.IsErr() {
		return r.Error()
	}
	return err
}

func (r Result[T]) GetRaw() (T, error) {
	return r.Value(), r.Error()
}
func (r Result[T]) Get() (Optional[T], error) {
	return r.value, r.err
}

func (r Result[T]) Then(fn func(T) Result[T]) Result[T] {
	if r.IsOk() {
		return fn(r.Value())
	}
	return r
}
func (r Result[T]) Or(fn func(error) Result[T]) Result[T] {
	if r.IsErr() {
		return fn(r.Error())
	}
	return r
}
