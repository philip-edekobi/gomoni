package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/philip-edekobi/gomoni/depmanager"
	"github.com/philip-edekobi/gomoni/filewatcher"
	"github.com/philip-edekobi/gomoni/processmanager"
)

const mainFile = "main.go"

// ExitCh is a channel that controls the life of the program.
// It waits for a value which when sent indicated that the program should end
var ExitCh = make(chan int, 1)

// KillCh is a channel that waits for a signal to kill a process
var KillCh = make(chan int, 1)

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
	depmanager.GlobalPkgMap["main"] = workDir

	filewatcher.Initialize()
	go filewatcher.WatchFiles(depmanager.GlobalPkgMap)

	fmt.Println("[gomoni] - Starting...")
	proc, err := processmanager.Run(workDir+"/"+mainFile, workDir)
	if err != nil {
		panic(err)
	}

	go processmanager.WatchForEnd(proc, workDir)
	go processmanager.Kill(proc, KillCh)

	t := time.After(4 * time.Second)
	<-t

	// KillCh <- 1
	<-ExitCh
}
