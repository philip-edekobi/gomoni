package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/philip-edekobi/gomoni/depmanager"
)

const mainFile = "main.go"

func main() {
	workFolder, err := os.Getwd()
	if err != nil {
		log.Fatalf("an error occured: %v", err)
	}

	if len(os.Args) > 1 {
		workFolder, err = filepath.Abs(os.Args[1])
		if err != nil {
			log.Fatalf("an error occured: %v", err)
		}
	}

	workFile := workFolder + "/" + mainFile

	err = depmanager.BuildGlobalDirMap(workFolder)
	if err != nil {
		log.Fatalf("could not build global directory map: %v", err)
	}

	depmanager.BuildDeps(workFile)
}
