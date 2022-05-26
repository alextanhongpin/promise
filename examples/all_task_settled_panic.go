package main

import (
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

	asyncTask := func() (int, error) {
		n := rand.Intn(10)
		panic(fmt.Errorf("intended panic: %d", n))
	}

	res, err := promise.AllTaskSettled(asyncTask, asyncTask).Await()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
