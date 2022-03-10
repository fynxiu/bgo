package git

import (
	"reflect"
	"testing"
)

func Test_changedFilefList(t *testing.T) {
	type args struct {
		prevCommit string
		currCommit string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "hp",
			args: args{
				prevCommit: "bbcad2e2a127ced30dfd9c3e89ff098adc46ce61",
				currCommit: "f9c8e27fcdc7cb977f7e84b0b470a6c0978581cd",
			},
			want: []string{
				"gitignore/gitignore.go",
				"gitignore/gitignore_test.go",
				"main.go",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := changedFilefList(tt.args.prevCommit, tt.args.currCommit)
			if (err != nil) != tt.wantErr {
				t.Errorf("changedFilefList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("changedFilefList() = %v, want %v", got, tt.want)
			}
		})
	}
}
