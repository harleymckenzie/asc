package ssm

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewSSMRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssm",
		Short:   "Perform SSM Parameter Store operations",
		GroupID: "service",
	}

	// Add commands
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)
	cmd.AddCommand(cpCmd)
	cmd.AddCommand(mvCmd)
	cmd.AddCommand(rmCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
