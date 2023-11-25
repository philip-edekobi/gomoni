package processmanager

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testLoc string = "/home/luxurydev/Desktop/work/projects/gomoni/processmanager/test/"

func TestRun(t *testing.T) {
	testCase := struct {
		expectedOutput string
		testFile       string
	}{
		expectedOutput: "Hello World",
		testFile:       testLoc + "main.go",
	}

	proc, err := Run(testCase.testFile)
	require.Nil(t, err)

	err = proc.Release()
	require.Nil(t, err)
}
