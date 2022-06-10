package charset

import "github.com/spf13/cobra"

func CmdCharset() *cobra.Command {
	cmd := &cobra.Command{
		Use: "charset",
	}
	cmd.AddCommand(addBomCmd())
	return cmd
}
