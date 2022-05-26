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
		fmt.Println("running task", n)
		time.Sleep(1 * time.Second)

		return n, nil
	}

	asyncThenTask := func(n int) *promise.Promise[string] {
		fmt.Println("running then", n)
		time.Sleep(1 * time.Second)

		return promise.Resolve(fmt.Sprintf("got number: %d", n))
	}

	p1 := promise.New(asyncTask)
	p2 := promise.Then(p1, asyncThenTask)

	res, err := p2.Await()

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
