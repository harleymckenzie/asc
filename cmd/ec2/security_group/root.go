package security_group

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// Root command
func NewSecurityGroupRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "security-group",
		Short:   "Perform security group operations",
		Aliases: []string{"security-groups", "sg", "secgroup"},
		GroupID: "subcommands",
	}

	// Add the subcommands to the command
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)

	return cmd
}
