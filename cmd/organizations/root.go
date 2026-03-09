package organizations

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// NewOrganizationsRootCmd creates and configures the root command for Organizations operations
func NewOrganizationsRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "organizations",
		Short:   "Perform AWS Organizations operations",
		Long:    "Manage AWS Organizations including accounts, organizational units, and organization details",
		Aliases: []string{"org", "organisations"},
		GroupID: "service",
	}

	// Add action commands
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	// Add command groups for better organization
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
