package utils

import (
	"path/filepath"
	"runtime"
)

func CallOrigin(depth int) (string, uint64) {
	_, file, line, ok := runtime.Caller(depth)

	if ok {
		_, filePath := filepath.Split(file)
		file = filePath
	} else {
		file = "UnknownFile"
		line = 0
	}

	return file, uint64(line)
}
