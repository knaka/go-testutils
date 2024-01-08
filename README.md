# Go-fsassert

Go-fsassert is a library for testing file system operations.

## `fsassert` package

See it in action:

```go
package main

import (
	"github.com/knaka/go-fsassert"
)

func TestSomething(t *testing.T) {
	fsassert.FilesAreEqual(t, "testdata/expected.txt", "testdata/actual.txt")
	fsassert.DirsAreEqual(t, "testdata/expected", "testdata/actual"
}
```
