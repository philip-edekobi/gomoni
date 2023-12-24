package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/philip-edekobi/gomoni/depmanager"
)

const mainFile = "main.go"

var proc *os.Process

func main() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("an error occured: %v", err)
	}

	if len(os.Args) > 1 {
		workDir, err = filepath.Abs(os.Args[1])
		if err != nil {
			log.Fatalf("an error occured: %v", err)
		}
	}

	err = depmanager.BuildGlobalDirMap(workDir)
	if err != nil {
		log.Fatalf("could not build global directory map: %v", err)
	}

	depmanager.BuildDeps(workDir)
	fmt.Println(depmanager.GlobalPkgMap)

	/*
		  go func() {
				proc, err = processmanager.Run(workDir + mainFile)

				if err != nil {
					panic(err)
				}
			}()
	*/
}
