package igw

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"igws", "igw", "internet-gateways", "internet-gateway"}

// Root command
func NewIGWRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "igw",
		Short:   "Perform internet gateway operations",
		Aliases: CmdAliases,
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
