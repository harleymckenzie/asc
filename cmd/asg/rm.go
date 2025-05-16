// The rm command acts as an umbrella command for all rm commands.
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
func init() {}

// Command variable
var rmCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Remove scheduled actions from an Auto Scaling Group",
	GroupID: "actions",
	Run:     func(cobraCmd *cobra.Command, args []string) {},
}

