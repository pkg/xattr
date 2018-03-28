// +build linux

package xattr

import (
	"syscall"
)

func getxattr(path string, name string, data []byte) (int, error) {
	return syscall.Getxattr(path, name, data)
}

func setxattr(path string, name string, data []byte, flags int) error {
	return syscall.Setxattr(path, name, data, flags)
}

func removexattr(path string, name string) error {
	return syscall.Removexattr(path, name)
}

func listxattr(path string, data []byte) (int, error) {
	return syscall.Listxattr(path, data)
}

// stringsFromByteSlice converts a sequence of attributes to a []string.
// On Darwin and Linux, each entry is a NULL-terminated string.
func stringsFromByteSlice(buf []byte) (result []string) {
	offset := 0
	for index, b := range buf {
		if b == 0 {
			result = append(result, string(buf[offset:index]))
			offset = index + 1
		}
	}
	return
}
