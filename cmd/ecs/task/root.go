package task

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"tasks"}

// NewTaskRootCmd creates and configures the root command for ECS task operations
func NewTaskRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "task",
		Short:   "Perform ECS task operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
