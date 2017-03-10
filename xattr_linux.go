// +build linux

package xattr

import "syscall"

// Get retrieves extended attribute data associated with path.
func Get(path, name string) ([]byte, error) {
	// find size.
	size, err := syscall.Getxattr(path, name, nil)
	if err != nil {
		return nil, &Error{"getxattr", path, name, err}
	}
	if size > 0 {
		data := make([]byte, size)
		// Read into buffer of that size.
		read, err := syscall.Getxattr(path, name, data)
		if err != nil {
			return nil, &Error{"getxattr", path, name, err}
		}
		return data[:read], nil
	}
	return []byte{}, nil
}

// List retrieves a list of names of extended attributes associated
// with the given path in the file system.
func List(path string) ([]string, error) {
	// find size.
	size, err := syscall.Listxattr(path, nil)
	if err != nil {
		return nil, &Error{"listxattr", path, "", err}
	}
	if size > 0 {
		buf := make([]byte, size)
		// Read into buffer of that size.
		read, err := syscall.Listxattr(path, buf)
		if err != nil {
			return nil, &Error{"listxattr", path, "", err}
		}
		return nullTermToStrings(buf[:read]), nil
	}
	return []string{}, nil
}

// Set associates name and data together as an attribute of path.
func Set(path, name string, data []byte) error {
	if err := syscall.Setxattr(path, name, data, 0); err != nil {
		return &Error{"setxattr", path, name, err}
	}
	return nil
}

// Remove removes the attribute associated
// with the given path.
func Remove(path, name string) error {
	if err := syscall.Removexattr(path, name); err != nil {
		return &Error{"removexattr", path, name, err}
	}
	return nil
}
