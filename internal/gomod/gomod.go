package gomod

import (
	"io/ioutil"
	"path"

	"golang.org/x/mod/modfile"
)

const FileName = "go.mod"

// ModulePath returns go module path.
func ModulePath(filePath string) (string, error) {
	modBytes, err := ioutil.ReadFile(path.Join(filePath, FileName))
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(modBytes), nil
}