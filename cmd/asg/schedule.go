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

var scheduleLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List scheduled actions for an Auto Scaling Group",
	Run: func(cobraCmd *cobra.Command, args []string) {
		lsSchedules(cobraCmd, args)
	},
}

var scheduleAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a scheduled action to an Auto Scaling Group",
	Long:  "Add a scheduled action to an Auto Scaling Group\n",
	Run: func(cobraCmd *cobra.Command, args []string) {
		addSchedule(cobraCmd, args)
	},
	Example: "asc asg add schedule my-schedule --asg-name my-asg --min-size 4 --start-time 'Friday 10:00'\n" +
		"asc asg add schedule my-schedule --asg-name my-asg --desired-capacity 8 --start-time '10:00am 25/04/2025'",
}

var scheduleRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a scheduled action from an Auto Scaling Group",
    Example: "asc asg rm schedule my-schedule --asg-name my-asg",
	Run: func(cobraCmd *cobra.Command, args []string) {
		rmSchedule(cobraCmd, args)
	},
}

func init() {
	addScheduleLsFlags(scheduleLsCmd)
	addScheduleAddFlags(scheduleAddCmd)
	addScheduleRmFlags(scheduleRmCmd)
	scheduleCmd.AddCommand(scheduleAddCmd)
	scheduleCmd.AddCommand(scheduleRmCmd)
    scheduleCmd.AddCommand(scheduleLsCmd)
}
