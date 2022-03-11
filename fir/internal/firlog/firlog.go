package firlog

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/fynxiu/bgo/fir/internal/config"
	"github.com/fynxiu/bgo/fir/internal/git"
	"github.com/fynxiu/bgo/internal/fs"
	"github.com/fynxiu/bgo/internal/hash"
	"gopkg.in/yaml.v3"
)

var (
	ErrMalformedLogFile = errors.New("malformed log file")
)

const commitLen = 40

type (
	FirLog struct {
		Commit        string
		CurrentCommit string `yaml:"-"`
		Exes          []Exe
	}

	Exe struct {
		Name   string
		MD5    string
		NewMD5 string `yaml:"-"`
	}
)

func FromFile(filename string) (*FirLog, error) {
	hash, err := git.GetHeadHash()
	if err != nil {
		return nil, err
	}
	if !fs.FileExists(filename) {
		return &FirLog{
			CurrentCommit: hash,
		}, nil
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var firLog FirLog
	if err := unmarshal(contents, &firLog); err != nil {
		return nil, err
	}
	firLog.CurrentCommit = hash
	return &firLog, nil
}

func (f *FirLog) NeedBuildAll() bool {
	return f.Commit == ""
}

func (f *FirLog) NoChange() bool {
	return strings.TrimSpace(f.Commit) == strings.TrimSpace(f.CurrentCommit)
}

func (f *FirLog) ExeChangedServices(c *config.Config) ([]string, error) {
	IsHashChanged := func(s config.Service, hash string) bool {
		for _, exe := range f.Exes {
			if exe.Name == s.Name {
				exe.NewMD5 = hash
				return exe.MD5 != hash
			}
		}
		return true
	}
	var services []string
	for _, s := range c.Services {
		hash, err := hash.MD5File(path.Join(c.BuildPath, s.Name))
		if err != nil {
			return nil, err
		}
		if IsHashChanged(s, string(hash)) {
			services = append(services, s.Name)
		}
	}
	return services, nil
}

func (f *FirLog) Overwrite(filename string) error {
	f.Commit = f.CurrentCommit
	for i := range f.Exes {
		f.Exes[i].MD5 = f.Exes[i].NewMD5
	}
	contents, err := yaml.Marshal(f)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, contents, 0644)
}

func unmarshal(contents []byte, firLog *FirLog) error {
	return yaml.Unmarshal(contents, firLog)
}
