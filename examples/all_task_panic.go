package main

import (
	"context"
	"fmt"
	"log"
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
		panic(fmt.Errorf("intended panic: %d", n))
	}

	ctx := context.Background()
	res, err := promise.AllTask(ctx, asyncTask, asyncTask).Await()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
