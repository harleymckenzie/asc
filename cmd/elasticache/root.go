package elasticache

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)


func NewElasticacheRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "elasticache",
		Short:   "Perform Elasticache operations",
		GroupID: "service",
	}

	// Add commands
	cmd.AddCommand(lsCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
