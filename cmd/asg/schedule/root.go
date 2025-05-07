package schedule

import (
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
	"github.com/spf13/cobra"
)

// ScheduleCmd is the root command for schedule subcommands.
func NewScheduleRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "schedule",
		Aliases: []string{"schedules"},
		Short:   "Manage scheduled actions for an Auto Scaling Group",
		GroupID: "subcommands",
	}

	// Add the subcommands to the command
	cmd.AddCommand(addCmd)
	cmd.AddCommand(rmCmd)
	cmd.AddCommand(lsCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
