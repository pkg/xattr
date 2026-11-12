//go:build linux
// +build linux

package xattr

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"syscall"
	"testing"
)

func TestListGrowingAttributes(t *testing.T) {
	for _, family := range []string{"List", "LList", "FList"} {
		t.Run(family, func(t *testing.T) {
			f, err := ioutil.TempFile("", "xattr-list-growth-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(f.Name())
			defer f.Close()
			names := []string{UserPrefix + "initial", UserPrefix + "added-after-first-size-query", UserPrefix + "added-after-second-size-query"}
			checkIfError(t, Set(f.Name(), names[0], []byte("value")))
			syscallList := func(data []byte) (int, error) { return listxattr(f.Name(), data) }
			publicList := func() ([]string, error) { return List(f.Name()) }
			switch family {
			case "LList":
				syscallList = func(data []byte) (int, error) { return llistxattr(f.Name(), data) }
				publicList = func() ([]string, error) { return LList(f.Name()) }
			case "FList":
				syscallList = func(data []byte) (int, error) { return flistxattr(f, data) }
				publicList = func() ([]string, error) { return FList(f) }
			}
			probes, rangeErrors := 0, 0
			got, err := list(f.Name(), func(data []byte) (int, error) {
				n, err := syscallList(data)
				if err != nil {
					if err == syscall.ERANGE {
						rangeErrors++
					}
					return n, err
				}
				// Reproduce another process adding a name between the actual kernel
				// size query and retrieval, without depending on goroutine timing.
				if data == nil {
					probes++
					if probes < len(names) {
						checkIfError(t, Set(f.Name(), names[probes], []byte("value")))
					}
				}
				return n, err
			})
			if err != nil {
				t.Fatal(err)
			}
			if rangeErrors != 2 {
				t.Fatalf("saw %d real ERANGE results, want 2", rangeErrors)
			}
			sort.Strings(got)
			sort.Strings(names)
			if !reflect.DeepEqual(got, names) {
				t.Fatalf("listed %q, want %q", got, names)
			}
			got, err = publicList()
			if err != nil {
				t.Fatal(err)
			}
			sort.Strings(got)
			if !reflect.DeepEqual(got, names) {
				t.Fatalf("public list %q, want %q", got, names)
			}
		})
	}
}

func TestListPreservesErrors(t *testing.T) {
	for _, at := range []int{1, 2, 3} {
		calls := 0
		_, err := list("list-error-path", func(data []byte) (int, error) {
			calls++
			if calls == at {
				return -1, syscall.EACCES
			}
			if data == nil {
				return 1, nil
			}
			return -1, syscall.ERANGE
		})
		var xerr *Error
		if !errors.As(err, &xerr) || !errors.Is(err, syscall.EACCES) || xerr.Op != "xattr.list" || xerr.Path != "list-error-path" || xerr.Name != "" {
			t.Fatalf("call %d error = %#v, want wrapped EACCES", at, err)
		}
		if calls != at {
			t.Fatalf("made %d calls, want %d", calls, at)
		}
	}
}

func TestListRepeatedRangeError(t *testing.T) {
	calls := 0
	_, err := list("unstable-list", func(data []byte) (int, error) {
		calls++
		// A faulty or continuously changing backend must not cause an endless
		// retry. Stop the fixture itself if the helper fails to stop first.
		if calls > 100 {
			return -1, syscall.EIO
		}
		if data == nil {
			return 1, nil
		}
		return -1, syscall.ERANGE
	})
	if !errors.Is(err, syscall.ERANGE) {
		t.Fatalf("unstable list error = %v, want ERANGE", err)
	}
}
