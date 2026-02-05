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
	cmd.AddCommand(catCmd)
	cmd.AddCommand(setCmd)
	cmd.AddCommand(editCmd)
	cmd.AddCommand(cpCmd)
	cmd.AddCommand(mvCmd)
	cmd.AddCommand(rmCmd)
	cmd.AddCommand(historyCmd)
	cmd.AddCommand(labelCmd)
	cmd.AddCommand(unlabelCmd)
	cmd.AddCommand(revertCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
