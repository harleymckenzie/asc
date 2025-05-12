// add.go defines the 'add' subcommand for schedule operations.
package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/harleymckenzie/asc/pkg/service/asg"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/asg/types"
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
	"github.com/harleymckenzie/asc/pkg/shared/format"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/harleymckenzie/asc/pkg/shared/utils"
	"github.com/spf13/cobra"
)

var (
	asgName         string
	minSize         int
	maxSize         int
	desiredCapacity int
	recurrence      string
	startTimeStr    string
	endTimeStr      string
	tableOpts       tableformat.RenderOptions
)

// addCmd defines the 'add' subcommand for schedule operations.
var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a scheduled action to an Auto Scaling Group",
	Example: "asc asg add schedule my-schedule --asg-name my-asg --min-size 4 --start-time 'Friday 10:00'\n" +
		"asc asg add schedule my-schedule --asg-name my-asg --desired-capacity 8 --start-time '10:00am 25/04/2025'",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(AddSchedule(cobraCmd, args))
	},
}

// NewAddFlags adds flags for the add subcommand.
func NewAddFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().
		StringVarP(&asgName, "asg-name", "a", "", "The name of the Auto Scaling Group to add the schedule to.")
	cobraCmd.MarkFlagRequired("asg-name")
	cobraCmd.Flags().
		IntVarP(&minSize, "min-size", "m", 0, "The minimum size of the Auto Scaling Group.")
	cobraCmd.Flags().
		IntVarP(&maxSize, "max-size", "M", 0, "The maximum size of the Auto Scaling Group.")
	cobraCmd.Flags().
		IntVarP(&desiredCapacity, "desired-capacity", "d", 0, "The desired capacity of the Auto Scaling Group.")
	cobraCmd.Flags().
		StringVarP(&recurrence, "recurrence", "R", "", "The recurrence of the schedule.")
	cobraCmd.Flags().
		StringVarP(&startTimeStr, "start-time", "s", "", "The start time of the schedule.")
	cobraCmd.Flags().StringVarP(&endTimeStr, "end-time", "e", "", "The end time of the schedule.")
}

func init() {
	NewAddFlags(addCmd)
}

//
// Command functions
//

// AddSchedule is the handler for the add schedule subcommand.
func AddSchedule(cobraCmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	svc, err := asg.NewAutoScalingService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new Auto Scaling Group service: %w", err)
	}

	var startTime *time.Time
	var endTime *time.Time

	scheduledActionName := args[0]

	if startTimeStr != "" {
		t, err := format.ParseTime(startTimeStr)
		if err != nil {
			return fmt.Errorf("parse start time: %w", err)
		}
		startTime = &t
	}

	if endTimeStr != "" {
		t, err := format.ParseTime(endTimeStr)
		if err != nil {
			return fmt.Errorf("parse end time: %w", err)
		}
		endTime = &t
	}

	input := &ascTypes.AddAutoScalingGroupScheduleInput{
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

	err = svc.AddAutoScalingGroupSchedule(ctx, input)
	if err != nil {
		return fmt.Errorf("add schedule: %w", err)
	}

	// List the schedule and print it
	schedules, err := svc.GetAutoScalingGroupSchedules(
		ctx,
		&ascTypes.GetAutoScalingGroupSchedulesInput{
			AutoScalingGroupName: asgName,
			ScheduledActionNames: []string{scheduledActionName},
		},
	)
	if err != nil {
		return fmt.Errorf("get schedule: %w", err)
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(schedules),
		Fields:    asgScheduleFields(),
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return asg.GetScheduleAttributeValue(fieldID, instance)
		},
	}, tableOpts)
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}
	return nil
}
