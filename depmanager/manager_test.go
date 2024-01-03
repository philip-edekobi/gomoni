package depmanager

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGlobalDirMap(t *testing.T) {
	expectedDirMap := make(map[string]bool)
	expectedDirMap["/home/luxurydev/Desktop/work/projects/gomoni/depmanager/test_proj"] = true
	expectedDirMap["/home/luxurydev/Desktop/work/projects/gomoni/depmanager/test_proj/calc"] = true
	expectedDirMap["/home/luxurydev/Desktop/work/projects/gomoni/depmanager/test_proj/out"] = true
	expectedDirMap["/home/luxurydev/Desktop/work/projects/gomoni/depmanager/test_proj/calc/mul"] = true

	BuildGlobalDirMap("/home/luxurydev/Desktop/work/projects/gomoni/depmanager/test_proj")

	require.Equal(t, expectedDirMap, GlobalDirMap)
}

func TestBuildDeps(t *testing.T) {
	expectedPkgMap := make(map[string]string)
	expectedPkgMap["test_proj/calc"] = "/home/luxurydev/Desktop/work/projects/gomoni/depmanager/test_proj/calc"
	expectedPkgMap["test_proj/calc/mul"] = "/home/luxurydev/Desktop/work/projects/gomoni/depmanager/test_proj/calc/mul"
	expectedPkgMap["test_proj/out"] = "/home/luxurydev/Desktop/work/projects/gomoni/depmanager/test_proj/out"

	BuildDeps("/home/luxurydev/Desktop/work/projects/gomoni/depmanager/test_proj")

	require.Equal(t, expectedPkgMap, GlobalPkgMap)
}

func TestExtractImports(t *testing.T) {
	testCases := []struct {
		FileName string
		Imports  []string
	}{
		{"test_file1.go", []string{"fmt"}},
		{"test_file2.go", []string{"fmt", "net", "net/http", "time"}},
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
		// fmt.Println(tc.dep)
		require.Equal(t, tc.isValid, isValidDep(tc.dep, "."))
	}
}

func TestFetchFiles(t *testing.T) {
	testCases := []struct {
		dir           string
		expectedFiles []string
	}{
		{"emp", []string{}},
		{
			"test_proj",
			[]string{
				"test_proj/calc/calc.go",
				"test_proj/calc/mul/mul.go",
				"test_proj/go.mod",
				"test_proj/main.go",
				"test_proj/out/out.go",
			},
		},
	}

	for _, testCase := range testCases {
		files, err := FetchFiles(testCase.dir)

		require.Nil(t, err)
		require.Equal(t, testCase.expectedFiles, files)
	}
}

func TestEmptyPkgMap(t *testing.T) {
	EmptyPkgMap()

	require.Empty(t, GlobalPkgMap)
}
