// +build linux darwin

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

	err = SetWithFlags(tmp.Name(), UserPrefix+"flags-test", []byte("flags-test-attr-value"), XATTR_CREATE)
	checkIfError(t, err)

	err = SetWithFlags(tmp.Name(), UserPrefix+"flags-test", []byte("flags-test-attr-value"), XATTR_REPLACE)
	checkIfError(t, err)
}
