package fsassert

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	. "github.com/knaka/go-utils"
	"io"
	"io/fs"
	"os"
)

func filesAreEqual(t testingT, path1, path2 string) (err error) {
	info1, err := Bind1(err, func() (fs.FileInfo, error) { return os.Stat(path1) })
	info2, err := Bind1(err, func() (fs.FileInfo, error) { return os.Stat(path2) })
	if err != nil {
		return err
	}
	if info1.IsDir() {
		return errors.New(fmt.Sprintf("%s is a directory", path1))
	}
	if info2.IsDir() {
		return errors.New(fmt.Sprintf("%s is a directory", path2))
	}
	if info1.Size() != info2.Size() {
		return errors.New(fmt.Sprintf("%s and %s differ in size", path1, path2))
	}
	reader1, err := Bind(err, func() (io.ReadCloser, error) { return os.Open(path1) })
	defer Let0(err, func() { Assert(reader1.Close() == nil) })
	hash1 := md5.New()
	_, err = Bind(err, func() (int64, error) { return io.Copy(hash1, reader1) })
	hashVal1 := hash1.Sum(nil)
	reader2, err := Bind(err, func() (io.ReadCloser, error) { return os.Open(path2) })
	defer Let0(err, func() { Assert(reader2.Close() == nil) })
	hash2 := md5.New()
	_, err = Bind(err, func() (int64, error) { return io.Copy(hash2, reader2) })
	hashVal2 := hash2.Sum(nil)
	if err != nil {
		return
	}
	if bytes.Compare(hashVal1, hashVal2) != 0 {
		return errors.New(fmt.Sprintf("%s differs from %s", path1, path2))
	}
	return
}

func FilesAreEqual(t testingT, path1, path2 string) bool {
	err := filesAreEqual(t, path1, path2)
	if err != nil {
		t.Errorf("%v", err)
		return false
	}
	return true
}
