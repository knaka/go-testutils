package fs

import (
	. "github.com/knaka/go-utils"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

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
		defer Let0(err, func() { Assert(reader.Close() == nil) })
		outFile, err := Bind(err, func() (*os.File, error) { return os.Create(outPath) })
		defer Let0(err, func() { Assert(outFile.Close() == nil) })
		err = Then(err, func() error { return outFile.Chmod(info.Mode()) })
		_, err = Bind(err, func() (int64, error) { return io.Copy(outFile, reader) })
		return
	})
}
