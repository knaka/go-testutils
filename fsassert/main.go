package fsassert

type tHelper interface {
	Helper()
}

type testingT interface {
	Errorf(format string, args ...interface{})
}
