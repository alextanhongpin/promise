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
		panic(fmt.Errorf("panic worker: %d", n))
	}

	ctx := context.Background()
	p1 := promise.New(ctx, asyncTask)
	p2 := promise.New(ctx, asyncTask)

	result, err := promise.AllSettled(ctx, p1, p2).Await()
	if err != nil {
		panic(err)
	}

	for _, res := range result {
		fmt.Println(res.Unwrap())
	}
}