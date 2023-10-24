package depmanager

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildDeps(t *testing.T) {
	BuildDeps("test_file1.go")
}

func TestExtractImports(t *testing.T) {
	testCases := []struct {
		FileName string
		Imports  []string
	}{
		{"test_file1.go", []string{"fmt"}},
		{"test_file2.go", []string{"fmt", "net/http", "time"}},
		{"test_file3.go", []string{}},
	}

	for _, tc := range testCases {

		f, err := os.Open(tc.FileName)
		require.Nil(t, err)

		imports, err := extractImports(f)
		require.Nil(t, err)
		require.Equal(t, tc.Imports, imports)
	}
}

func TestRemoveQuote(t *testing.T) {
	testCases := []struct {
		line  string
		quote string
	}{
		{"import \"fmt\"", "fmt"},
		{"impo\"box\"", "box"},
		{"the man is here", ""},
	}

	for _, tc := range testCases {
		quote := removeQuote(tc.line)

		require.Equal(t, tc.quote, quote)
	}
}

func TestIsValidDep(t *testing.T) {
	testCases := []struct {
		dep     string
		isValid bool
	}{
		{"github.com/philip-edekobi/gomoni/depmanager", true},
		{"fmt", false},
		{"", false},
	}

	for _, tc := range testCases {
		//fmt.Println(tc.dep)
		require.Equal(t, tc.isValid, isValidDep(tc.dep))
	}
}
