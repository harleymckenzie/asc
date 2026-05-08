package iam

import (
	"github.com/harleymckenzie/asc/cmd/iam/policy"
	"github.com/harleymckenzie/asc/cmd/iam/role"
	"github.com/harleymckenzie/asc/cmd/iam/user"

	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	lsCmd.AddCommand(roleLsCmd)
	lsCmd.AddCommand(policyLsCmd)
	lsCmd.AddCommand(userLsCmd)

	showCmd.AddCommand(roleShowCmd)
	showCmd.AddCommand(policyShowCmd)
	showCmd.AddCommand(userShowCmd)

	role.NewLsFlags(roleLsCmd)
	policy.NewLsFlags(policyLsCmd)
	user.NewLsFlags(userLsCmd)

	role.NewShowFlags(roleShowCmd)
	policy.NewShowFlags(policyShowCmd)
	user.NewShowFlags(userShowCmd)

	lsCmd.AddGroup(cmdutil.SubcommandGroups()...)
	showCmd.AddGroup(cmdutil.SubcommandGroups()...)
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List IAM resources",
	Aliases: []string{"list"},
	GroupID: "actions",
}

var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an IAM resource",
	Aliases: []string{"describe"},
	GroupID: "actions",
}

var roleLsCmd = &cobra.Command{
	Use:     "roles",
	Short:   "List IAM roles",
	Aliases: role.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(role.ListRoles(cmd, args))
	},
}

var policyLsCmd = &cobra.Command{
	Use:     "policies",
	Short:   "List IAM policies",
	Aliases: policy.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(policy.ListPolicies(cmd, args))
	},
}

var userLsCmd = &cobra.Command{
	Use:     "users",
	Short:   "List IAM users",
	Aliases: user.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(user.ListUsers(cmd, args))
	},
}

var roleShowCmd = &cobra.Command{
	Use:     "roles",
	Short:   "Show detailed information about an IAM role",
	Aliases: role.CmdAliases,
	GroupID: "subcommands",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(role.ShowRole(cmd, args[0]))
	},
}

var policyShowCmd = &cobra.Command{
	Use:     "policies",
	Short:   "Show detailed information about an IAM policy",
	Aliases: policy.CmdAliases,
	GroupID: "subcommands",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(policy.ShowPolicy(cmd, args[0]))
	},
}

var userShowCmd = &cobra.Command{
	Use:     "users",
	Short:   "Show detailed information about an IAM user",
	Aliases: user.CmdAliases,
	GroupID: "subcommands",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(user.ShowUser(cmd, args[0]))
	},
}
