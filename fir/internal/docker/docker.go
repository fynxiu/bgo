package docker

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/fynxiu/bgo/fir/internal/config"
	"github.com/fynxiu/bgo/fir/internal/git"
	"github.com/fynxiu/bgo/internal/fs"
)

const dockerfile = "Dockerfile"

func Build(services []config.Service, c *config.Config) error {
	if c.BuildPath == "" {
		return fmt.Errorf("buildPath should not be empty")
	}

	var registgry Registry
	if c.AliyunRegistry != nil {
		registgry = NewAliyunRegistry(*c.AliyunRegistry)
	} else {
		return fmt.Errorf("no registry config")
	}

	for _, x := range services {
		if err := build(x, c, registgry); err != nil {
			return err
		}
	}

	if len(services) > 0 {
		var req *http.Request
		var err error
		for _, x := range c.UpdatedHooks {
			req, err = http.NewRequest(http.MethodPost, x.URL, nil)
			if err != nil {
				return err
			}
			for k, vs := range x.Header {
				for _, v := range vs {
					req.Header.Set(k, v)
				}
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("UpdatedHooks failed, %v, status_code=%d", x, resp.StatusCode)
			}
		}
	}
	return nil
}

func build(service config.Service, c *config.Config, registry Registry) error {
	// find default Dockerfile
	var dockerfilePath string
	if service.Dockerfile != "" {
		dockerfilePath = service.Dockerfile
	} else if c.DockerfileDir != "" {
		dockerfilePath = path.Join(c.DockerfileDir, service.Name, dockerfile)
	} else {
		err := os.MkdirAll(c.BuildPath, 0777)
		if err != nil {
			return err
		}
		dockerfilePath = path.Join(c.BuildPath, fmt.Sprintf("%s-%s", service.Name, dockerfile))
	}
	// if no exist, generates a default
	if !fs.FileExists(dockerfilePath) {
		if err := genDefaultDockerfile(dockerfilePath, c.BuildPath, service.Name); err != nil {
			return err
		}
	}

	version, err := git.GetVersion()
	if err != nil {
		return err
	}

	imageName := fmt.Sprintf("%s/%s-%s", registry.Namesapce(), c.Project, service.Name)
	versionTag := fmt.Sprintf("%s:%s", imageName, version)
	devTag := fmt.Sprintf("%s:dev", imageName)

	// run docker build
	if err := runDockerBuild(dockerfilePath, versionTag, devTag); err != nil {
		return err
	}

	if err := registry.Login(); err != nil {
		return err
	}

	if err := runDockerPush(versionTag); err != nil {
		return err
	}

	return runDockerPush(devTag)
}

func runDockerPush(imageName string) error {
	cmd := exec.Command("docker", "push", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDockerBuild(dockerfilePath string, versionTag, devTag string) error {

	cmd := exec.Command("docker", "build", "-t", versionTag, "-t", devTag, "-f", dockerfilePath, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
