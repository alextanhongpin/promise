package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alextanhongpin/promise"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	asyncTask := func(ctx context.Context) (int, error) {
		panic("something happened")
	}

	ctx := context.Background()
	res, err := promise.New(ctx, asyncTask).Await()
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
