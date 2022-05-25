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
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	p := promise.New(ctx, asyncTask)

	res, err := p.Await()
	fmt.Println(res, err)

	res, err = p.Await()
	fmt.Println(res, err)
}
