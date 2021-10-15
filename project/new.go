package project

import (
	"fmt"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"

	"github.com/fyn.xiu/bgo/internal/fs"
	"github.com/fyn.xiu/bgo/internal/repo"
)

// Project is a project template.
type Project struct {
	Name string
	Path string
}

// New a project from remote repo.
func (p *Project) New(dir string, layout string, branch string) error {
	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "📂 Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return e
		}
		if !override {
			return err
		}
		err := os.RemoveAll(to)
		if err != nil {
			return err
		}
	}
	fmt.Printf("🚀 Creating service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)
	r := repo.New(layout, branch)
	if err := r.CopyTo(to, p.Path, []string{".git", ".github"}); err != nil {
		return err
	}
	fs.Tree(to, dir)

	fmt.Printf("\n🍺 Project creation succeeded %s\n", color.GreenString(p.Name))
	return nil
}
