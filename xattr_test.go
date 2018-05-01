// +build linux darwin freebsd

package xattr

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"
)

const UserPrefix = "user."

func Test(t *testing.T) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	err = Set(tmp.Name(), UserPrefix+"test", []byte("test-attr-value"), 0)
	checkIfError(t, err)

	list, err := List(tmp.Name())
	checkIfError(t, err)

	found := false
	for _, name := range list {
		if name == UserPrefix+"test" {
			found = true
		}
	}

	if !found {
		t.Fatal("Listxattr did not return test attribute")
	}

	var data []byte
	data, err = Get(tmp.Name(), UserPrefix+"test")
	checkIfError(t, err)

	value := string(data)
	t.Log(value)
	if "test-attr-value" != value {
		t.Fail()
	}

	err = Remove(tmp.Name(), UserPrefix+"test")
	checkIfError(t, err)
}

func TestNoData(t *testing.T) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	err = Set(tmp.Name(), UserPrefix+"test", []byte{}, 0)
	checkIfError(t, err)

	list, err := List(tmp.Name())
	checkIfError(t, err)

	found := false
	for _, name := range list {
		if name == UserPrefix+"test" {
			found = true
		}
	}

	if !found {
		t.Fatal("Listxattr did not return test attribute")
	}
}

func checkIfError(t *testing.T, err error) {
	if err == nil {
		return
	}

	errno := err.(*Error)
	if errno == nil {
		t.Fatal(err)
	}

	// check if filesystem supports extended attributes
	if errno.Err == syscall.Errno(syscall.ENOTSUP) || errno.Err == syscall.Errno(syscall.EOPNOTSUPP) {
		t.Skip("Skipping test - filesystem does not support extended attributes")
	} else {
		t.Fatal(err)
	}
}
