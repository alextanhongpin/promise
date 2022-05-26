package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alextanhongpin/promise"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	asyncTask := func() (int, error) {
		panic("something happened")
	}

	res, err := promise.New(asyncTask).Await()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
