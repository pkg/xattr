// +build darwin

package xattr

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRootless(t *testing.T) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	err = Set(tmp.Name(), "com.apple.rootless", []byte{})
	if err != nil {
		// expected: operation not permitted
		t.Log(err)
	}

	list, err := List(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range list {
		t.Fatal("xattr: " + name + " not expected")
	}
}
