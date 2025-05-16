// The add command acts as an umbrella command for all add commands.
// It re-uses existing functions and flags from the relevant commands.

package asg

import (
	"github.com/spf13/cobra"
)

// Variables
//
// (No variables for this command)
//
// Init function
func init() {
}

// Command variable
var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add scheduled actions to an Auto Scaling Group",
	GroupID: "actions",
	Run:     func(cobraCmd *cobra.Command, args []string) {},
}
