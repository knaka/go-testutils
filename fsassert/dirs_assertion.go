package fsassert

import (
	"errors"
	"fmt"
	. "github.com/knaka/go-utils"
	"github.com/samber/lo"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type dirsAreEqualParams struct {
	//IgnoreSymlinks bool
	IgnoreGlobs []string
}

type Option func(params *dirsAreEqualParams)

func IgnoreGlobs(globs ...string) Option {
	return func(params *dirsAreEqualParams) {
		params.IgnoreGlobs = globs
	}
}

// DirsAreEqual asserts that two directories are equal.
func DirsAreEqual(t testingT, dir0 string, dir1 string, opts ...Option) (f bool) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	params := &dirsAreEqualParams{}
	for _, opt := range opts {
		opt(params)
	}
	err := dirsAreEqual(dir0, dir1, params)
	if err != nil {
		t.Errorf("%v", err)
		return false
	}
	return true
}

func dirsAreEqual(dir0 string, dir1 string, params *dirsAreEqualParams) (err error) {
	defer Catch(&err)
	absDirs := lo.Map([]string{dir0, dir1}, func(dir string, _ int) string {
		if !V(os.Stat(dir)).IsDir() {
			Throw(errors.New(fmt.Sprintf("%s is not a directory", dir)))
		}
		return V(filepath.Abs(dir))
	})
	if absDirs[0] == absDirs[1] {
		return errors.New(fmt.Sprintf("%s is the same directory as %s", dir0, dir1))
	}
	var dirLhs, dirRhs string
	fn := func(path string, info fs.FileInfo, errGiven error) (err error) {
		if errGiven != nil {
			return errGiven
		}
		tgtPath := filepath.Join(dirRhs, strings.TrimPrefix(path, dirLhs))
		if params != nil {
			for _, glob := range params.IgnoreGlobs {
				if matched, errSub := filepath.Match(glob, filepath.Base(path)); errSub != nil {
					return errSub
				} else if matched {
					return nil
				}
				if matched, errSub := filepath.Match(glob, filepath.Base(tgtPath)); errSub != nil {
					return errSub
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
		return filesAreEqual(path, tgtPath)
	}
	dirLhs = dir0
	dirRhs = dir1
	if err = filepath.Walk(dirLhs, fn); err != nil {
		return
	}
	dirLhs = dir1
	dirRhs = dir0
	if err = filepath.Walk(dirLhs, fn); err != nil {
		return err
	}
	return nil
}
