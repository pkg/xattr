package xattr

import (
	"syscall"
	"testing"
)

func TestIgnoringEINTR(t *testing.T) {
	eintrs := 100
	err := ignoringEINTR(func() error {
		if eintrs == 0 {
			return nil
		}
		eintrs--
		return syscall.EINTR
	})

	if err != nil {
		t.Fatal(err)
	}
}
