# Go-testutils

Go-testutils is a collection of utilities for testing.

## `fs` package

See it in action:

```go
package main

import (
	testfsutils "github.com/knaka/go-testutils/fs"
)

func TestSomething(t *testing.T) {
	err := testfsutils.CopyDir("dest/dir", "src/dir")
}
```

## `fsassert` package

See it in action:

```go
package main

import (
	"github.com/knaka/go-testutils/fsassert"
)

func TestSomething(t *testing.T) {
	fsassert.FilesAreEqual(t, "testdata/expected.txt", "testdata/actual.txt")
	fsassert.DirsAreEqual(t, "testdata/expected", "testdata/actual"
}
```
