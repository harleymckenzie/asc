package route_table

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"route-tables", "routetable", "routetables"}

// NewRouteTableRootCmd returns the root command for Route Table operations.
func NewRouteTableRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "route-table",
		Short:   "Perform Route Table operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)

	return cmd
}
