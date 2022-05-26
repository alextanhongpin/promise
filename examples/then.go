package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/alextanhongpin/promise"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	asyncTask := func(ctx context.Context) (int, error) {
		n := rand.Intn(10)
		fmt.Println("running task", n)
		time.Sleep(1 * time.Second)

		return n, nil
	}

	asyncThenTask := func(ctx context.Context, n int) (string, error) {
		fmt.Println("running then", n)
		time.Sleep(1 * time.Second)

		return fmt.Sprintf("got number: %d", n), nil
	}

	ctx := context.Background()
	p1 := promise.New(ctx, asyncTask)
	p2 := promise.Then(ctx, p1, asyncThenTask)

	res, err := p2.Await()

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
