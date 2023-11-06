package main

import (
	"log"
	"os"
	"path/filepath"
)

const mainFile = "main.go"

func main() {
	initFolder, err := os.Getwd()
	if err != nil {
		log.Fatalf("an error occured: %v", err)
	}

	if len(os.Args) > 1 {
		initFolder, err = filepath.Abs(os.Args[1])
		if err != nil {
			log.Fatalf("an error occured: %v", err)
		}
	}

	workFile := initFolder + "/" + mainFile

	log.Println(workFile)
}
