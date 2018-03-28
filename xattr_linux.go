// +build linux

package xattr

import (
	"syscall"
)

func getxattr(path string, attr string, dest []byte) (sz int, err error) {
	return syscall.Getxattr(path, attr, dest)
}

func setxattr(path string, attr string, data []byte, flags int) (err error) {
	return syscall.Setxattr(path, attr, data, flags)
}

func removexattr(path string, attr string) (err error) {
	return syscall.Removexattr(path, attr)
}

func listxattr(path string, dest []byte) (sz int, err error) {
	return syscall.Listxattr(path, dest)
}

// attrListToStrings converts a sequence of attribute name entries to a
// []string.
// On Darwin and Linux, each entry is a NUL-terminated string.
func attrListToStrings(buf []byte) (result []string) {
	offset := 0
	for index, b := range buf {
		if b == 0 {
			result = append(result, string(buf[offset:index]))
			offset = index + 1
		}
	}
	return
}
