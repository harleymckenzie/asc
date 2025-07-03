package iam

import (
	"github.com/harleymckenzie/asc/cmd/iam/role"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	// Add subcommands
	lsCmd.AddCommand(roleLsCmd)
	
	showCmd.AddCommand(roleShowCmd)

	// Add flags
	role.NewLsFlags(roleLsCmd)

	role.NewShowFlags(roleShowCmd)

	// Add groups
	lsCmd.AddGroup(cmdutil.SubcommandGroups()...)
	showCmd.AddGroup(cmdutil.SubcommandGroups()...)
}

// Subcommand variables
var roleLsCmd = &cobra.Command{
	Use:     "roles",
	Short:   "List roles",
	Aliases: role.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(role.ListRoles(cmd, args))
	},
}

var roleShowCmd = &cobra.Command{
	Use:     "role",
	Short:   "Show role",
	Aliases: role.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(role.ShowIAMRole(cmd, args[0]))
	},
}