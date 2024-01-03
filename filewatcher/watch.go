package filewatcher

import (
	"os"

	"github.com/fsnotify/fsnotify"

	"github.com/philip-edekobi/gomoni/depmanager"
	"github.com/philip-edekobi/gomoni/processmanager"
)

var (
	globalWatcher *fsnotify.Watcher = nil
	dirCtx        string
	process       *os.Process
)

func Initialize(proc *os.Process, dir string) {
	dirCtx = dir

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	globalWatcher = watcher
	process = proc
}

func Close() {
	globalWatcher.Close()
}

func attachWatchers(fileMap map[string]string) {
	for _, dir := range fileMap {
		err := globalWatcher.Add(dir)
		if err != nil {
			panic(err)
		}
	}
}

func removeWatchers(fileMap map[string]string) {
	for _, dir := range fileMap {
		err := globalWatcher.Remove(dir)
		if err != nil {
			panic(err)
		}
	}
}

func WatchFiles(fileMap map[string]string) {
	attachWatchers(fileMap)

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
				// fmt.Println("CHANGEEEEEE")
				removeWatchers(depmanager.GlobalPkgMap)
				depmanager.BuildDeps(dirCtx)
				attachWatchers(depmanager.GlobalPkgMap)
				processmanager.KillCh <- 1
				processmanager.Restart(process, "main.go", dirCtx)

				go processmanager.WatchForEnd(process, dirCtx)
				go processmanager.Kill(process)
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
	for fileName := range depmanager.GlobalFileHashMap {
		if file == fileName {
			return true
		}
	}
	return false
}
