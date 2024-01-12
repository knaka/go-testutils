# Go-fsassert

Go-fsassert is a assertion library for testing files equiality.

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
