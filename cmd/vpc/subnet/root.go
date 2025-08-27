package subnet

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"subnets"}

// NewSubnetRootCmd returns the root command for Subnet operations.
func NewSubnetRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subnet",
		Short:   "Perform Subnet operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	cmd.AddCommand(lsCmd)
	// cmd.AddCommand(showCmd) // Disabled - show.go.disabled

	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)

	return cmd
}
