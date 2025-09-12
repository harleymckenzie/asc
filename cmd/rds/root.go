package rds

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// NewRDSRootCmd creates and configures the root command for RDS operations
func NewRDSRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rds",
		Short:   "Perform RDS operations",
		Long:    "Manage Amazon RDS database instances including MySQL, PostgreSQL, and other database engines",
		GroupID: "service",
	}

	// Add action commands
	cmd.AddCommand(lsCmd)

	// Add command groups for better organization
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
