package config

import (
	"fmt"
	"strings"
)

type Service struct {
	Name         string
	Path         string   // should be relative path without './'
	Relatives    []string // docker build relatives, eg., config, asset, etc.
	Embeded      []string // go embeded files, affects go build
	BuildCommand []string `yaml:"buildCommand"`
	Env          []string // env, eg., "CGO_ENABLED=1", etc.
	Dockerfile   string   // dockerfile path
}

func (s Service) String() string {
	return s.Name
}

func (s Service) validate() error {
	if strings.HasPrefix(s.Path, "./") {
		return fmt.Errorf("Service %s Path %q should not contain './'", s.Name, s.Path)
	}
	for _, x := range s.Relatives {
		if strings.HasPrefix(x, "./") {
			return fmt.Errorf("Service %s Realtives path %q should not contain './'", s.Name, x)
		}
	}
	return nil
}

func (s Service) IsRelative(filename string) bool {
	for _, x := range s.Relatives {
		if strings.HasPrefix(filename, x) {
			return true
		}
	}
	return false
}

func (s Service) IsDockerfile(filename string) bool {
	return s.Dockerfile == filename
}

func (s Service) IsEmbeded(filename string) bool {
	for _, x := range s.Embeded {
		if strings.HasPrefix(filename, x) {
			return true
		}
	}
	return false
}
