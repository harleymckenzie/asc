// The rm command acts as an umbrella command for all rm commands.
// It re-uses existing functions and flags from the relevant commands.

package asg

import (
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove scheduled actions from an Auto Scaling Group",
	Run: func(cobraCmd *cobra.Command, args []string) {},
	GroupID: "actions",
}

var rmScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Remove scheduled actions from an Auto Scaling Group",
	Run: func(cobraCmd *cobra.Command, args []string) {
		rmSchedule(cobraCmd, args)
	},
}

func init() {
	addScheduleRmFlags(rmScheduleCmd)
	rmCmd.AddCommand(rmScheduleCmd)
}