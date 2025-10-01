package fs

import (
	"os"
	"testing"

	"github.com/knaka/go-testutils/fsassert"

	//revive:disable-next-line:dot-imports
	. "github.com/knaka/go-utils"
)

func TestCopyDir(t *testing.T) {
	tempDir := V(os.MkdirTemp("", "copy_dir_test"))
	V0(CopyDir(tempDir, "testdata/dir1"))
	fsassert.DirsAreEqual(t, "testdata/dir1", tempDir)
}
