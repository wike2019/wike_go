package result

type Result[T any] struct {
	err  error
	data T
}

func NewResult[T any](data T, err error) *Result[T] {
	return &Result[T]{
		err:  err,
		data: data,
	}
}
func (r *Result[T]) Unwrap() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.data
}
func (r *Result[T]) UnwrapDefault(data T) T {
	if r.err != nil {
		return data
	}
	return r.data
}
