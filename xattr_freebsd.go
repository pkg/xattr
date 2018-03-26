// +build freebsd

package xattr

import (
	"syscall"
	"unsafe"
)

const (
	EXTATTR_NAMESPACE_USER = 1
)

func getxattr(path string, attr string, dest []byte) (sz int, err error) {
	data, nbytes := sliceToPtr(dest)
	/*
		ssize_t extattr_get_file(const char *path,
			int attrnamespace, const char *attrname,
			void *data, size_t nbytes);
	*/
	r0, _, e1 := syscall.Syscall6(syscall.SYS_EXTATTR_GET_FILE, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		EXTATTR_NAMESPACE_USER, uintptr(unsafe.Pointer(syscall.StringBytePtr(attr))),
		uintptr(unsafe.Pointer(data)), uintptr(nbytes), 0)
	if e1 != 0 {
		return int(r0), err
	}
	return int(r0), nil
}

func setxattr(path string, attr string, data []byte, flags int) (err error) {
	data2, nbytes := sliceToPtr(data)
	/*
		ssize_t extattr_set_file(const char *path,
			int attrnamespace, const char *attrname,
			const void *data, size_t nbytes);
	*/
	r0, _, e1 := syscall.Syscall6(syscall.SYS_EXTATTR_SET_FILE, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		EXTATTR_NAMESPACE_USER, uintptr(unsafe.Pointer(syscall.StringBytePtr(attr))),
		uintptr(unsafe.Pointer(data2)), uintptr(nbytes), 0)
	if e1 != 0 {
		return e1
	}
	if int(r0) != nbytes {
		return syscall.E2BIG
	}
	return
}

func removexattr(path string, attr string) (err error) {
	/*
		int extattr_delete_file(const char *path,
			int attrnamespace, const char *attrname);
	*/
	_, _, e1 := syscall.Syscall(syscall.SYS_EXTATTR_DELETE_FILE, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		EXTATTR_NAMESPACE_USER, uintptr(unsafe.Pointer(syscall.StringBytePtr(attr))),
	)
	if e1 != 0 {
		return e1
	}
	return
}

func listxattr(path string, dest []byte) (sz int, err error) {
	data, nbytes := sliceToPtr(dest)
	/*
	   ssize_t extattr_list_file(const char *path, int attrnamespace, void *data,
	   size_t nbytes);
	*/
	r0, _, e1 := syscall.Syscall6(syscall.SYS_EXTATTR_LIST_FILE, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		EXTATTR_NAMESPACE_USER, uintptr(unsafe.Pointer(data)), uintptr(nbytes), 0, 0)
	if e1 != 0 {
		return int(r0), e1
	}
	return int(r0), nil
}

// attrListToStrings converts a sequence of attribute name entries to a
// []string.
// On FreeBSD, each entry consists of a single byte containing the length
// of the attribute name, followed by the attribute name.
// The name is _not_ terminated by NUL.
func attrListToStrings(buf []byte) []string {
	var result []string
	index := 0
	for index < len(buf) {
		next := index + 1 + int(buf[index])
		result = append(result, string(buf[index+1:next]))
		index = next
	}
	return result
}
