package promise

import "context"

func Then[T any, R any](ctx context.Context, promise *Promise[T], resolver func(context.Context, T) (R, error)) *Promise[R] {
	t, err := promise.Await()
	if err != nil {
		return reject[R](err)
	}

	return New(ctx, func(ctx context.Context) (R, error) {
		return resolver(ctx, t)
	})
}

func Catch[T any](promise *Promise[T], catcher func(error)) error {
	_, err := promise.Await()
	if err != nil {
		catcher(err)
		return err
	}
	return nil
}

func reject[T any](err error) *Promise[T] {
	return &Promise[T]{
		err: err,
	}
}

func resolve[T any](t T) *Promise[T] {
	return &Promise[T]{
		res: t,
	}
}
