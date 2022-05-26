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

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	asyncTask := func() (int, error) {
		n := rand.Intn(10)
		fmt.Println("running task", n)
		// Each children must handle cancellation.
		select {
		case <-time.After(1 * time.Second):
			return n, nil
		case <-ctx.Done():
			fmt.Println("aborting")
			return 0, ctx.Err()
		}
	}

	// All child must be passed in the context explicitly for cancellation.
	p1 := promise.New(asyncTask)
	p2 := promise.New(asyncTask)

	all := promise.AllSettled(p1, p2)
	result, err := all.Await()
	if err != nil {
		fmt.Println("err", err)
	}

	for _, res := range result {
		fmt.Println(res.Unwrap())
	}
}
