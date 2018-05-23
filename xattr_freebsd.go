// +build freebsd

package xattr

import (
	"syscall"
	"unsafe"
)

const (
	EXTATTR_NAMESPACE_USER = 1

	// ENOATTR is not exported by the syscall package on Linux, because it is
	// an alias for ENODATA. We export it here so it is available on all
	// our supported platforms.
	ENOATTR = syscall.ENOATTR
)

func getxattr(path string, name string, data []byte) (int, error) {
	ptr, nbytes := bytePtrFromSlice(data)
	/*
		ssize_t extattr_get_file(
			const char *path,
			int attrnamespace,
			const char *attrname,
			void *data,
			size_t nbytes);
	*/
	r0, _, err := syscall.Syscall6(syscall.SYS_EXTATTR_GET_FILE, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		EXTATTR_NAMESPACE_USER, uintptr(unsafe.Pointer(syscall.StringBytePtr(name))),
		uintptr(unsafe.Pointer(ptr)), uintptr(nbytes), 0)
	if err != syscall.Errno(0) {
		return int(r0), err
	}
	return int(r0), nil
}

func lgetxattr(path string, name string, data []byte) (int, error) {
	ptr, nbytes := bytePtrFromSlice(data)
	/*
		ssize_t extattr_get_link(
			const char *path,
			int attrnamespace,
			const char *attrname,
			void *data,
			size_t nbytes);
	*/
	r0, _, err := syscall.Syscall6(syscall.SYS_EXTATTR_GET_LINK, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		EXTATTR_NAMESPACE_USER, uintptr(unsafe.Pointer(syscall.StringBytePtr(name))),
		uintptr(unsafe.Pointer(ptr)), uintptr(nbytes), 0)
	if err != syscall.Errno(0) {
		return int(r0), err
	}
	return int(r0), nil
}

func setxattr(path string, name string, data []byte, flags int) error {
	ptr, nbytes := bytePtrFromSlice(data)
	/*
		ssize_t extattr_set_file(
			const char *path,
			int attrnamespace,
			const char *attrname,
			const void *data,
			size_t nbytes
		);
	*/
	r0, _, err := syscall.Syscall6(syscall.SYS_EXTATTR_SET_FILE, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		EXTATTR_NAMESPACE_USER, uintptr(unsafe.Pointer(syscall.StringBytePtr(name))),
		uintptr(unsafe.Pointer(ptr)), uintptr(nbytes), 0)
	if err != syscall.Errno(0) {
		return err
	}
	if int(r0) != nbytes {
		return syscall.E2BIG
	}
	return nil
}

func lsetxattr(path string, name string, data []byte, flags int) error {
	ptr, nbytes := bytePtrFromSlice(data)
	/*
		ssize_t extattr_set_link(
			const char *path,
			int attrnamespace,
			const char *attrname,
			const void *data,
			size_t nbytes
		);
	*/
	r0, _, err := syscall.Syscall6(syscall.SYS_EXTATTR_SET_LINK, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		EXTATTR_NAMESPACE_USER, uintptr(unsafe.Pointer(syscall.StringBytePtr(name))),
		uintptr(unsafe.Pointer(ptr)), uintptr(nbytes), 0)
	if err != syscall.Errno(0) {
		return err
	}
	if int(r0) != nbytes {
		return syscall.E2BIG
	}
	return nil
}

func removexattr(path string, name string) error {
	/*
		int extattr_delete_file(
			const char *path,
			int attrnamespace,
			const char *attrname
		);
	*/
	_, _, err := syscall.Syscall(syscall.SYS_EXTATTR_DELETE_FILE, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		EXTATTR_NAMESPACE_USER, uintptr(unsafe.Pointer(syscall.StringBytePtr(name))),
	)
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func lremovexattr(path string, name string) error {
	/*
		int extattr_delete_link(
			const char *path,
			int attrnamespace,
			const char *attrname
		);
	*/
	_, _, err := syscall.Syscall(syscall.SYS_EXTATTR_DELETE_LINK, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		EXTATTR_NAMESPACE_USER, uintptr(unsafe.Pointer(syscall.StringBytePtr(name))),
	)
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func listxattr(path string, data []byte) (int, error) {
	ptr, nbytes := bytePtrFromSlice(data)
	/*
		   ssize_t extattr_list_file(
			   const char *path,
			   int attrnamespace,
			   void *data,
			   size_t nbytes
			);
	*/
	r0, _, err := syscall.Syscall6(syscall.SYS_EXTATTR_LIST_FILE, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		EXTATTR_NAMESPACE_USER, uintptr(unsafe.Pointer(ptr)), uintptr(nbytes), 0, 0)
	if err != syscall.Errno(0) {
		return int(r0), err
	}
	return int(r0), nil
}

func llistxattr(path string, data []byte) (int, error) {
	ptr, nbytes := bytePtrFromSlice(data)
	/*
		   ssize_t extattr_list_link(
			   const char *path,
			   int attrnamespace,
			   void *data,
			   size_t nbytes
			);
	*/
	r0, _, err := syscall.Syscall6(syscall.SYS_EXTATTR_LIST_LINK, uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		EXTATTR_NAMESPACE_USER, uintptr(unsafe.Pointer(ptr)), uintptr(nbytes), 0, 0)
	if err != syscall.Errno(0) {
		return int(r0), err
	}
	return int(r0), nil
}

// stringsFromByteSlice converts a sequence of attributes to a []string.
// On FreeBSD, each entry consists of a single byte containing the length
// of the attribute name, followed by the attribute name.
// The name is _not_ terminated by NULL.
func stringsFromByteSlice(buf []byte) (result []string) {
	index := 0
	for index < len(buf) {
		next := index + 1 + int(buf[index])
		result = append(result, string(buf[index+1:next]))
		index = next
	}
	return
}
