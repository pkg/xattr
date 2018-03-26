// +build darwin

package xattr

import (
	"syscall"
	"unsafe"
)

func getxattr(path string, attr string, dest []byte) (sz int, err error) {
	value, size := sliceToPtr(dest)
	/*
		ssize_t getxattr(const char *path, const char *name,
			void *value, size_t size, u_int32_t position, int options)
	*/
	r0, _, e1 := syscall.Syscall6(syscall.SYS_GETXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(attr))), uintptr(unsafe.Pointer(value)), uintptr(size), 0, 0)
	if e1 != 0 {
		return int(r0), e1
	}
	return int(r0), nil
}

func setxattr(path string, attr string, data []byte, flags int) (err error) {
	value, size := sliceToPtr(data)
	/*
		int setxattr(const char *path,
			const char *name, void *value, size_t size, u_int32_t position, int options);
	*/
	_, _, e1 := syscall.Syscall6(syscall.SYS_SETXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(attr))), uintptr(unsafe.Pointer(value)), uintptr(size), 0, 0)
	if e1 != 0 {
		return e1
	}
	return
}

func removexattr(path string, attr string) (err error) {
	/*
		int removexattr(const char *path,
			const char *name, int options);
	*/
	_, _, e1 := syscall.Syscall(syscall.SYS_REMOVEXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(attr))), 0)
	if e1 != 0 {
		return e1
	}
	return
}

func listxattr(path string, dest []byte) (sz int, err error) {
	namebuf, size := sliceToPtr(dest)
	/*
		ssize_t listxattr(const char *path, char *namebuf,
			size_t size, int options);
	*/
	r0, _, e1 := syscall.Syscall6(syscall.SYS_LISTXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		uintptr(unsafe.Pointer(namebuf)), uintptr(size), 0, 0, 0)
	if e1 != 0 {
		return int(r0), e1
	}
	return int(r0), nil
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
