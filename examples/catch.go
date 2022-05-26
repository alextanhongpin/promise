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
		panic(fmt.Errorf("running task: %d", n))
	}

	ctx := context.Background()
	p := promise.New(ctx, asyncTask)
	_ = promise.Catch(p, func(err error) {
		fmt.Println("got error", err)
	})
}
