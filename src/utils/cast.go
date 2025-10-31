package utils

import "unsafe"

// https://dev.to/chigbeef_77/bool-int-but-stupid-in-go-3jb3
func BoolToInt(value bool) int {
	return int(*(*byte)(unsafe.Pointer(&value)))
}
