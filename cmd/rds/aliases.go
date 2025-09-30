package rds

import (
	"github.com/harleymckenzie/asc/cmd/rds/cluster"

	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	// Add subcommands
	showCmd.AddCommand(clusterShowCmd)

	// Add flags
	cluster.NewShowFlags(clusterShowCmd)

	// Add groups
	showCmd.AddGroup(cmdutil.SubcommandGroups()...)
}

// Subcommand variables
var clusterShowCmd = &cobra.Command{
	Use:     "cluster",
	Short:   "Show detailed information about an RDS cluster",
	Aliases: cluster.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(cluster.ShowRDSCluster(cmd, args))
	},
}
