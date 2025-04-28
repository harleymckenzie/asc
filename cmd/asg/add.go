// The add command acts as an umbrella command for all add commands.
// It re-uses existing functions and flags from the relevant commands.

package asg

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add scheduled actions to an Auto Scaling Group",
	Run: func(cobraCmd *cobra.Command, args []string) {},
	GroupID: "actions",
}

var addScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Add scheduled actions to an Auto Scaling Group",
	Run: func(cobraCmd *cobra.Command, args []string) {
		addSchedule(cobraCmd, args)
	},
}

func init() {
	addScheduleAddFlags(addScheduleCmd)
	addCmd.AddCommand(addScheduleCmd)
}