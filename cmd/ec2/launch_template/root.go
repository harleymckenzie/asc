package launch_template

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"launch-templates", "lt"}

// Root command
func NewLaunchTemplateRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "launch-template",
		Short:   "Perform Launch Template operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	// Add the subcommands to the command
	cmd.AddCommand(showCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
