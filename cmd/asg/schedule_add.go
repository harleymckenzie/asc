package asg

import (
	"context"
	"fmt"
	"github.com/harleymckenzie/asc/pkg/service/asg"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/asg/types"
	"github.com/harleymckenzie/asc/pkg/shared/timeformat"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var (
	asgName         string
	minSize         int
	maxSize         int
	desiredCapacity int
	recurrence      string
	startTimeStr    string
	endTimeStr      string
)

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

func addSchedule(cobraCmd *cobra.Command, args []string) {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	svc, err := asg.NewAutoScalingService(ctx, profile, region)
	if err != nil {
		log.Fatalf("Failed to initialize Auto Scaling Group service: %v", err)
	}

	var startTime *time.Time
	var endTime *time.Time

	scheduledActionName := args[0]

	if startTimeStr != "" {
		t, err := timeformat.ParseTime(startTimeStr)
		if err != nil {
			log.Fatalf("Failed to parse start time: %v", err)
		}
		startTime = &t
		fmt.Println("startTime: ", startTime)
	}

	if endTimeStr != "" {
		t, err := timeformat.ParseTime(endTimeStr)
		if err != nil {
			log.Fatalf("Failed to parse end time: %v", err)
		}
		endTime = &t
		fmt.Println("endTime: ", endTime)
	}

	input := &ascTypes.AddScheduleInput{
		ScheduledActionName:  scheduledActionName,
		AutoScalingGroupName: asgName,
	}

	if minSize != 0 {
		minSizeInt32 := int32(minSize)
		input.MinSize = &minSizeInt32
	}

	if maxSize != 0 {
		maxSizeInt32 := int32(maxSize)
		input.MaxSize = &maxSizeInt32
	}

	if desiredCapacity != 0 {
		desiredCapacityInt32 := int32(desiredCapacity)
		input.DesiredCapacity = &desiredCapacityInt32
	}

	if recurrence != "" {
		input.Recurrence = &recurrence
	}

	if startTime != nil {
		input.StartTime = startTime
	}

	if endTime != nil {
		input.EndTime = endTime
	}

	err = svc.AddSchedule(ctx, input)
	if err != nil {
		log.Fatalf("Failed to add schedule: %v", err)
	}
}

func addScheduleAddFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().StringVarP(&asgName, "asg-name", "a", "", "The name of the Auto Scaling Group to add the schedule to.")
	cobraCmd.MarkFlagRequired("asg-name")
	cobraCmd.Flags().IntVarP(&minSize, "min-size", "m", 0, "The minimum size of the Auto Scaling Group.")
	cobraCmd.Flags().IntVarP(&maxSize, "max-size", "M", 0, "The maximum size of the Auto Scaling Group.")
	cobraCmd.Flags().IntVarP(&desiredCapacity, "desired-capacity", "d", 0, "The desired capacity of the Auto Scaling Group.")
	cobraCmd.Flags().StringVarP(&recurrence, "recurrence", "R", "", "The recurrence of the schedule.")
	cobraCmd.Flags().StringVarP(&startTimeStr, "start-time", "s", "", "The start time of the schedule.")
	cobraCmd.Flags().StringVarP(&endTimeStr, "end-time", "e", "", "The end time of the schedule.")
}

func init() {
	addScheduleAddFlags(scheduleAddCmd)
}
