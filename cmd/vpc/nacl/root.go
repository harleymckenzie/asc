package nacl

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"nacls", "nacl", "network-acl", "network-acls"}

// NewNACLRootCmd returns the root command for NACL operations.
func NewNACLRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nacl",
		Short:   "Perform Network ACL operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	cmd.AddCommand(lsCmd)
	// cmd.AddCommand(showCmd) // Disabled - show.go.disabled

	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)

	return cmd
}
