// +build linux darwin freebsd

package xattr

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"syscall"
	"testing"
)

const UserPrefix = "user."

// Test Get, Set, List, Remove on a regular file
func TestRegularFile(t *testing.T) {
	doTestRegularFile(t, false)
}

// Test LGet, LSet, LList, LRemove on a regular file
func TestRegularFileL(t *testing.T) {
	doTestRegularFile(t, true)
}

func doTestRegularFile(t *testing.T, followSymlinks bool) {
	setFunc := LSet
	setWithFlagsFunc := LSetWithFlags
	getFunc := LGet
	removeFunc := LRemove
	listFunc := LList
	if followSymlinks {
		setFunc = Set
		setWithFlagsFunc = SetWithFlags
		getFunc = Get
		removeFunc = Remove
		listFunc = List
	}
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	xName := UserPrefix + "test"
	xVal := []byte("test-attr-value")

	// Test that SetWithFlags succeeds and that the xattr shows up in List()
	err = setWithFlagsFunc(tmp.Name(), xName, xVal, 0)
	checkIfError(t, err)
	list, err := listFunc(tmp.Name())
	checkIfError(t, err)
	found := false
	for _, name := range list {
		if name == xName {
			found = true
		}
	}
	if !found {
		t.Fatalf("List/LList did not return test attribute: %q", list)
	}
	err = removeFunc(tmp.Name(), xName)
	checkIfError(t, err)

	// Test that Set succeeds and that the the xattr shows up in List()
	err = setFunc(tmp.Name(), xName, xVal)
	checkIfError(t, err)
	list, err = listFunc(tmp.Name())
	checkIfError(t, err)
	found = false
	for _, name := range list {
		if name == xName {
			found = true
		}
	}
	if !found {
		t.Fatalf("List/LList did not return test attribute: %q", list)
	}

	var data []byte
	data, err = getFunc(tmp.Name(), xName)
	checkIfError(t, err)

	value := string(data)
	t.Log(value)
	if string(xVal) != value {
		t.Fail()
	}

	err = removeFunc(tmp.Name(), xName)
	checkIfError(t, err)
}

// Test that setting an xattr with an empty value works.
func TestNoData(t *testing.T) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	err = Set(tmp.Name(), UserPrefix+"test", []byte{})
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

// Test that Get/LGet, Set/LSet etc operate as expected on symlinks.
func TestSymlink(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	s := dir + "/symlink1"
	err = os.Symlink("/some/nonexistent/path", s)
	if err != nil {
		t.Fatal(err)
	}
	xName := UserPrefix + "TestSymlink"
	xVal := []byte("test")

	// Test Set/LSet
	if err := Set(s, xName, xVal); err == nil {
		t.Error("Set on a broken symlink should fail, but did not")
	}
	err = LSet(s, xName, xVal)
	errno := unpackSysErr(err)
	setOk := true
	if runtime.GOOS == "linux" && errno == syscall.EPERM {
		// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/fs/xattr.c?h=v4.17-rc5#n122 :
		// In the user.* namespace, only regular files and directories can have
		// extended attributes.
		t.Log("got EPERM, adjusting test scope")
		setOk = false
	} else {
		checkIfError(t, err)
	}

	// Test List/LList
	_, err = List(s)
	errno = unpackSysErr(err)
	if errno != syscall.ENOENT {
		t.Errorf("List() on a broken symlink should fail with ENOENT, got %q", errno)
	}
	data, err := LList(s)
	checkIfError(t, err)
	if setOk {
		found := false
		for _, n := range data {
			if n == xName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("xattr %q did not show up in Llist output: %q", xName, data)
		}
	}

	// Test Get/LGet
	_, err = Get(s, xName)
	errno = unpackSysErr(err)
	if errno != syscall.ENOENT {
		t.Errorf("Get() on a broken symlink should fail with ENOENT, got %q", errno)
	}
	val, err := LGet(s, xName)
	if setOk {
		checkIfError(t, err)
		if !bytes.Equal(xVal, val) {
			t.Errorf("wrong xattr value: want=%q have=%q", xVal, val)
		}
	} else {
		errno = unpackSysErr(err)
		if errno != ENOATTR {
			t.Errorf("expected ENOATTR, got %q", errno)
		}
	}

	// Test Remove/Lremove
	err = Remove(s, xName)
	errno = unpackSysErr(err)
	if errno != syscall.ENOENT {
		t.Errorf("Remove() on a broken symlink should fail with ENOENT, got %q", errno)
	}
	err = LRemove(s, xName)
	if setOk {
		checkIfError(t, err)
	} else {
		errno = unpackSysErr(err)
		if errno != syscall.EPERM {
			t.Errorf("expected EPERM, got %q", errno)
		}
	}
}

// checkIfError calls t.Skip() if the underlying syscall.Errno is
// ENOTSUP or EOPNOTSUPP. It calls t.Fatal() on any other non-zero error.
func checkIfError(t *testing.T, err error) {
	errno := unpackSysErr(err)
	if errno == syscall.Errno(0) {
		return
	}
	// check if filesystem supports extended attributes
	if errno == syscall.Errno(syscall.ENOTSUP) || errno == syscall.Errno(syscall.EOPNOTSUPP) {
		t.Skip("Skipping test - filesystem does not support extended attributes")
	} else {
		t.Fatal(err)
	}
}

// unpackSysErr unpacks the underlying syscall.Errno from an error value
// returned by Get/Set/...
func unpackSysErr(err error) syscall.Errno {
	if err == nil {
		return syscall.Errno(0)
	}
	err2, ok := err.(*Error)
	if !ok {
		log.Panicf("cannot unpack err=%#v", err)
	}
	return err2.Err.(syscall.Errno)
}
