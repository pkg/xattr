/*
Package xattr provides support for extended attributes on linux, darwin and freebsd.
Extended attributes are name:value pairs associated permanently with files and directories,
similar to the environment strings associated with a process.
An attribute may be defined or undefined. If it is defined, its value may be empty or non-empty.
More details you can find here: https://en.wikipedia.org/wiki/Extended_file_attributes
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

// Get retrieves extended attribute data associated with path.
func Get(path, name string) ([]byte, error) {
	// find size.
	size, err := getxattr(path, name, nil)
	if err != nil {
		return nil, &Error{"xattr.Get", path, name, err}
	}
	if size > 0 {
		data := make([]byte, size)
		// Read into buffer of that size.
		read, err := getxattr(path, name, data)
		if err != nil {
			return nil, &Error{"xattr.Get", path, name, err}
		}
		return data[:read], nil
	}
	return []byte{}, nil
}

// Set associates name and data together as an attribute of path.
func Set(path, name string, data []byte) error {
	return SetWithFlags(path, name, data, 0)
}

// SetWithFlags associates name and data together as an attribute of path. Forwards the flags parameter to the syscall layer.
func SetWithFlags(path, name string, data []byte, flags int) error {
	if err := setxattr(path, name, data, flags); err != nil {
		return &Error{"xattr.SetWithFlags", path, name, err}
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

// List retrieves a list of names of extended attributes associated
// with the given path in the file system.
func List(path string) ([]string, error) {
	// find size.
	size, err := listxattr(path, nil)
	if err != nil {
		return nil, &Error{"xattr.List", path, "", err}
	}
	if size > 0 {
		// `size + 1` because of ERANGE error when reading
		// from a SMB1 mount point (https://github.com/pkg/xattr/issues/16).
		buf := make([]byte, size+1)
		// Read into buffer of that size.
		read, err := listxattr(path, buf)
		if err != nil {
			return nil, &Error{"xattr.List", path, "", err}
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
