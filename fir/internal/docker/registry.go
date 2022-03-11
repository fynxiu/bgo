package docker

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/fynxiu/bgo/fir/internal/config"
)

type Registry interface {
	Namesapce() string
	Login() error
}

var _ Registry = (*aliyunRegistry)(nil)

func NewAliyunRegistry(config config.AliyunRegistry) Registry {
	return &aliyunRegistry{
		config: config,
	}
}

type aliyunRegistry struct {
	config config.AliyunRegistry
}

func (r *aliyunRegistry) Namesapce() string {
	return path.Join(r.config.Endpoint, r.config.Namesapce)
}

func (r *aliyunRegistry) Login() error {
	if r.alreadyAuthed() {
		return nil
	}
	cmd := exec.Command("docker", "login", "--username", r.config.Username, "--password", r.config.Password, r.config.Endpoint)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *aliyunRegistry) alreadyAuthed() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	contents, err := ioutil.ReadFile(path.Join(homeDir, ".docker/config.json"))
	if err != nil {
		return false
	}
	return containsAuth(contents, r.config.Endpoint)
}

func containsAuth(contents []byte, endpoint string) bool {
	var c struct {
		Auths map[string]struct {
			Auth string
		}
	}
	if err := json.Unmarshal(contents, &c); err != nil {
		return false
	}
	auth, ok := c.Auths[endpoint]
	if !ok {
		return false
	}
	return auth.Auth != ""
}
