package project

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)
// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a project template",
	Long:  "Create a project using the repository template. Example: bgo new hello",
	Run:   run,
}


var (
	repoURL string
	branch  string
)

func init() {
	if repoURL = os.Getenv("BGO_LAYOUT_REPO"); repoURL == "" {
		repoURL = "https://github.com/fynxiu/bgo-layout.git"
	}
	CmdNew.Flags().StringVarP(&repoURL, "repo", "r", repoURL, "url of layout repo")
	CmdNew.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
}

func run(_ *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	name := ""
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "What is project name ?",
			Help:    "Created project name.",
		}
		err = survey.AskOne(prompt, &name)
		if err != nil || name == "" {
			return
		}
	} else {
		name = args[0]
	}
	p := &Project{Name: path.Base(name), Path: name}
	if err := p.New(wd, repoURL, branch); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
		if err != nil {
			log.Printf("fmt.Fprintf failed, %v", err)
		}
		return
	}
}