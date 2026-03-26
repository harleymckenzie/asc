package efs

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewEFSRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "efs",
		Short:   "Perform EFS operations",
		GroupID: "service",
	}

	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
