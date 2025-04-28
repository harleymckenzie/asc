package asg

import (
	"log"
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage scheduled actions for an Auto Scaling Group",
	GroupID: "subcommands",
	Run: func(cobraCmd *cobra.Command, args []string) {
        
		if len(args) == 0 {
			log.Fatalf("Please provide a subcommand: add or rm")
		}

		switch args[0] {
		case "add":
			scheduleAddCmd.Run(cobraCmd, args[1:])
		case "rm":
			scheduleRmCmd.Run(cobraCmd, args[1:])
		case "ls":
			scheduleLsCmd.Run(cobraCmd, args[1:])
		default:
			log.Fatalf("Invalid subcommand: %s", args[0])
		}
	},
}

func init() {
	scheduleCmd.AddCommand(scheduleAddCmd)
	scheduleCmd.AddCommand(scheduleRmCmd)
    scheduleCmd.AddCommand(scheduleLsCmd)
}
