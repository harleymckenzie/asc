// The add command acts as an umbrella command for all add commands.
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
	addCmd.AddCommand(scheduleAddCmd)
	addCmd.AddGroup(cmdutil.SubcommandGroups()...)
	schedule.NewAddFlags(scheduleAddCmd)
}

// Command variable
var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add scheduled actions to an Auto Scaling Group",
	GroupID: "actions",
	Run:     func(cobraCmd *cobra.Command, args []string) {},
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
