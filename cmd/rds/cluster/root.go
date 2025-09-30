package cluster

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"clusters"}

// Root command
func NewClusterRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cluster",
		Short:   "Perform cluster operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	// Add the subcommands to the command
	cmd.AddCommand(showCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
