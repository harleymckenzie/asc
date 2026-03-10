package taskdefinition

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"task-definition", "task-definitions", "taskdefinitions", "td"}

// NewTaskDefinitionRootCmd creates and configures the root command for ECS task definition operations
func NewTaskDefinitionRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "task-definition",
		Short:   "Perform ECS task definition operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
