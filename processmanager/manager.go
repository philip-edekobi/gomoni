package processmanager

import (
	"os"
	"os/exec"
)

var arg0 string = "go"

// RunFile is a function that starts executing a file in a new process and
// returns a pointer to a exec.Cmd and an error
func Run(file, dirCtx string) (*os.Process, error) {
	args := []string{"", "run", file}

	path, err := exec.LookPath(arg0)
	if err != nil {
		return nil, err
	}

	args[0] = path

	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	procAttr.Dir = dirCtx

	proc, err := os.StartProcess(args[0], args, &procAttr)
	if err != nil {
		return nil, err
	}

	return proc, nil
}

func Kill(proc *os.Process, killCh <-chan int, exitCh chan<- int) error {
	<-killCh

	err := proc.Kill()
	if err != nil {
		return err
	}

	exitCh <- 1

	return nil
}
