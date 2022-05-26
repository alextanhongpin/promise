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
		panic(fmt.Errorf("running task: %d", n))
	}

	p := promise.New(asyncTask)
	_ = promise.Catch(p, func(err error) {
		fmt.Println("got error", err)
	})
}
