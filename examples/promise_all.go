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

	ctx := context.Background()
	p1 := promise.New(ctx, asyncTask)
	p2 := promise.New(ctx, asyncTask)

	res, err := promise.All(ctx, p1, p2).Await()
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
