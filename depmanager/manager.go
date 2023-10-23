package depmanager

import (
	"bufio"
	"io"
	"os"
	"strings"
)

const (
	initialFile = "../main.go"
	bufSize     = 1024
)

// GetDeps takes a file and returns a list of the dependencies of the
// file along with their dependencies
func GetDeps(file string) []*os.File {
	var deps []*os.File

	workFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer workFile.Close()

	err = buildDepPackages(workFile, deps)
	if err != nil {
		panic(err)
	}

	return deps
}

func buildDepPackages(file *os.File, dest []*os.File) error {
	//buf := make([]byte, bufSize)
	return nil
}

func extractImports(file *os.File) ([]string, error) {
	imports := []string{}

	fileReader := bufio.NewReader(file)
	for {
		line, err := fileReader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return imports, err
		}
		if err == io.EOF {
			break
		}

		trimmedLine := strings.TrimSpace(string(line))

		if strings.HasPrefix(trimmedLine, "package") {
			continue
		}

		if strings.HasPrefix(trimmedLine, "import") {
			if strings.HasSuffix(trimmedLine, "\"") {
				imports = append(imports, removeQuote(trimmedLine))
				break
			} else if strings.HasSuffix(trimmedLine, "(") {
				for strings.TrimSpace(string(line)) != ")" {
					line, err = fileReader.ReadBytes('\n')
					if err != nil && err != io.EOF {
						return imports, err
					}
					pkgImport := removeQuote(strings.TrimSpace(string(line)))

					if pkgImport != "" {
						imports = append(imports, pkgImport)
					}
				}

				break
			}
		}
	}

	return imports, nil
}

func removeQuote(line string) string {
	chars := []rune(line)
	quote := []rune{}
	inQuote := false

	for i := 0; i < len(chars); i++ {
		if chars[i] == '"' {
			inQuote = !inQuote
			continue
		}

		if inQuote {
			quote = append(quote, chars[i])
		} else {
			continue
		}
	}

	return string(quote)
}
