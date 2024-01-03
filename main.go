package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/philip-edekobi/gomoni/depmanager"
	"github.com/philip-edekobi/gomoni/filewatcher"
	"github.com/philip-edekobi/gomoni/processmanager"
	"github.com/philip-edekobi/gomoni/types"
)

var monitor = &types.Monitor{
	MainFile: "main.go",
	ExitCh:   make(chan int, 1),
}

var sigs = make(chan os.Signal, 1)

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

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go handleTerminate()

	err = depmanager.BuildGlobalDirMap(workDir)
	if err != nil {
		log.Fatalf("could not build global directory map: %v", err)
	}

	depmanager.BuildDeps(workDir)
	depmanager.GlobalPkgMap["main"] = workDir

	filewatcher.Initialize(monitor, workDir)
	go filewatcher.WatchFiles(depmanager.GlobalPkgMap)

	fmt.Println("[gomoni] - Starting...")
	monitor.CurrentProcess, err = processmanager.Run(workDir+"/"+monitor.MainFile, workDir)
	if err != nil {
		panic(err)
	}

	go processmanager.WatchForEnd(monitor.CurrentProcess, workDir)
	go processmanager.Kill(monitor.CurrentProcess)

	// t := time.After(4 * time.Second)
	// <-t
	// processmanager.KillCh <- 1
	<-monitor.ExitCh
}

func handleTerminate() {
	<-sigs

	processmanager.KillCh <- 1
	monitor.ExitCh <- 1
}
