package role

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"role", "roles"}

// Root command
func NewRoleRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "role",
		Short:   "Perform role operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	// Add the subcommands to the command
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
