package config

import (
	"fmt"
	"strings"
)

type Service struct {
	Name         string
	Path         string   // should be relative path without './'
	Relatives    []string // should be relative path without './'
	BuildCommand []string `yaml:"buildCommand"`
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
		if x == filename || strings.HasPrefix(filename, x) {
			return true
		}
	}
	return false
}
