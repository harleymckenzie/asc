// The rm command acts as an umbrella command for all rm commands.
// It re-uses existing functions and flags from the relevant commands.

package asg

import (
	schedule "github.com/harleymckenzie/asc/cmd/asg/schedule"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// Variables
//
// (No variables for this command)
//
// Init function
func init() {
	rmCmd.AddCommand(scheduleRmCmd)
	rmCmd.AddGroup(cmdutil.SubcommandGroups()...)
	schedule.NewRmFlags(scheduleRmCmd)
}

// Command variable
var rmCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Remove scheduled actions from an Auto Scaling Group",
	GroupID: "actions",
	Run:     func(cobraCmd *cobra.Command, args []string) {},
}

// Subcommand variable
var scheduleRmCmd = &cobra.Command{
	Use:     "schedule",
	Short:   "Remove scheduled actions from an Auto Scaling Group",
	GroupID: "subcommands",
	Run: func(cobraCmd *cobra.Command, args []string) {
		schedule.RemoveSchedule(cobraCmd, args)
	},
}
