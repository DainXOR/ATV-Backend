package utils

import (
	"fmt"
	"runtime"
	"strings"
)

func CallOrigin(depth int) string {
	originFile := ""
	_, file, line, ok := runtime.Caller(depth)

	if ok {
		splitPath := strings.Split(file, "/")
		file = splitPath[len(splitPath)-1] // Get the last part of the path
		originFile = fmt.Sprintf("%s:%d", file, line)
	} else {
		originFile = "UnknownFile:0"
	}

	return originFile
}
