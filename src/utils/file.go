package utils

import (
	"runtime"
	"strings"
)

func CallOrigin(depth int) (string, int) {
	_, file, line, ok := runtime.Caller(depth)

	if ok {
		splitPath := strings.Split(file, "/")
		file = splitPath[len(splitPath)-1] // Get the last part of the path
	} else {
		file = "UnknownFile"
		line = 0
	}

	return file, line
}
