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
	res, err := promise.AllTask(ctx, asyncTask, asyncTask).Await()
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
