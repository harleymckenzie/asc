package cloudformation

import (
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
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

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
