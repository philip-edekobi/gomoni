package processmanager

import (
	"bufio"
	"os"
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
	}{
		{
			expectedOutput: "Hello World",
			testFile:       testLoc + "main.go",
		},
		{
			isInfinite: true,
			testFile:   testLoc + "inf.go",
		},
	}

	for _, testCase := range testCases {
		proc, err := Run(testCase.testFile)
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

			err = os.Remove("file.test")
			require.Nil(t, err)
		}
	}
}
