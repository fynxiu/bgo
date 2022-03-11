package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPrint(t *testing.T) {
	ss := []Service{
		{
			Name: "foo",
		},
		{
			Name: "bar",
		},
	}
	expected := "[foo bar]"
	actual := fmt.Sprintf("%v", ss)

	assert.Equal(t, expected, actual)
}

func TestService_IsEmbeded(t *testing.T) {
	type fields struct {
		Embeded []string
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			fields: fields{
				Embeded: []string{"foo/bar"},
			},
			args: args{
				filename: "foo/bar",
			},
			want: true,
		},
		{
			fields: fields{
				Embeded: []string{"foo/"},
			},
			args: args{
				filename: "foo/bar",
			},
			want: true,
		},
		{
			fields: fields{
				Embeded: []string{"foo/"},
			},
			args: args{
				filename: "foo",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				Embeded: tt.fields.Embeded,
			}
			if got := s.IsEmbeded(tt.args.filename); got != tt.want {
				t.Errorf("Service.IsEmbeded() = %v, want %v", got, tt.want)
			}
		})
	}
}
