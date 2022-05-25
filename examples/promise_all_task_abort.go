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
		// Each children must handle cancellation.
		n := rand.Intn(10)
		fmt.Println("running task", n)
		select {
		case <-time.After(1 * time.Second):
			return n, nil
		case <-ctx.Done():
			fmt.Println("done")
			return 0, ctx.Err()
		}
		//time.Sleep(1 * time.Second)
		//return n, nil
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	all := promise.AllTask(ctx, asyncTask, asyncTask)

	result, err := all.Await()
	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println(result)
}
