package build

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/fynxiu/bgo/fir/internal/config"
	"github.com/fynxiu/bgo/fir/internal/git"
)

const (
	defaultBuildPath = "build"
)

func BuildWithVersion(services []config.Service, buildPath string) error {
	version, err := git.GetVersion()
	if err != nil {
		return err
	}
	_buildPath := buildPath
	if _buildPath == "" {
		buildPath = defaultBuildPath
	}
	for _, x := range services {
		if err := build(x, _buildPath, version); err != nil {
			return err
		}
	}

	return nil
}

func Build(services []config.Service, buildPath string) error {
	_buildPath := buildPath
	if _buildPath == "" {
		buildPath = defaultBuildPath
	}
	for _, x := range services {
		if err := build(x, _buildPath, ""); err != nil {
			return err
		}
	}

	return nil
}

func build(service config.Service, buildPath string, version string) error {
	var buildComand = service.BuildCommand
	buildPath = normalizeBuildPath(buildPath, service.Name)
	servicePath := normalizeServicePath(service.Path)
	ldflags := "-s -w"
	if version != "" {
		ldflags = fmt.Sprintf("%s -X main.version=%s -X main.buildTime=%s", ldflags, version, time.Now().Format(time.RFC3339))
	}
	if buildComand == nil || len(buildComand) == 0 {

		buildComand = []string{"go", "build", "-ldflags", ldflags, "-o", buildPath, servicePath}
	}
	cmd := exec.Command(buildComand[0], buildComand[1:]...)
	env := []string{"GO111MODULE=on", "GOPROXY=https://goproxy.cn,direct", "CGO_ENABLED=0"}
	cmd.Env = append(env, os.Environ()...)
	cmd.Env = append(cmd.Env, service.Env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func normalizeBuildPath(buildPath string, serviceName string) string {
	return path.Join(buildPath, serviceName)
}

func normalizeServicePath(servicePath string) string {
	if !strings.HasPrefix(servicePath, "./") {
		servicePath = "./" + servicePath
	}
	return servicePath
}
