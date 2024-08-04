package fsassert

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_filesAreEqual(t *testing.T) {
	type args struct {
		path1 string
		path2 string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "Equals",
			args: args{
				filepath.Join("testdata", "test1.txt"),
				filepath.Join("testdata", "test2.txt"),
			},
			wantNil: true,
		},
		{
			name: "NotEquals",
			args: args{
				filepath.Join("testdata", "test1.txt"),
				filepath.Join("testdata", "test3.txt"),
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, filesAreEqual(tt.args.path1, tt.args.path2) == nil == tt.wantNil)
		})
	}
}
