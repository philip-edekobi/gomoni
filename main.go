package main

import (
	"fmt"
	"os"
)

const mainFile = "main.go"

func main() {
	var initFolder string = "."

	if len(os.Args) > 1 {
		initFolder = os.Args[1]
	}

	workFile := initFolder + "/" + mainFile

	fmt.Println(workFile)
}
