package types

type Result[T any] struct {
	value Optional[T]
	err   error
}

func ResultErr[T any](err error) Result[T] {
	return Result[T]{value: OptionalEmpty[T](), err: err}
}
func ResultOk[T any](value T) Result[T] {
	return Result[T]{value: OptionalOf(value), err: nil}
}
func ResultOf[T any](value T, err error, condition bool) Result[T] {
	if condition {
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

func (r *Result[T]) Value() T {
	return r.value.Get()
}
func (r *Result[T]) ValueOr(value T) T {
	if r.IsOk() {
		return r.Value()
	}
	return value
}

func (r *Result[T]) Error() error {
	return r.err
}
func (r *Result[T]) ErrorOr(err error) error {
	if r.IsErr() {
		return r.Error()
	}
	return err
}

func (r *Result[T]) GetRaw() (T, error) {
	return r.Value(), r.Error()
}
func (r *Result[T]) Get() (Optional[T], error) {
	return r.value, r.err
}
