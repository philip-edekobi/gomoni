package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("file.test")
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintln(file, "running")
	file.Close()

	a := make(chan int)

	for {
		select {
		case <-a:
			return
		default:
			continue
		}
	}
}
