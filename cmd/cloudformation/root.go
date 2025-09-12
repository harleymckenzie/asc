package cloudformation

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// NewCloudFormationRootCmd creates and configures the root command for CloudFormation operations
func NewCloudFormationRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cloudformation",
		Short:   "Perform CloudFormation operations",
		Long:    "Manage AWS CloudFormation stacks and templates for infrastructure as code",
		GroupID: "service",
	}

	// Add action commands
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	// Add command groups for better organization
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
