package fsassert

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"

	//revive:disable-next-line:dot-imports
	. "github.com/knaka/go-utils"
)

type dirsAreEqualParams struct {
	//IgnoreSymlinks bool
	IgnoreGlobs []string
}

// Option is a function type for configuring directory comparison options.
type Option func(params *dirsAreEqualParams)

// IgnoreGlobs returns an Option that configures directory comparison to ignore files matching the given glob patterns.
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
			Throw(fmt.Errorf("%s is not a directory", dir))
		}
		return V(filepath.Abs(dir))
	})
	if absDirs[0] == absDirs[1] {
		return fmt.Errorf("%s is the same directory as %s", dir0, dir1)
	}
	var dirLHS, dirRHS string
	fn := func(path string, info fs.FileInfo, errGiven error) (err error) {
		if errGiven != nil {
			return errGiven
		}
		tgtPath := filepath.Join(dirRHS, strings.TrimPrefix(path, dirLHS))
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
				return fmt.Errorf("%s does not exist in %s",
					tgtPath,
					dirRHS)
			}
			return err
		}
		if info.IsDir() {
			if !tgtInfo.IsDir() {
				return fmt.Errorf("%s in %s is a directory while %s in %s is not", path, dirLHS, tgtPath, dirRHS)
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
					return fmt.Errorf("%s in %s is a symlink to %s while %s in %s is a symlink to %s", path, dirLHS, link, tgtPath, dirRHS, link2)
				}
				return nil
			default:
				return nil
			}
		}
		return filesAreEqual(path, tgtPath)
	}
	dirLHS = dir0
	dirRHS = dir1
	if err = filepath.Walk(dirLHS, fn); err != nil {
		return
	}
	dirLHS = dir1
	dirRHS = dir0
	if err = filepath.Walk(dirLHS, fn); err != nil {
		return err
	}
	return nil
}
