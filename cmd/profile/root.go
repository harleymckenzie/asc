package profile

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewProfileRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profile",
		Short:   "Manage AWS CLI profiles",
		GroupID: "service",
	}

	cmd.AddCommand(lsCmd)

	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
