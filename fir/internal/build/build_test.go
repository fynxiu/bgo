package build

import "testing"

func Test_normalizeBuildPath(t *testing.T) {
	type args struct {
		buildPath   string
		serviceName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				buildPath:   "build",
				serviceName: "foo",
			},
			want: "build/foo",
		},
		{
			args: args{
				buildPath:   "build/",
				serviceName: "foo",
			},
			want: "build/foo",
		},
		{
			args: args{
				buildPath:   "./build/",
				serviceName: "foo",
			},
			want: "build/foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeBuildPath(tt.args.buildPath, tt.args.serviceName); got != tt.want {
				t.Errorf("normalizeBuildPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
