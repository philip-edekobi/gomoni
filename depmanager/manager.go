package depmanager

import (
	"bufio"
	"crypto/sha1"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	initialFile = "../main.go"
)

var (
	// GlobalPkgMap is the map of all the identified packages in the project dir
	GlobalPkgMap map[string]bool = make(map[string]bool)

	// GlobalFileHashMap is a map of files and the hashes of their imports.
	// It is used to determine when to rebuild the dependency array(i.e when
	// the imports change)
	GlobalFileHashMap map[string][]byte = make(map[string][]byte)
)

// BuildDeps takes a file and returns a list of the dependencies of the
// file along with their dependencies
func BuildDeps(file string) {
	workFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer workFile.Close()

	err = buildDepPackages(workFile)
	if err != nil {
		panic(err)
	}
}

func buildDepPackages(file *os.File) error {
	tempDepArr := []string{}

	fileDeps, err := extractImports(file)
	if err != nil {
		return err
	}

	GlobalFileHashMap[file.Name()] = hashFileDeps(fileDeps)

	for _, dep := range fileDeps {
		if !isValidDep(dep) {
			continue
		} else {
			if !GlobalPkgMap[dep] {
				GlobalPkgMap[dep] = true
				tempDepArr = append(tempDepArr, dep)
			} else {
				continue
			}
		}
	}

	for _, pkg := range tempDepArr {
		for _, file := range fetchFiles(pkg) {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			buildDepPackages(f)
		}
	}

	return nil
}

func fetchFiles(pkg string) []string {}

func hashFileDeps(fileDeps []string) []byte {
	hasher := sha1.New()

	for _, fileDep := range fileDeps {
		io.WriteString(hasher, fileDep)
	}

	return hasher.Sum(nil)
}

func isValidDep(dep string) bool {
	dep = strings.TrimSpace(dep)

	cmd := exec.Command("go", "list", "-m")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	modName := strings.TrimSpace(string(out))

	return strings.HasPrefix(dep, modName)
}

/*
 * another approach to validation that i might take in future
func isValidDep() {
	// construct recursive map of folders in main dir
	// match dep with the map
}
*/

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
