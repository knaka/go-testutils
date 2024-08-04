package fs

import (
	"github.com/knaka/go-testutils/fsassert"
	"os"
	"path/filepath"
	"testing"

	. "github.com/knaka/go-utils"
)

func TestCopyDir(t *testing.T) {
	tempDir := V(os.MkdirTemp("", "copy_dir_test"))
	V0(CopyDir(filepath.Join(tempDir), "testdata/dir1"))
	fsassert.DirsAreEqual(t, "testdata/dir1", tempDir)
}
