package promise

import (
	"errors"
)

var ErrNoResult = errors.New("not result")

type Result[T any] struct {
	res   T
	err   error
	dirty bool
}

func resolve[T any](t T) *Result[T] {
	return &Result[T]{res: t, dirty: true}
}

func reject[T any](err error) *Result[T] {
	return &Result[T]{err: err, dirty: true}
}

func (r *Result[T]) Result() (t T) {
	if r.IsZero() {
		return
	}

	return r.res
}

func (r *Result[T]) Error() error {
	if r.IsZero() {
		return ErrNoResult
	}

	return r.err
}

func (r *Result[T]) IsZero() bool {
	return r == nil || !r.dirty
}

func (r *Result[T]) Ok() bool {
	return r.Error() == nil
}

func (r *Result[T]) Unwrap() (T, error) {
	return r.Result(), r.Error()
}
