package fsassert

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"io"
	"os"

	. "github.com/knaka/go-utils"
)

func filesAreEqual(path1 string, path2 string) (err error) {
	defer Catch(&err)
	hashVals := lo.Map([]string{path1, path2}, func(path string, _ int) []byte {
		if V(os.Stat(path)).IsDir() {
			Throw(errors.New(fmt.Sprintf("%s is a directory", path)))
		}
		reader := V(os.Open(path))
		defer (func() { V0(reader.Close()) })()
		hash := sha1.New()
		V0(io.Copy(hash, reader))
		return hash.Sum(nil)
	})
	return TernaryF(bytes.Compare(hashVals[0], hashVals[1]) == 0,
		func() error { return nil },
		func() error { return errors.New(fmt.Sprintf("%s differs from %s", path1, path2)) },
	)
}

//goland:noinspection GoUnusedExportedFunction
func FilesAreEqual(t testingT, path1 string, path2 string) bool {
	err := filesAreEqual(path1, path2)
	if err != nil {
		t.Errorf("%v", err)
		return false
	}
	return true
}
