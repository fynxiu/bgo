package temp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_tempDir(t *testing.T) {
	type args struct {
		group string
		now   func() time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				group: "foo",
				now: func() time.Time {
					now, err := time.Parse(time.RFC3339, "2022-03-11T15:04:05Z")
					assert.NoError(t, err)
					return now
				},
			},
			want: "/tmp/bgo-foo-20220311",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tempDir(tt.args.group, tt.args.now); got != tt.want {
				t.Errorf("tempDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
