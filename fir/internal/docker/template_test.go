package docker

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_genDefaultDockerfile(t *testing.T) {
	tempDir := t.TempDir()
	type args struct {
		filename    string
		buildPath   string
		serviceName string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		validate func()
	}{
		{
			args: args{
				filename: "foo",
			},
			wantErr: true,
		},
		{
			args: args{
				filename:    path.Join(tempDir, "foo"),
				buildPath:   "build",
				serviceName: "bar",
			},
			wantErr: false,
			validate: func() {
				contents, err := ioutil.ReadFile(path.Join(tempDir, "foo"))
				assert.NoError(t, err)
				assert.Equal(t, `FROM alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY build/bar /app/bar

CMD ["./bar"]
`, string(contents))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := genDefaultDockerfile(tt.args.filename, tt.args.buildPath, tt.args.serviceName); (err != nil) != tt.wantErr {
				t.Errorf("genDefaultDockerfile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.validate != nil {
				tt.validate()
			}
		})
	}
}
