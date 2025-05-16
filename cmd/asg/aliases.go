package asg

import (
	"github.com/harleymckenzie/asc/cmd/asg/schedule"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	// Add subcommands
    addCmd.AddCommand(scheduleAddCmd)
    lsCmd.AddCommand(scheduleLsCmd)
    rmCmd.AddCommand(scheduleRmCmd)

	// Add flags
    schedule.NewAddFlags(scheduleAddCmd)
    schedule.NewLsFlags(scheduleLsCmd)
	schedule.NewRmFlags(scheduleRmCmd)
	
    // Add groups
	addCmd.AddGroup(cmdutil.SubcommandGroups()...)
	lsCmd.AddGroup(cmdutil.SubcommandGroups()...)
	rmCmd.AddGroup(cmdutil.SubcommandGroups()...)
}

// Subcommand variable
var scheduleAddCmd = &cobra.Command{
	Use:     "schedule",
	Short:   "Add scheduled actions to an Auto Scaling Group",
	GroupID: "subcommands",
	Run: func(cobraCmd *cobra.Command, args []string) {
		schedule.AddSchedule(cobraCmd, args)
	},
}

// scheduleLsCmd is the command for listing schedules for an Auto Scaling Group
var scheduleLsCmd = &cobra.Command{
	Use:     "schedules",
	Short:   "List schedules for an Auto Scaling Group",
	GroupID: "subcommands",
	Run: func(cobraCmd *cobra.Command, args []string) {
		schedule.ListSchedules(cobraCmd, args)
	},
}

// scheduleRmCmd is the command for removing schedules for an Auto Scaling Group
var scheduleRmCmd = &cobra.Command{
	Use:     "schedule",
	Short:   "Remove scheduled actions from an Auto Scaling Group",
	GroupID: "subcommands",
	Run: func(cobraCmd *cobra.Command, args []string) {
		schedule.RemoveSchedule(cobraCmd, args)
	},
}
