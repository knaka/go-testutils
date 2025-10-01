// Package fs is a set of file utilities for testing.
package fs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	//revive:disable-next-line:dot-imports
	. "github.com/knaka/go-utils"
)

// CanonAbs returns the canonical absolute path of the given value.
func CanonAbs(s string) (ret string, err error) {
	ret, err = filepath.Abs(s)
	if err != nil {
		return
	}
	ret, err = filepath.EvalSymlinks(ret)
	if err != nil {
		return
	}
	ret = filepath.Clean(ret)
	return
}

// CopyFile copies a file from src to dst.
func CopyFile(dst, src string) (err error) {
	reader := V(os.Open(src))
	defer (func() { V0(reader.Close()) })()
	writer := V(os.Create(dst))
	defer (func() { V0(writer.Close()) })()
	_, err = io.Copy(writer, reader)
	return
}

// CopyDir copies a directory recursively.
func CopyDir(dst, src string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, errGiven error) (err error) {
		if errGiven != nil {
			return errGiven
		}
		outPath := filepath.Join(dst, strings.TrimPrefix(path, src))
		if info.IsDir() {
			return os.MkdirAll(outPath, info.Mode())
		}
		if !info.Mode().IsRegular() {
			switch info.Mode().Type() & os.ModeType {
			case os.ModeSymlink:
				link, errSub := os.Readlink(path)
				if errSub != nil {
					return errSub
				}
				return os.Symlink(link, outPath)
			default:
				return nil
			}
		}
		reader, err := Bind(err, func() (io.ReadCloser, error) { return os.Open(path) })
		defer Let0(err, func() {
			if reader.Close() != nil {
				panic("?")
			}
		})
		outFile, err := Bind(err, func() (*os.File, error) { return os.Create(outPath) })
		defer Let0(err, func() {
			if outFile.Close() != nil {
				panic("?")
			}
		})
		err = Then(err, func() error { return outFile.Chmod(info.Mode()) })
		_, err = Bind(err, func() (int64, error) { return io.Copy(outFile, reader) })
		return
	})
}
