package upgrade

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fynxiu/bgo/internal/goexec"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the bgo tools",
	Long:  "Upgrade the bgo tools. Example: bgo upgrade",
	Run:   Run,
}

// Run upgrade the bgo tools.
func Run(_ *cobra.Command, _ []string) {
	err := goexec.Install(
		"github.com/fynxiu/bgo",
	)
	if err != nil {
		fmt.Println(err)
	}
}
