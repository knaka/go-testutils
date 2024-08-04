package fsassert

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type doubleT struct {
	testing.T
	delim string
	logs  string
}

func (ct *doubleT) Errorf(format string, args ...any) {
	ct.logs = ct.logs + ct.delim + fmt.Sprintf(format, args...)
	ct.delim = "\n"
	ct.T.Errorf(format, args...)
}

func TestAssertDirsEqual(t *testing.T) {
	type args struct {
		dir1   string
		dir2   string
		params []Option
	}
	tests := []struct {
		name           string
		args           args
		errMsgFragment string
	}{
		{"ok",
			args{"test/foo", "test/foo2", nil},
			"",
		},
		{"no glob pattern",
			args{"test/foo", "test/foo-with-metainfo", nil},
			"does not exist in",
		},
		{
			"ok with glob pattern",
			args{
				"test/foo",
				"test/foo-with-metainfo",
				[]Option{IgnoreGlobs("metainfo.*", "*.metainfo")},
			},
			"",
		},
		{"ok with glob pattern 2",
			args{
				"test/foo-with-metainfo",
				"test/foo",
				[]Option{IgnoreGlobs("metainfo.*", "*.metainfo")},
			},
			"",
		},
		{"dir not exist",
			args{"test/no", "test/foo", nil},
			"no such file or directory",
		},
		{"dir not exist 2",
			args{"test/foo", "test/no", nil},
			"no such file or directory",
		},
		{"file specified",
			args{"test/hello.txt", "test/foo", nil},
			"not a directory",
		},
		{"file specified 2",
			args{"test/foo", "test/hello.txt", nil},
			"not a directory",
		},
		{"same directory",
			args{"test/foo", "test/foo", nil},
			"same directory",
		},
		{"file not exist",
			args{"test/foo", "test/foobar", nil},
			"bar.txt does not exist",
		},
		{"file not exist 2",
			args{"test/foobar", "test/foo", nil},
			"bar.txt does not exist",
		},
		{"content differs",
			args{"test/foo", "test/fakefoo", nil},
			"differs",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tDouble := doubleT{}
			success := DirsAreEqual(&tDouble, tt.args.dir1, tt.args.dir2, tt.args.params...)
			if success && tt.errMsgFragment != "" || !success && tt.errMsgFragment == "" {
				t.Errorf("DirsAreEqual() = %v, want %v, logs %v", success, tt.errMsgFragment, tDouble.logs)
			}
			if !success && tt.errMsgFragment != "" {
				assert.Contains(t, tDouble.logs, tt.errMsgFragment)
			}
		})
	}
}
