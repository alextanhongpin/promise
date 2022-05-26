package promise

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const MaxGoroutine = 1_000

var ErrPanic = errors.New("panic")

type Promise[T any] struct {
	res  T
	err  error
	wg   sync.WaitGroup
	once sync.Once
}

type Task[T any] func(ctx context.Context) (T, error)

func New[T any](ctx context.Context, fn Task[T]) *Promise[T] {
	p := &Promise[T]{}
	p.wg.Add(2)

	done := make(chan bool)

	go func() {
		defer p.wg.Done()
		defer close(done)
		defer func() {
			if err := recover(); err != nil {
				if perr, ok := err.(error); ok {
					p.reject(perr)
				} else {
					p.reject(fmt.Errorf("%w: %v", ErrPanic, err))
				}
			}
		}()

		res, err := fn(ctx)
		if err != nil {
			p.reject(err)
		} else {
			p.resolve(res)
		}
	}()

	go func() {
		defer p.wg.Done()

		select {
		case <-ctx.Done():
			p.reject(ctx.Err())
		case <-done:
		}
	}()

	return p
}

func (p *Promise[T]) Await() (T, error) {
	p.wg.Wait()

	return p.res, p.err
}

func (p *Promise[T]) AwaitResult() *Result[T] {
	res, err := p.Await()
	if err != nil {
		return Reject[T](err)
	}

	return Resolve(res)
}

func (p *Promise[T]) resolve(t T) {
	p.once.Do(func() {
		p.res = t
	})
}

func (p *Promise[T]) reject(err error) {
	p.once.Do(func() {
		p.err = err
	})
}

// All attempts to resolves all promises and rejects on the first error.
// The context is not propagated to the children promises.
func All[T any](ctx context.Context, promises ...*Promise[T]) *Promise[[]T] {
	return New(ctx, func(ctx context.Context) ([]T, error) {
		result := make([]T, len(promises))

		g, ctx := errgroup.WithContext(ctx)
		g.SetLimit(MaxGoroutine)

		for i, promise := range promises {
			i, promise := i, promise

			g.Go(func() error {
				res, err := promise.Await()
				if err != nil {
					return err
				}

				result[i] = res

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return nil, err
		}

		return result, nil
	})
}

// AllTasks attempts to resolves all tasks, and rejects on the first error.
// The context is passed to each task.
func AllTask[T any](ctx context.Context, tasks ...Task[T]) *Promise[[]T] {
	return New(ctx, func(ctx context.Context) ([]T, error) {
		result := make([]T, len(tasks))

		g, ctx := errgroup.WithContext(ctx)
		g.SetLimit(MaxGoroutine)

		for i, task := range tasks {
			i, task := i, task

			g.Go(func() (err error) {
				defer func() {
					if in := recover(); in != nil {
						if perr, ok := in.(error); ok {
							err = perr
						} else {
							err = fmt.Errorf("%w: %v", ErrPanic, in)
						}
					}
				}()

				res, err := task(ctx)
				if err != nil {
					return err
				}

				result[i] = res

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return nil, err
		}

		return result, nil
	})
}

func AllSettled[T any, K []*Result[T]](ctx context.Context, promises ...*Promise[T]) *Promise[K] {
	return New(ctx, func(ctx context.Context) (K, error) {
		var wg sync.WaitGroup
		wg.Add(len(promises))

		result := make([]*Result[T], len(promises))
		for i, promise := range promises {

			go func(i int, promise *Promise[T]) {
				defer wg.Done()

				res, err := promise.Await()
				if err != nil {
					result[i] = Reject[T](err)
				} else {
					result[i] = Resolve(res)
				}
			}(i, promise)
		}

		wg.Wait()

		return result, nil
	})
}

func AllTaskSettled[T any, K []*Result[T]](ctx context.Context, tasks ...Task[T]) *Promise[K] {
	return New(ctx, func(ctx context.Context) (K, error) {
		var wg sync.WaitGroup
		wg.Add(len(tasks))

		result := make([]*Result[T], len(tasks))
		for i, task := range tasks {
			go func(i int, task Task[T]) {
				defer wg.Done()

				defer func() {
					if in := recover(); in != nil {
						if err, ok := in.(error); ok {
							result[i] = Reject[T](err)
						} else {
							result[i] = Reject[T](fmt.Errorf("%w: %v", ErrPanic, in))
						}
					}
				}()

				res, err := task(ctx)
				if err != nil {
					result[i] = Reject[T](err)
				} else {
					result[i] = Resolve(res)
				}
			}(i, task)
		}

		wg.Wait()

		return result, nil
	})
}
