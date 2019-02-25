package main

import (
	"strconv"
)

func MakePathsFromDirectory(directory, filename, fileextension string, start, end int) []string {
	var paths []string
	for i := start; i <= end; i++ {
		path := directory + filename + strconv.Itoa(i) + "." + fileextension
		paths = append(paths, path)
	}
	return paths
}
