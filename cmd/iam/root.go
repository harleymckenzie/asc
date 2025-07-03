package iam

import (
	"github.com/harleymckenzie/asc/cmd/iam/role"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewIAMRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "iam",
		Short:   "Perform IAM operations",
		GroupID: "service",
	}

	// Add commands
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	// Subcommands
	cmd.AddCommand(role.NewRoleRootCmd())

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)
	
	return cmd
}
