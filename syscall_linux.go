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

func parseXattrList(buf []byte) []string {
	return nullTermToStrings(buf)
}
