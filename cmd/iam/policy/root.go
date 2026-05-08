package policy

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"policy", "policies"}

func NewPolicyRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "policy",
		Short:   "Perform IAM policy operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
