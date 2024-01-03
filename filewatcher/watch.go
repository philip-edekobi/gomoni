package filewatcher

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var globalWatcher *fsnotify.Watcher = nil

func Initialize() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	globalWatcher = watcher
}

func Close() {
	globalWatcher.Close()
}

func WatchFiles(fileMap map[string]string) {
	for _, dir := range fileMap {
		err := globalWatcher.Add(dir)
		if err != nil {
			panic(err)
		}
	}

	for {
		select {
		case event, ok := <-globalWatcher.Events:
			if !ok {
				return
			}

			if !fileIsValid(event.Name) {
				continue
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("CHANGEEEEEE")
			}
		case err, ok := <-globalWatcher.Errors:
			if !ok {
				return
			}

			print("ERROR:", err)
		}
	}
}

func fileIsValid(file string) bool {
	return strings.HasSuffix(file, ".go") && !strings.HasSuffix(file, "_test.go")
}
