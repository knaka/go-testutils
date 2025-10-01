// Package fsassert is a set of assertions for filesystems.
package fsassert

type tHelper interface {
	Helper()
}

type testingT interface {
	Errorf(format string, args ...interface{})
}
