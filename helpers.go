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

// nullTermToStrings converts an array of NULL terminated UTF-8 strings to a
// []string. Used on Darwin and Linux.
func nullTermToStrings(buf []byte) (result []string) {
	offset := 0
	for index, b := range buf {
		if b == 0 {
			result = append(result, string(buf[offset:index]))
			offset = index + 1
		}
	}
	return
}
