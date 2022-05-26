package promise

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

const MaxGoroutine = 1_000

var ErrPanic = errors.New("panic")

type Promise[T any] struct {
	res  T
	err  error
	wg   sync.WaitGroup
	once sync.Once
}

type Task[T any] func() (T, error)

func TaskContext[T any](ctx context.Context, fn func(ctx context.Context) (T, error)) Task[T] {
	return func() (T, error) {
		return fn(ctx)
	}
}

func New[T any](fn Task[T]) *Promise[T] {
	p := &Promise[T]{}
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				if perr, ok := err.(error); ok {
					p.reject(perr)
				} else {
					p.reject(fmt.Errorf("%w: %v", ErrPanic, err))
				}
			}
		}()

		res, err := fn()
		if err != nil {
			p.reject(err)
		} else {
			p.resolve(res)
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
		return reject[T](err)
	}

	return resolve(res)
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

func AllSettled[T any, K []*Result[T]](promises ...*Promise[T]) *Promise[K] {
	return New(func() (K, error) {
		var wg sync.WaitGroup
		wg.Add(len(promises))

		result := make([]*Result[T], len(promises))
		for i, promise := range promises {

			go func(i int, promise *Promise[T]) {
				defer wg.Done()

				res, err := promise.Await()
				if err != nil {
					result[i] = reject[T](err)
				} else {
					result[i] = resolve(res)
				}
			}(i, promise)
		}

		wg.Wait()

		return result, nil
	})
}

func AllTaskSettled[T any, K []*Result[T]](tasks ...Task[T]) *Promise[K] {
	return New(func() (K, error) {
		var wg sync.WaitGroup
		wg.Add(len(tasks))

		result := make([]*Result[T], len(tasks))
		for i, task := range tasks {
			go func(i int, task Task[T]) {
				defer wg.Done()

				defer func() {
					if in := recover(); in != nil {
						if err, ok := in.(error); ok {
							result[i] = reject[T](err)
						} else {
							result[i] = reject[T](fmt.Errorf("%w: %v", ErrPanic, in))
						}
					}
				}()

				res, err := task()
				if err != nil {
					result[i] = reject[T](err)
				} else {
					result[i] = resolve(res)
				}
			}(i, task)
		}

		wg.Wait()

		return result, nil
	})
}
