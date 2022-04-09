//go:build linux || darwin || solaris
// +build linux darwin solaris

package xattr

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFlags(t *testing.T) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	err = SetWithFlags(tmp.Name(), UserPrefix+"flags-test", []byte("flags-test-attr-value"), 0)
	checkIfError(t, err)

	err = SetWithFlags(tmp.Name(), UserPrefix+"flags-test", []byte("flags-test-attr-value"), XATTR_CREATE)
	if err == nil {
		t.Fatalf("XATTR_CREATE should have failed because the xattr already exists")
	}
	t.Log(err)

	err = SetWithFlags(tmp.Name(), UserPrefix+"flags-test", []byte("flags-test-attr-value"), XATTR_REPLACE)
	checkIfError(t, err)

	err = Remove(tmp.Name(), UserPrefix+"flags-test")
	checkIfError(t, err)

	err = SetWithFlags(tmp.Name(), UserPrefix+"flags-test", []byte("flags-test-attr-value"), XATTR_REPLACE)
	if err == nil {
		t.Fatalf("XATTR_REPLACE should have failed because there is nothing to replace")
	}
	t.Log(err)
}
