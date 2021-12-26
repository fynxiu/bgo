package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/fynxiu/bgo/gitignore"
	"github.com/fynxiu/bgo/project"
	"github.com/fynxiu/bgo/upgrade"
)

const version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "bgo",
	Short:   "bgo: tool for go project generation.",
	Long:    `bgo: A tool for go project generation.`,
	Version: version,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
	rootCmd.AddCommand(gitignore.NewCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
