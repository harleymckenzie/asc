package schedule

import (
	"github.com/spf13/cobra"
)

// ScheduleCmd is the root command for schedule subcommands.
var ScheduleCmd = &cobra.Command{
	Use:     "schedule",
	Short:   "Manage scheduled actions for an Auto Scaling Group",
	GroupID: "subcommands",
}

func init() {
	ScheduleCmd.AddCommand(addCmd)
	ScheduleCmd.AddCommand(rmCmd)
	ScheduleCmd.AddCommand(lsCmd)
}
