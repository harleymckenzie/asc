package cluster

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"clusters"}

// NewClusterRootCmd creates and configures the root command for ECS cluster operations
func NewClusterRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cluster",
		Short:   "Perform ECS cluster operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
