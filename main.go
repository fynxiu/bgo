package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/fynxiu/bgo/charset"
	"github.com/fynxiu/bgo/fir"
	"github.com/fynxiu/bgo/gitignore"
	"github.com/fynxiu/bgo/project"
	"github.com/fynxiu/bgo/upgrade"
)

const version = "0.2.2"

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
	rootCmd.AddCommand(fir.CmdFir)
	rootCmd.AddCommand(charset.CmdCharset())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
