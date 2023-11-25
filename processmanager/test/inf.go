package main

import (
	"fmt"
)

func main() {
	fmt.Println("running")
	a := make(chan int)

	for {
		select {
		case <-a:
			fmt.Print("ahoy")
			break
		default:
		}
	}
}
