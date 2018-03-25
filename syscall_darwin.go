// +build darwin

package xattr

import (
	"syscall"
	"unsafe"
)

func getxattr(path string, attr string, dest []byte) (sz int, err error) {
	size := len(dest)
	var value *byte
	if size > 0 {
		value = &dest[0]
	}
	return darwinGetxattr(path, attr, value, size, 0, 0)
}

func setxattr(path string, attr string, data []byte, flags int) (err error) {
	size := len(data)
	var value *byte
	if size > 0 {
		value = &data[0]
	}
	return darwinSetxattr(path, attr, value, size, 0, flags)
}

func removexattr(path string, attr string) (err error) {
	return darwinRemovexattr(path, attr, 0)
}

func listxattr(path string, dest []byte) (sz int, err error) {
	size := len(dest)
	var namebuf *byte
	if size > 0 {
		namebuf = &dest[0]
	}
	return darwinListxattr(path, namebuf, size, 0)
}

func parseXattrList(buf []byte) []string {
	return nullTermToStrings(buf)
}

/*
	ssize_t
	getxattr(const char *path, const char *name, void *value, size_t size,
	         u_int32_t position, int options);

	int
	setxattr(const char *path, const char *name, void *value, size_t size,
	         u_int32_t position, int options);

	ssize_t
	listxattr(const char *path, char *namebuf, size_t size, int options);

	int
	removexattr(const char *path, const char *name, int options);
*/
func darwinGetxattr(path string, name string, value *byte, size int, pos int, options int) (int, error) {
	r0, _, e1 := syscall.Syscall6(syscall.SYS_GETXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))), uintptr(unsafe.Pointer(syscall.StringBytePtr(name))), uintptr(unsafe.Pointer(value)), uintptr(size), uintptr(pos), uintptr(options))
	if e1 != syscall.Errno(0) {
		return int(r0), e1
	}
	return int(r0), nil
}

func darwinSetxattr(path string, name string, value *byte, size int, pos int, options int) error {
	if _, _, e1 := syscall.Syscall6(syscall.SYS_SETXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))), uintptr(unsafe.Pointer(syscall.StringBytePtr(name))), uintptr(unsafe.Pointer(value)), uintptr(size), uintptr(pos), uintptr(options)); e1 != syscall.Errno(0) {
		return e1
	}
	return nil
}

func darwinRemovexattr(path string, name string, options int) error {
	if _, _, e1 := syscall.Syscall(syscall.SYS_REMOVEXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))), uintptr(unsafe.Pointer(syscall.StringBytePtr(name))), uintptr(options)); e1 != syscall.Errno(0) {
		return e1
	}
	return nil
}

func darwinListxattr(path string, namebuf *byte, size int, options int) (int, error) {
	r0, _, e1 := syscall.Syscall6(syscall.SYS_LISTXATTR, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))), uintptr(unsafe.Pointer(namebuf)), uintptr(size), uintptr(options), 0, 0)
	if e1 != syscall.Errno(0) {
		return int(r0), e1
	}
	return int(r0), nil
}
