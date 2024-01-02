package processmanager

import (
	"fmt"
	"os"
	"os/exec"
)

var (
	arg0    string = "go"
	tmpFile string = "temp_prog_00000"
)

// RunFile is a function that starts executing a file in a new process and
// returns a pointer to a exec.Cmd and an error
func Run(file, dirCtx string) (*os.Process, error) {
	fmt.Println("[gomoni] - building")
	cmd := exec.Command("go", "build", "-o", tmpFile)
	cmd.Dir = dirCtx
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	fmt.Println("[gomoni] - running...")
	cmd = exec.Command("./" + tmpFile)
	cmd.Dir = dirCtx
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	// args := []string{"", "run", file}
	//
	// path, err := exec.LookPath(arg0)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// args[0] = path
	//
	// var procAttr os.ProcAttr
	// procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	// procAttr.Dir = dirCtx
	//
	// proc, err := os.StartProcess(args[0], args, &procAttr)
	// if err != nil {
	// 	return nil, err
	// }

	// return proc, nil
	return cmd.Process, nil
}

func Kill(proc *os.Process, killCh <-chan int) error {
	<-killCh

	err := proc.Kill()
	if err != nil {
		panic(err)
	}

	return nil
}

func WatchForEnd(proc *os.Process, dirCtx string) {
	procState, err := proc.Wait()
	if err != nil {
		panic(err)
	}

	if procState.Exited() {
		fmt.Println("[gomoni] - exit... waiting for changes to restart")
	}

	err = os.Remove(dirCtx + "/" + tmpFile)
	if err != nil {
		panic(err)
	}
}
