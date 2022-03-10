package service

import (
	"reflect"
	"testing"
)

func Test_filterNoSourceFiles(t *testing.T) {
	type args struct {
		changedFileList []string
	}
	tests := []struct {
		name    string
		args    args
		wantRet []string
	}{
		{
			args: args{
				changedFileList: []string{"foo.go", "foo", "foo_test.go"},
			},
			wantRet: []string{"foo.go"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := filterNoSourceFiles(tt.args.changedFileList); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("filterNoSourceFiles() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
