package main

import (
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

	asyncTask := func() (int, error) {
		n := rand.Intn(10)
		panic(fmt.Errorf("panic worker: %d", n))
	}

	p1 := promise.New(asyncTask)
	p2 := promise.New(asyncTask)

	result, err := promise.AllSettled(p1, p2).Await()
	if err != nil {
		panic(err)
	}

	for _, res := range result {
		fmt.Println(res.Unwrap())
	}
}
