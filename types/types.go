package types

import "os"

type Monitor struct {
	MainFile string
	// ExitCh is a channel that controls the life of the program.
	// It waits for a value which when sent indicated that the program should end
	ExitCh         chan int
	CurrentProcess *os.Process
}
