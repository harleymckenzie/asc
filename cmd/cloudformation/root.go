package cloudformation

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewCloudFormationRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cloudformation",
		Short:   "Perform CloudFormation operations",
		GroupID: "service",
	}

	// Add commands
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
