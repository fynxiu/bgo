package repo

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/fynxiu/bgo/internal/fs"
	"github.com/fynxiu/bgo/internal/gomod"
)

// Repo is git repository manager.
type Repo struct {
	url    string
	home   string
	branch string
}

// New a repository manager.
func New(url string, branch string) *Repo {
	var start int
	start = strings.Index(url, "//")
	if start == -1 {
		start = strings.Index(url, ":") + 1
	} else {
		start += 2
	}
	end := strings.LastIndex(url, "/")
	return &Repo{
		url:    url,
		home:   fs.BgoHomeWithDir(path.Join("repo", url[start:end])),
		branch: branch,
	}
}

// Path returns the repository cache path.
func (r *Repo) Path() string {
	start := strings.LastIndex(r.url, "/")
	end := strings.LastIndex(r.url, ".git")
	if end == -1 {
		end = len(r.url)
	}
	var branch string
	if r.branch == "" {
		branch = "@main"
	} else {
		branch = "@" + r.branch
	}
	return path.Join(r.home, r.url[start+1:end]+branch)
}

// Pull fetch the repository from remote url.
func (r *Repo) Pull() error {
	cmd := exec.Command("git", "symbolic-ref", "HEAD")
	cmd.Dir = r.Path()
	err := cmd.Run()
	if err != nil {
		return nil
	}
	cmd = exec.Command("git", "pull")
	cmd.Dir = r.Path()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	return err
}

// Clone clones the repository to cache path.
func (r *Repo) Clone() error {
	if _, err := os.Stat(r.Path()); !os.IsNotExist(err) {
		return r.Pull()
	}
	var cmd *exec.Cmd
	if r.branch == "" {
		cmd = exec.Command("git", "clone", r.url, r.Path())
	} else {
		cmd = exec.Command("git", "clone", "-b", r.branch, r.url, r.Path())
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err
}

// CopyTo copies the repository to project path.
func (r *Repo) CopyTo(to string, modPath string, ignores []string) error {
	if err := r.Clone(); err != nil {
		return err
	}
	mod, err := gomod.ModulePath(r.Path())
	if err != nil {
		return err
	}
	return fs.CopyDir(r.Path(), to, []string{mod, modPath}, ignores)
}
