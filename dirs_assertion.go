package fsassert

import (
	"errors"
	"fmt"
	. "github.com/knaka/go-utils"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type DirsAreEqualParams struct {
	//IgnoreSymlinks bool
	IgnoreGlobs []string
}

// DirsAreEqual asserts that two directories are equal.
func DirsAreEqual(t testingT, dir1 string, dir2 string, params *DirsAreEqualParams) (b bool) {
	var err error
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	stat1, err := Bind(err, func() (fs.FileInfo, error) { return os.Stat(dir1) })
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if !stat1.IsDir() {
		t.Errorf("%s is not a directory", dir1)
		return
	}
	stat2, err := Bind(err, func() (fs.FileInfo, error) { return os.Stat(dir2) })
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if !stat2.IsDir() {
		t.Errorf("%s is not a directory", dir2)
		return
	}
	absDir1, err := Bind(err, func() (string, error) { return filepath.Abs(dir1) })
	absDir2, err := Bind(err, func() (string, error) { return filepath.Abs(dir2) })
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if absDir1 == absDir2 {
		t.Errorf("%s and %s are the same directory", dir1, dir2)
		return
	}
	var dirLhs, dirRhs string
	f := func(path string, info fs.FileInfo, errGiven error) (err error) {
		if errGiven != nil {
			return errGiven
		}
		tgtPath := filepath.Join(dirRhs, strings.TrimPrefix(path, dirLhs))
		if params != nil {
			for _, glob := range params.IgnoreGlobs {
				if matched, err := filepath.Match(glob, filepath.Base(path)); err != nil {
					return err
				} else if matched {
					return nil
				}
				if matched, err := filepath.Match(glob, filepath.Base(tgtPath)); err != nil {
					return err
				} else if matched {
					return nil
				}
			}
		}
		tgtInfo, err := os.Stat(tgtPath)
		if err != nil {
			if os.IsNotExist(err) {
				return errors.New(fmt.Sprintf("%s does not exist in %s",
					tgtPath,
					dirRhs,
				))
			}
			return err
		}
		if info.IsDir() {
			if !tgtInfo.IsDir() {
				return errors.New(fmt.Sprintf("%s in %s is a directory while %s in %s is not", path, dirLhs, tgtPath, dirRhs))
			}
			return nil
		}
		if !info.Mode().IsRegular() {
			switch info.Mode().Type() & os.ModeType {
			case os.ModeSymlink:
				link, errSub := os.Readlink(path)
				if errSub != nil {
					return errSub
				}
				link2, errSub := os.Readlink(tgtPath)
				if errSub != nil {
					return errSub
				}
				if link != link2 {
					return errors.New(fmt.Sprintf("%s in %s is a symlink to %s while %s in %s is a symlink to %s", path, dirLhs, link, tgtPath, dirRhs, link2))
				}
				return nil
			default:
				return nil
			}
		}
		return filesAreEqual(t, path, tgtPath)
	}
	err = Then(err, func() error {
		dirLhs = dir1
		dirRhs = dir2
		return filepath.Walk(dirLhs, f)
	})
	err = Then(err, func() error {
		dirLhs = dir2
		dirRhs = dir1
		return filepath.Walk(dirLhs, f)
	})
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	b = true
	return
}
