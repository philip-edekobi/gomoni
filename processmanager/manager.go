package processmanager

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/philip-edekobi/gomoni/types"
)

var (
	arg0    string = "go"
	tmpFile string = "temp_prog_00000"

	// KillCh is a channel that waits for a signal to kill a process
	KillCh chan int = make(chan int, 1)
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
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
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

func Kill(proc *os.Process) error {
	<-KillCh

	// err := proc.Kill()
	err := syscall.Kill(-proc.Pid, syscall.SIGKILL)
	if err != nil {
		panic(err)
	}

	return nil
}

func Restart(monitor *types.Monitor, file, dirCtx string) error {
	p, err := Run(file, dirCtx)
	if err != nil {
		panic(err)
	}

	monitor.CurrentProcess = p
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
