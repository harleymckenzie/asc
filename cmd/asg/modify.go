// The modify command allows updating min, max, or desired capacity for an Auto Scaling Group.

package asg

import (
	"context"
	"fmt"
	"time"

	"github.com/harleymckenzie/asc/internal/service/asg"
	ascTypes "github.com/harleymckenzie/asc/internal/service/asg/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

// Variables
var (
	minSizeStr         string
	maxSizeStr         string
	desiredCapacityStr string
	durationStr        string
)

// Init function
func init() {
	addModifyFlags(modifyCmd)
}

func asgScheduleFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Name", Display: true, Sort: false},
		{ID: "Recurrence", Display: true, Sort: false},
		{ID: "Start Time", Display: true, Sort: true, DefaultSort: true},
		{ID: "End Time", Display: true, Sort: false},
		{ID: "Desired Capacity", Display: true, Sort: false},
		{ID: "Min", Display: true, Sort: false},
		{ID: "Max", Display: true, Sort: false},
	}
}

// Command variable
var modifyCmd = &cobra.Command{
	Use:     "modify",
	Short:   "Modify an Auto Scaling Group min, max, or desired capacity",
	Long:    "Modify an Auto Scaling Group min, max, or desired capacity",
	Args:    cobra.ExactArgs(1),
	GroupID: "actions",
	Aliases: []string{"edit", "update"},
	Example: "  asc asg modify my-asg --min 3       # Set the minimum capacity to 3\n" +
		"  asc asg modify my-asg --max -6           # Decrease the maximum capacity by 6\n" +
		"  asc asg modify my-asg --desired +5       # Increase the desired capacity by 5\n" +
		"  asc asg modify my-asg -m +5 --duration 2h  # Increase the minimum capacity by 5 and revert the change after 2 hours",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ModifyAutoScalingGroup(cmd, args))
	},
}

// Flag function
func addModifyFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().SortFlags = false
	cobraCmd.Flags().StringVarP(&minSizeStr, "min", "m", "", "The minimum capacity (absolute or relative, e.g. 3, +1, -2)")
	cobraCmd.Flags().StringVarP(&maxSizeStr, "max", "M", "", "The maximum capacity (absolute or relative, e.g. 3, +3, -3)")
	cobraCmd.Flags().StringVarP(&desiredCapacityStr, "desired", "d", "", "The desired capacity (absolute or relative, e.g. 3, +1, -2)")
	cobraCmd.Flags().StringVarP(&durationStr, "duration", "D", "", "Duration after which to revert the changes (e.g. '2h', '30m', '1d')")
}

// Command functions
func ModifyAutoScalingGroup(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	profile, region := cmdutil.GetPersistentFlags(cmd)

	svc, err := asg.NewAutoScalingService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new Auto Scaling Service: %w", err)
	}

	// Get current information about the Auto Scaling Group
	getInput := &ascTypes.GetAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{args[0]},
	}
	asgOutput, err := svc.GetAutoScalingGroups(ctx, getInput)
	if err != nil {
		return fmt.Errorf("get Auto Scaling Groups: %w", err)
	}

	// Create a ModifyAutoScalingGroupInput struct to be updated with the new information
	input := &ascTypes.ModifyAutoScalingGroupInput{
		AutoScalingGroupName: args[0],
	}

	// Apply the relative or absolute values to the ModifyAutoScalingGroupInput struct
	if minSizeStr != "" {
		minSizeInt32, err := utils.ApplyRelativeOrAbsolute(minSizeStr, *asgOutput[0].MinSize)
		if err != nil {
			return fmt.Errorf("apply relative or absolute min size: %w", err)
		}
		input.MinSize = &minSizeInt32
	}
	if maxSizeStr != "" {
		maxSizeInt32, err := utils.ApplyRelativeOrAbsolute(maxSizeStr, *asgOutput[0].MaxSize)
		if err != nil {
			return fmt.Errorf("apply relative or absolute max size: %w", err)
		}
		input.MaxSize = &maxSizeInt32
	}
	if desiredCapacityStr != "" {
		desiredCapacityInt32, err := utils.ApplyRelativeOrAbsolute(
			desiredCapacityStr,
			*asgOutput[0].DesiredCapacity,
		)
		if err != nil {
			return fmt.Errorf("apply relative or absolute desired capacity: %w", err)
		}
		input.DesiredCapacity = &desiredCapacityInt32
	}

	// Modify the Auto Scaling Group
	err = svc.ModifyAutoScalingGroup(ctx, input)
	if err != nil {
		return fmt.Errorf("modify Auto Scaling Group: %w", err)
	}

	// Add a schedule to revert the change after a given duration
	if durationStr != "" {
		// Convert the duration string into time.Time, which will be passed to addRevertSchedule
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			return fmt.Errorf("parse duration: %w", err)
		}
		timeToRevert := time.Now().Add(duration)

		fmt.Printf("Creating scheduled action to revert changes on %s\n", timeToRevert.Format("Monday January 2 2006 at 15:04:05"))
		addScheduleInput := &ascTypes.AddAutoScalingGroupScheduleInput{
			AutoScalingGroupName: args[0],
			ScheduledActionName:  fmt.Sprintf("temporary-scaling-change-%s", time.Now().Format("2006-01-02-15-04-05")),
			StartTime:            &timeToRevert,
		}

		if minSizeStr != "" {
			addScheduleInput.MinSize = asgOutput[0].MinSize
		}

		if maxSizeStr != "" {
			addScheduleInput.MaxSize = asgOutput[0].MaxSize
		}

		if desiredCapacityStr != "" {
			addScheduleInput.DesiredCapacity = asgOutput[0].DesiredCapacity
		}

		err = addRevertSchedule(ctx, svc, addScheduleInput)
		if err != nil {
			return fmt.Errorf("add scheduled revert action: %w", err)
		}
	}

	return nil
}


func addRevertSchedule(ctx context.Context, svc *asg.AutoScalingService, input *ascTypes.AddAutoScalingGroupScheduleInput) error {
	err := svc.AddAutoScalingGroupSchedule(ctx, input)
	if err != nil {
		return fmt.Errorf("add scheduled revert action: %w", err)
	}

	schedules, err := svc.GetAutoScalingGroupSchedules(ctx, &ascTypes.GetAutoScalingGroupSchedulesInput{
		AutoScalingGroupName: input.AutoScalingGroupName,
		ScheduledActionNames: []string{input.ScheduledActionName},
	})
	if err != nil {
		return fmt.Errorf("get scheduled revert action: %w", err)
	}

	opts := tableformat.RenderOptions{
		Title:  "Scheduled Revert Action",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(asgScheduleFields(), false),
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(schedules),
		Fields:    asgScheduleFields(),
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return asg.GetScheduleAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}