package rds

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewRDSRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rds",
		Short:   "Perform RDS operations",
		GroupID: "service",
	}

	// Add commands
	cmd.AddCommand(lsCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
