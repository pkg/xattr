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
	nbytes := len(dest)
	var dataPtr *byte
	if nbytes > 0 {
		dataPtr = &dest[0]
	}
	return extattr_get_file(path, EXTATTR_NAMESPACE_USER, attr, dataPtr, nbytes)
}

func setxattr(path string, attr string, data []byte, flags int) (err error) {
	nbytes := len(data)
	var dataPtr *byte
	if nbytes > 0 {
		dataPtr = &data[0]
	}
	written, err := extattr_set_file(path, EXTATTR_NAMESPACE_USER, attr, dataPtr, nbytes)
	if err != nil {
		return err
	}
	if written != nbytes {
		return syscall.E2BIG
	}
	return nil
}

func removexattr(path string, attr string) (err error) {
	return extattr_delete_file(path, EXTATTR_NAMESPACE_USER, attr)
}

func listxattr(path string, dest []byte) (sz int, err error) {
	nbytes := len(dest)
	var dataPtr *byte
	if nbytes > 0 {
		dataPtr = &dest[0]
	}
	return extattr_list_file(path, EXTATTR_NAMESPACE_USER, dataPtr, nbytes)
}

func parseXattrList(buf []byte) []string {
	return attrListToStrings(buf)
}

/*
   ssize_t
   extattr_get_file(const char *path, int attrnamespace,
       const char *attrname, void *data, size_t nbytes);

   ssize_t
   extattr_set_file(const char *path, int attrnamespace,
       const char *attrname, const void *data, size_t nbytes);

   int
   extattr_delete_file(const char *path, int attrnamespace,
       const char *attrname);

   ssize_t
   extattr_list_file(const char *path, int attrnamespace, void *data,
       size_t nbytes);
*/

func extattr_get_file(path string, attrnamespace int, attrname string, data *byte, nbytes int) (int, error) {
	r, _, e := syscall.Syscall6(
		syscall.SYS_EXTATTR_GET_FILE,
		uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		uintptr(attrnamespace),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(attrname))),
		uintptr(unsafe.Pointer(data)),
		uintptr(nbytes),
		0,
	)
	var err error
	if e != 0 {
		err = e
	}
	return int(r), err
}

func extattr_set_file(path string, attrnamespace int, attrname string, data *byte, nbytes int) (int, error) {
	r, _, e := syscall.Syscall6(
		syscall.SYS_EXTATTR_SET_FILE,
		uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		uintptr(attrnamespace),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(attrname))),
		uintptr(unsafe.Pointer(data)),
		uintptr(nbytes),
		0,
	)
	var err error
	if e != 0 {
		err = e
	}
	return int(r), err
}

func extattr_delete_file(path string, attrnamespace int, attrname string) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_EXTATTR_DELETE_FILE,
		uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		uintptr(attrnamespace),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(attrname))),
	)
	var err error
	if e != 0 {
		err = e
	}
	return err
}

func extattr_list_file(path string, attrnamespace int, data *byte, nbytes int) (int, error) {
	r, _, e := syscall.Syscall6(
		syscall.SYS_EXTATTR_LIST_FILE,
		uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
		uintptr(attrnamespace),
		uintptr(unsafe.Pointer(data)),
		uintptr(nbytes),
		0,
		0,
	)
	var err error
	if e != 0 {
		err = e
	}
	return int(r), err
}

// attrListToStrings converts a sequnce of attribute name entries to a []string.
// Each entry consists of a single byte containing the length
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
