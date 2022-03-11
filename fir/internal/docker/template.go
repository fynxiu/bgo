package docker

import (
	"html/template"
	"os"
	"path"
)

var dockerfileTemplate = `FROM alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY {{ .BuildPath }} /app/{{ .ServiceName }}

CMD ["./{{ .ServiceName }}"]
`

// defaultDockerfile write default Dockerfile to filename
// filename must be abs file path
func genDefaultDockerfile(filename, buildPath, serviceName string) error {
	// if !path.IsAbs(filename) {
	// 	return fmt.Errorf("filename should be abs path but got %q", filename)
	// }
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	t := template.Must(template.New("dockerfileTemplate").Parse(dockerfileTemplate))
	return t.Execute(fp, map[string]string{
		"BuildPath":   normalizeBuildPath(buildPath, serviceName),
		"ServiceName": serviceName,
	})
}

func normalizeBuildPath(buildPath string, serviceName string) string {
	return path.Join(buildPath, serviceName)
}
