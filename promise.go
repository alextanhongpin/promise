package promise

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

const MaxGoroutine = 1_000

type Promise[T any] struct {
	res    T
	err    error
	status Status
	wg     sync.WaitGroup
	mu     sync.RWMutex
	ctx    context.Context
	once   sync.Once
}

type Task[T any] func(ctx context.Context) (T, error)

func New[T any](ctx context.Context, fn Task[T]) *Promise[T] {
	p := &Promise[T]{
		status: Pending,
		ctx:    ctx,
	}
	p.wg.Add(2)
	ch := make(chan Result[T], 1)

	go func() {
		defer p.wg.Done()

		res, err := fn(ctx)
		ch <- Result[T]{res: res, err: err, dirty: true}
	}()

	go func() {
		defer p.wg.Done()

		select {
		case <-ctx.Done():
			p.reject(ctx.Err())
		case result := <-ch:
			res, err := result.Unwrap()
			if err != nil {
				p.reject(err)
			} else {
				p.resolve(res)
			}
		}
	}()

	return p
}

func (p *Promise[T]) Await() (T, error) {
	p.wg.Wait()
	return p.res, p.err
}

func (p *Promise[T]) Then(fn func(context.Context, T) (T, error)) *Promise[T] {
	p.wg.Wait()
	if p.Status() == Rejected {
		return p
	}

	return New(p.ctx, func(ctx context.Context) (t T, err error) {
		return fn(ctx, p.res)
	})
}

func (p *Promise[T]) Catch(fn func(err error)) *Promise[T] {
	p.wg.Wait()
	if p.Status() == Rejected {
		fn(p.err)
	}
	return p
}

func (p *Promise[T]) resolve(t T) {
	p.once.Do(func() {
		p.res = t
		p.status = Fulfilled
	})
}

func (p *Promise[T]) reject(err error) {
	p.once.Do(func() {
		p.err = err
		p.status = Rejected
	})
}

func (p *Promise[T]) Status() Status {
	p.mu.RLock()
	status := p.status
	p.mu.RUnlock()

	return status
}

func All[T any](ctx context.Context, promises ...*Promise[T]) *Promise[[]T] {
	return New[[]T](ctx, func(ctx context.Context) ([]T, error) {
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

func AllTask[T any](ctx context.Context, tasks ...Task[T]) *Promise[[]T] {
	return New[[]T](ctx, func(ctx context.Context) ([]T, error) {
		result := make([]T, len(tasks))

		g, ctx := errgroup.WithContext(ctx)
		g.SetLimit(MaxGoroutine)

		for i, task := range tasks {
			i, task := i, task

			g.Go(func() error {
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
	return New[K](ctx, func(ctx context.Context) (K, error) {
		var wg sync.WaitGroup
		wg.Add(len(promises))

		result := make([]*Result[T], len(promises))
		for i, promise := range promises {
			i, promise := i, promise

			res, err := promise.Await()
			if err != nil {
				result[i] = Reject[T](err)
			} else {
				result[i] = Resolve(res)
			}
		}

		return result, nil
	})
}

func AllTaskSettled[T any, K []*Result[T]](ctx context.Context, tasks ...Task[T]) *Promise[K] {
	return New[K](ctx, func(ctx context.Context) (K, error) {
		var wg sync.WaitGroup
		wg.Add(len(tasks))

		result := make([]*Result[T], len(tasks))
		for i, task := range tasks {
			i, task := i, task

			res, err := task(ctx)
			if err != nil {
				result[i] = Reject[T](err)
			} else {
				result[i] = Resolve(res)
			}
		}

		return result, nil
	})
}
