package docker

import "testing"

func Test_containsAuth(t *testing.T) {
	type args struct {
		contents []byte
		endpoint string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "hp",
			args: args{
				contents: []byte(`{
					"auths": {
						"foo": {
							"auth": "bar"
						}
					}
				}`),
				endpoint: "foo",
			},
			want: true,
		},
		{
			name: "without auth",
			args: args{
				contents: []byte(`{
					"auths": {
						"foo": {
						}
					}
				}`),
				endpoint: "foo",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsAuth(tt.args.contents, tt.args.endpoint); got != tt.want {
				t.Errorf("containsAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}
