package iam

import (
	"github.com/harleymckenzie/asc/cmd/iam/policy"
	"github.com/harleymckenzie/asc/cmd/iam/role"
	"github.com/harleymckenzie/asc/cmd/iam/user"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewIAMRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "iam",
		Short:   "Perform IAM operations",
		GroupID: "service",
	}

	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	cmd.AddCommand(role.NewRoleRootCmd())
	cmd.AddCommand(policy.NewPolicyRootCmd())
	cmd.AddCommand(user.NewUserRootCmd())

	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)

	return cmd
}
