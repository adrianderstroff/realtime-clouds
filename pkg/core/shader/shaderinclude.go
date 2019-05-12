package shader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func loadFileWithIncludesRecursive(filename string, loadedpaths *map[string]bool) string {
	// simplify filepath to consistently find duplicate includes
	cleanpath := filepath.Clean(filename)

	// early return since file had been loaded before
	if _, ok := (*loadedpaths)[cleanpath]; ok {
		return ""
	}
	(*loadedpaths)[cleanpath] = true

	// create a file reader
	file, err := os.Open(cleanpath)
	defer file.Close()
	if err != nil {
		fmt.Println("Shaderinclude cannot find " + cleanpath)
		return ""
	}
	reader := bufio.NewReader(file)

	// regular expression for matching the relative paths
	exp := regexp.MustCompile(`"([^"]+)"`)

	// read the file line by line and load includes
	var shadersource string
	for {
		line, err := reader.ReadString('\n')

		if matched, _ := regexp.MatchString(`#include`, line); matched {
			// extract the relative file path
			relativepath := exp.FindString(line)
			relativepath = strings.Trim(relativepath, "\"")
			basepath, _ := filepath.Split(cleanpath)

			// load new shader include file recursively and attach
			// the result to shadersource
			newpath := basepath + relativepath
			childshadersource := loadFileWithIncludesRecursive(newpath, loadedpaths)
			shadersource += childshadersource
		} else {
			shadersource = shadersource + line
		}

		// an error is returned when the last line had been read
		if err != nil {
			shadersource += "\n"
			break
		}
	}

	return shadersource
}

func LoadFileWithIncludes(filepath string) (string, error) {
	loadedpaths := make(map[string]bool)
	return loadFileWithIncludesRecursive(filepath, &loadedpaths) + "\000", nil
}
