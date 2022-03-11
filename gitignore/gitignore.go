package gitignore

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fynxiu/bgo/internal/repo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const gitignoreUrl = "https://github.com/github/gitignore.git"

var (
	flagDir            string
	flagIgnoreElements elements
)

type elements []string

var _ pflag.Value = (*elements)(nil)

func (e *elements) String() string {
	return fmt.Sprint(*e)
}

func (e *elements) Set(value string) error {
	na := strings.Fields(value)
	*e = append(*e, na...)
	return nil
}

func (e *elements) Type() string {
	return "string"
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitignore",
		Short: "Geenerate .gitignore file",
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := newApp()
			app.flagDir = flagDir
			app.flagIgnoreElements = flagIgnoreElements
			if err != nil {
				panic(err)
			}
			filepath := path.Join(app.flagDir, ".gitignore")
			if _, err := os.Stat(filepath); !os.IsNotExist(err) {
				log.Fatalf("%s already exists.", filepath)
				return nil
			}

			f, err := os.Create(filepath)
			if err != nil {
				return err
			}
			defer f.Close()
			return app.generate(app.flagIgnoreElements, f)
		},
		Args: cobra.NoArgs,
	}
	cmd.Flags().StringVarP(&flagDir, "dir", "d", ".", "directory of the workdir")
	cmd.Flags().VarP(&flagIgnoreElements, "elements", "e", "elements for the .gitgnore")
	return cmd
}

type app struct {
	ignores map[string]string

	flagIgnoreElements elements
	flagDir            string
}

func newApp() (*app, error) {
	r := repo.New(gitignoreUrl, "main")
	if err := r.Clone(); err != nil {
		return nil, err
	}
	ignores, err := findGitignores(r.Path())
	if err != nil {
		return nil, err
	}
	return &app{
		ignores: ignores,
	}, nil
}

func findGitignores(repoPath string) (map[string]string, error) {
	var err error
	_, err = ioutil.ReadDir(repoPath)
	if err != nil {
		return nil, err
	}

	filelist := make(map[string]string)
	filepath.Walk(repoPath, func(filepath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".gitignore") {
			name := strings.ToLower(strings.Replace(info.Name(), ".gitignore", "", 1))
			filelist[name] = filepath
		}
		return nil
	})
	return filelist, nil
}

func (a *app) availableFiles() ([]string, error) {
	availableGitignores := []string{}
	for key := range a.ignores {
		availableGitignores = append(availableGitignores, key)
	}

	return availableGitignores, nil
}

func (a *app) generate(names []string, output io.Writer) error {
	notFound := []string{}
	for index, name := range names {
		if filepath, ok := a.ignores[strings.ToLower(name)]; ok {
			bytes, err := ioutil.ReadFile(filepath)
			if err == nil {
				if _, err = output.Write([]byte("\n#### " + name + " ####\n")); err != nil {
					return err
				}
				if _, err = output.Write(bytes); err != nil {
					return err
				}
				if index < len(names)-1 {
					if _, err = output.Write([]byte("\n")); err != nil {
						return err
					}
				}
				continue
			}
		} else {
			notFound = append(notFound, name)
		}
	}

	if len(notFound) > 0 {
		return fmt.Errorf("not found %q", strings.Join(notFound, ","))
	}
	return nil
}
