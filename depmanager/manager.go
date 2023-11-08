package depmanager

import (
	"bufio"
	"crypto/sha1"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GlobalDirMap is the map of all the subdirectories in the project
var GlobalDirMap map[string]bool = make(map[string]bool)

// GlobalPkgMap is the map of all the identified packages in the project dir
// to their file paths
var GlobalPkgMap map[string]string = make(map[string]string)

// GlobalFileHashMap is a map of files and the hashes of their imports.
// It is used to determine when to rebuild the dependency array(i.e when
// the imports change)
var GlobalFileHashMap map[string][]byte = make(map[string][]byte)

// BuildGlobalDirMap is a function that constructs a map that contains all
// the files and directories in the currDir argument passed to it
func BuildGlobalDirMap(currDir string) error {
	err := filepath.WalkDir(currDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			GlobalDirMap[path] = true
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// BuildDeps takes a file and returns a list of the dependencies of the
// file along with their dependencies
func BuildDeps(dir string) {
	mainFile := dir + "/main.go"
	workFile, err := os.Open(mainFile)
	if err != nil {
		panic(err)
	}
	defer workFile.Close()

	err = buildDepPackages(workFile, dir)
	if err != nil {
		panic(err)
	}
}

func buildDepPackages(file *os.File, dirCtx string) error {
	tempDepArr := []string{}

	fileDeps, err := extractImports(file)
	if err != nil {
		return err
	}

	GlobalFileHashMap[file.Name()] = hashFileDeps(fileDeps)

	for _, dep := range fileDeps {
		if !isValidDep(dep, dirCtx) {
			continue
		} else {
			depPath := findPath(dep)

			_, ok := GlobalPkgMap[dep]
			if !ok {
				GlobalPkgMap[dep] = depPath
				tempDepArr = append(tempDepArr, depPath)
			} else {
				continue
			}
		}
	}

	for _, pkg := range tempDepArr {
		files, err := fetchFiles(pkg)
		if err != nil {
			return err
		}

		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			buildDepPackages(f, dirCtx)
		}
	}

	return nil
}

func findPath(dep string) string {
	for k := range GlobalDirMap {
		if strings.HasSuffix(k, dep) {
			return k
		}
	}

	return ""
}

func fetchFiles(pkg string) ([]string, error) {
	files := []string{}

	err := filepath.WalkDir(pkg, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return files, err
	}

	return files, nil
}

func hashFileDeps(fileDeps []string) []byte {
	hasher := sha1.New()

	for _, fileDep := range fileDeps {
		io.WriteString(hasher, fileDep)
	}

	return hasher.Sum(nil)
}

func isValidDep(dep, dirCtx string) bool {
	dep = strings.TrimSpace(dep)

	cmd := exec.Command("go", "list", "-m")
	cmd.Dir = dirCtx
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
