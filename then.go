package promise

func Then[T any, R any](promise *Promise[T], resolver func(T) *Promise[R]) *Promise[R] {

	return New(func() (R, error) {
		res, err := promise.Await()
		if err != nil {
			return Reject[R](err).Await()
		}

		return resolver(res).Await()
	})
}

func Catch[T any](promise *Promise[T], catcher func(error)) error {
	if _, err := promise.Await(); err != nil {
		catcher(err)
		return err
	}

	return nil
}

func Reject[T any](err error) *Promise[T] {
	return &Promise[T]{
		err: err,
	}
}

func Resolve[T any](t T) *Promise[T] {
	return &Promise[T]{
		res: t,
	}
}
