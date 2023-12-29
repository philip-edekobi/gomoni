package processmanager

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testLoc string = "/home/luxurydev/Desktop/work/projects/gomoni/processmanager/test/"

func TestRun(t *testing.T) {
	testCases := []struct {
		expectedOutput string
		testFile       string
		isInfinite     bool
		dir            string
	}{
		{
			expectedOutput: "Hello World",
			testFile:       testLoc + "main.go",
			dir:            testLoc,
		},
		{
			isInfinite: true,
			testFile:   testLoc + "inf.go",
			dir:        testLoc,
		},
	}

	for _, testCase := range testCases {
		proc, err := Run(testCase.testFile, testCase.dir)
		require.Nil(t, err)

		if !testCase.isInfinite {
			err = proc.Kill()
			require.Nil(t, err)
		} else {
			time.Sleep(time.Duration(3) * time.Second)

			err = proc.Kill()
			require.Nil(t, err)

			file := testLoc + "file.test"
			lines := []string{}

			f, err := os.Open(file)
			require.Nil(t, err)
			defer f.Close()

			fScanner := bufio.NewScanner(f)
			for fScanner.Scan() {
				lines = append(lines, fScanner.Text())
			}

			text := strings.Join(lines, "")

			require.Equal(t, "running", strings.TrimSpace(text))
		}
	}
}

func TestKill(t *testing.T) {
	args := []string{""}

	path, err := exec.LookPath("cat")
	require.Nil(t, err)

	args[0] = path

	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	proc, err := os.StartProcess(args[0], args, &procAttr)
	require.Nil(t, err)

	k := make(chan int, 1)
	e := make(chan int, 1)

	k <- 1

	Kill(proc, k, e)

	val := <-e

	require.Equal(t, 1, val)
}
