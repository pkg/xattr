// +build darwin

package xattr

import (
	"syscall"
	"unsafe"
)

func getxattr(path string, name string, data []byte) (int, error) {
	value, size := bytePtrFromSlice(data)
	/*
		ssize_t getxattr(
			const char *path,
			const char *name,
			void *value,
			size_t size,
			u_int32_t position,
			int options
		)
	*/
	r0, _, err := syscall.Syscall6(syscall.SYS_GETXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(name))), uintptr(unsafe.Pointer(value)), uintptr(size), 0, 0)
	if err != syscall.Errno(0) {
		return int(r0), err
	}
	return int(r0), nil
}

func setxattr(path string, name string, data []byte, flags int) error {
	value, size := bytePtrFromSlice(data)
	/*
		int setxattr(
			const char *path,
			const char *name,
			void *value,
			size_t size,
			u_int32_t position,
			int options
		);
	*/
	_, _, err := syscall.Syscall6(syscall.SYS_SETXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(name))), uintptr(unsafe.Pointer(value)), uintptr(size), 0, 0)
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func removexattr(path string, name string) error {
	/*
		int removexattr(
			const char *path,
			const char *name,
			int options
		);
	*/
	_, _, err := syscall.Syscall(syscall.SYS_REMOVEXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(name))), 0)
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func listxattr(path string, data []byte) (int, error) {
	name, size := bytePtrFromSlice(data)
	/*
		ssize_t listxattr(
			const char *path,
			char *name,
			size_t size,
			int options
		);
	*/
	r0, _, err := syscall.Syscall6(syscall.SYS_LISTXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		uintptr(unsafe.Pointer(name)), uintptr(size), 0, 0, 0)
	if err != syscall.Errno(0) {
		return int(r0), err
	}
	return int(r0), nil
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
