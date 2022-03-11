package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

var (
	ErrUnsupportedFormat = errors.New("unsupported config format")
)

type Config struct {
	Project        string // procject name
	Dockerizing    bool
	BuildAll       bool     `yaml:"buildAll"`
	BuildPath      string   `yaml:"buildPath"`
	Excluded       []string // must be relative path without './'
	Services       []Service
	DockerfileDir  string          `yaml:"dockerfileDir"` // where to find Dockerfile, DockerfileDir/ServiceName/Dockerfile
	AliyunRegistry *AliyunRegistry `yaml:"aliyunRegistry"`
}

func FromFile(filename string) (*Config, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	unmarshaller, ok := unmarshallers[path.Ext(filename)]
	if !ok {
		return nil, ErrUnsupportedFormat
	}
	var c Config
	if err := unmarshaller(contents, &c); err != nil {
		return nil, err
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) validate() error {
	if len(c.Services) == 0 {
		return fmt.Errorf("no services")
	}

	for _, x := range c.Excluded {
		if strings.HasPrefix(x, "./") {
			return fmt.Errorf("excluded filepath %q should not contain './'", x)
		}
	}

	for _, s := range c.Services {
		if err := s.validate(); err != nil {
			return err
		}
	}

	return nil
}
