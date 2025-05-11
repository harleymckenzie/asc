package ec2

import (
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewEC2RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ec2",
		Short:   "Perform EC2 operations",
		GroupID: "service",
	}

	// Add commands
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)
	cmd.AddCommand(restartCmd)
	cmd.AddCommand(startCmd)
	cmd.AddCommand(stopCmd)
	cmd.AddCommand(terminateCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
