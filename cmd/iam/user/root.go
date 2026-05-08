package user

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"user", "users"}

func NewUserRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "Perform IAM user operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
