//go:build linux || darwin || freebsd || netbsd || solaris
// +build linux darwin freebsd netbsd solaris

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

type funcFamily struct {
	familyName   string
	get          func(path, name string) ([]byte, error)
	set          func(path, name string, data []byte) error
	setWithFlags func(path, name string, data []byte, flags int) error
	remove       func(path, name string) error
	list         func(path string) ([]string, error)
}

// Test Get, Set, List, Remove on a regular file
func TestRegularFile(t *testing.T) {
	families := []funcFamily{
		{
			familyName:   "Get and friends",
			get:          Get,
			set:          Set,
			setWithFlags: SetWithFlags,
			remove:       Remove,
			list:         List,
		},
		{
			familyName:   "LGet and friends",
			get:          LGet,
			set:          LSet,
			setWithFlags: LSetWithFlags,
			remove:       LRemove,
			list:         LList,
		},
		{
			familyName:   "FGet and friends",
			get:          wrapFGet,
			set:          wrapFSet,
			setWithFlags: wrapFSetWithFlags,
			remove:       wrapFRemove,
			list:         wrapFList,
		},
	}
	for _, ff := range families {
		t.Run(ff.familyName, func(t *testing.T) {
			t.Logf("Testing %q on a regular file", ff.familyName)
			testRegularFile(t, ff)
		})
	}
}

// testRegularFile is called with the "Get and friends" and the
// "LGet and friends" function family. Both families should behave
// the same on a regular file.
func testRegularFile(t *testing.T, ff funcFamily) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	xName := UserPrefix + "test"
	xVal := []byte("test-attr-value")

	// Test that SetWithFlags succeeds and that the xattr shows up in List()
	err = ff.setWithFlags(tmp.Name(), xName, xVal, 0)
	checkIfError(t, err)
	list, err := ff.list(tmp.Name())
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
	err = ff.remove(tmp.Name(), xName)
	checkIfError(t, err)

	// Test that Set succeeds and that the the xattr shows up in List()
	err = ff.set(tmp.Name(), xName, xVal)
	checkIfError(t, err)
	list, err = ff.list(tmp.Name())
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
	data, err = ff.get(tmp.Name(), xName)
	checkIfError(t, err)

	value := string(data)
	t.Log(value)
	if string(xVal) != value {
		t.Fail()
	}

	err = ff.remove(tmp.Name(), xName)
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

// Test that Get/LGet, Set/LSet etc operate as expected on symlinks. The
// functions should behave differently when operating on a symlink.
func TestSymlink(t *testing.T) {
	if runtime.GOOS == "solaris" || runtime.GOOS == "illumos" {
		t.Skipf("extended attributes aren't supported for symlinks on %s", runtime.GOOS)
	}
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	s := dir + "/symlink1"
	err = os.Symlink(dir+"/some/nonexistent/path", s)
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

// Verify that Get() handles values larger than the default buffer size (1 KB)
func TestLargeVal(t *testing.T) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	path := tmp.Name()

	key := UserPrefix + "TestERANGE"
	// On ext4, key + value length must be <= 4096. Use 4000 so we can test
	// reliably on ext4.
	val := bytes.Repeat([]byte("z"), 4000)
	err = Set(path, key, val)
	checkIfError(t, err)

	val2, err := Get(path, key)
	checkIfError(t, err)
	if !bytes.Equal(val, val2) {
		t.Errorf("wrong result from Get: want=%s have=%s", string(val), string(val2))
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
	err3, ok := err2.Err.(syscall.Errno)
	if !ok {
		log.Panicf("cannot unpack err2=%#v", err2)
	}
	return err3
}

// wrappers to adapt "F" variants to the test

func wrapFGet(path, name string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return FGet(f, name)
}

func wrapFSet(path, name string, data []byte) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return FSet(f, name, data)
}

func wrapFSetWithFlags(path, name string, data []byte, flags int) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return FSetWithFlags(f, name, data, flags)
}

func wrapFRemove(path, name string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return FRemove(f, name)
}

func wrapFList(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return FList(f)
}
