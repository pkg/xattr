/*
Package xattr provides support for extended attributes on linux, darwin and freebsd.
Extended attributes are name:value pairs associated permanently with files and directories,
similar to the environment strings associated with a process.
An attribute may be defined or undefined. If it is defined, its value may be empty or non-empty.
More details you can find here: https://en.wikipedia.org/wiki/Extended_file_attributes .

All functions are provided in pairs: Get/LGet, Set/LSet etc. The "L"
variant will not follow a symlink at the end of the path.

Example assuming path is "/symlink1/symlink2", where both components are
symlinks:
Get will follow "symlink1" and "symlink2" and operate on the target of
"symlink2". LGet will follow "symlink1" but operate directly on "symlink2".
*/
package xattr

// Error records an error and the operation, file path and attribute that caused it.
type Error struct {
	Op   string
	Path string
	Name string
	Err  error
}

func (e *Error) Error() string {
	return e.Op + " " + e.Path + " " + e.Name + ": " + e.Err.Error()
}

// Get retrieves extended attribute data associated with path. It will follow
// all symlinks along the path.
func Get(path, name string) ([]byte, error) {
	return doGet(path, name, true)
}

// LGet is like Get but does not follow a symlink at the end of the path.
func LGet(path, name string) ([]byte, error) {
	return doGet(path, name, false)
}

// doGet contains the buffer allocation logic used by both Get and LGet.
func doGet(path string, name string, followSymlinks bool) ([]byte, error) {
	myname := "xattr.Get"
	getxattrFunc := lgetxattr
	if followSymlinks {
		getxattrFunc = getxattr
		myname = "xattr.LGet"
	}
	// find size.
	size, err := getxattrFunc(path, name, nil)
	if err != nil {
		return nil, &Error{myname, path, name, err}
	}
	if size > 0 {
		data := make([]byte, size)
		// Read into buffer of that size.
		read, err := getxattrFunc(path, name, data)
		if err != nil {
			return nil, &Error{myname, path, name, err}
		}
		return data[:read], nil
	}
	return []byte{}, nil
}

// Set associates name and data together as an attribute of path.
func Set(path, name string, data []byte) error {
	if err := setxattr(path, name, data, 0); err != nil {
		return &Error{"xattr.Set", path, name, err}
	}
	return nil
}

// LSet is like Set but does not follow a symlink at
// the end of the path.
func LSet(path, name string, data []byte) error {
	if err := lsetxattr(path, name, data, 0); err != nil {
		return &Error{"xattr.LSet", path, name, err}
	}
	return nil
}

// SetWithFlags associates name and data together as an attribute of path.
// Forwards the flags parameter to the syscall layer.
func SetWithFlags(path, name string, data []byte, flags int) error {
	if err := setxattr(path, name, data, flags); err != nil {
		return &Error{"xattr.SetWithFlags", path, name, err}
	}
	return nil
}

// LSetWithFlags is like SetWithFlags but does not follow a symlink at
// the end of the path.
func LSetWithFlags(path, name string, data []byte, flags int) error {
	if err := lsetxattr(path, name, data, flags); err != nil {
		return &Error{"xattr.LSetWithFlags", path, name, err}
	}
	return nil
}

// Remove removes the attribute associated with the given path.
func Remove(path, name string) error {
	if err := removexattr(path, name); err != nil {
		return &Error{"xattr.Remove", path, name, err}
	}
	return nil
}

// LRemove is like Remove but does not follow a symlink at the end of the
// path.
func LRemove(path, name string) error {
	if err := lremovexattr(path, name); err != nil {
		return &Error{"xattr.LRemove", path, name, err}
	}
	return nil
}

// List retrieves a list of names of extended attributes associated
// with the given path in the file system.
func List(path string) ([]string, error) {
	return doList(path, true)
}

// LList is like List but does not follow a symlink at the end of the
// path.
func LList(path string) ([]string, error) {
	return doList(path, false)
}

// doList contains the buffer allocation logic used by both List and LList.
func doList(path string, followSymlinks bool) ([]string, error) {
	myname := "xattr.LList"
	listxattrFunc := llistxattr
	if followSymlinks {
		listxattrFunc = listxattr
		myname = "xattr.List"
	}
	// find size.
	size, err := listxattrFunc(path, nil)
	if err != nil {
		return nil, &Error{myname, path, "", err}
	}
	if size > 0 {
		// `size + 1` because of ERANGE error when reading
		// from a SMB1 mount point (https://github.com/pkg/xattr/issues/16).
		buf := make([]byte, size+1)
		// Read into buffer of that size.
		read, err := listxattrFunc(path, buf)
		if err != nil {
			return nil, &Error{myname, path, "", err}
		}
		return stringsFromByteSlice(buf[:read]), nil
	}
	return []string{}, nil
}

// bytePtrFromSlice returns a pointer to array of bytes and a size.
func bytePtrFromSlice(data []byte) (ptr *byte, size int) {
	size = len(data)
	if size > 0 {
		ptr = &data[0]
	}
	return
}
