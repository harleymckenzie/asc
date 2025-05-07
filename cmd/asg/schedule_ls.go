package asg

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/asg"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/asg/types"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

var (
	sortStartTime bool
	sortEndTime bool
	sortMinSize bool
	sortMaxSize bool
)

func lsSchedules(cobraCmd *cobra.Command, args []string) {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	svc, err := asg.NewAutoScalingService(ctx, profile, region)
	if err != nil {
		log.Fatalf("Failed to initialize Auto Scaling Group service: %v", err)
	}

	if len(args) > 0 {
		ListAutoScalingGroupSchedules(svc, args[0])
	} else {
		ListAllAutoScalingGroupSchedules(svc)
	}
}

// ListAutoScalingGroupSchedules lists all schedules for a given Auto Scaling Group
func ListAutoScalingGroupSchedules(svc *asg.AutoScalingService, asgName string) {
	ctx := context.TODO()
	schedules, err := svc.GetAutoScalingGroupSchedules(ctx, &ascTypes.GetAutoScalingGroupSchedulesInput{
		AutoScalingGroupName: asgName,
	})
	if err != nil {
		log.Fatalf("Failed to get schedules for Auto Scaling Group %s: %v", asgName, err)
	}

	// Define columns for schedules
	columns := []tableformat.Column{
		{ID: "Name", Visible: true, Sort: sortName},
		{ID: "Recurrence", Visible: true, Sort: false},
		{ID: "Start Time", Visible: true, Sort: sortStartTime, DefaultSort: true},
		{ID: "End Time", Visible: true, Sort: sortEndTime},
		{ID: "Desired Capacity", Visible: true, Sort: sortDesiredCapacity},
		{ID: "Min", Visible: true, Sort: sortMinSize},
		{ID: "Max", Visible: true, Sort: sortMaxSize},
	}
	selectedColumns, sortBy := tableformat.BuildColumns(columns)

	opts := tableformat.RenderOptions{
		SortBy: sortBy,
		List:   list,
	}

	tableformat.Render(&asg.AutoScalingSchedulesTable{
		Schedules:       schedules,
		SelectedColumns: selectedColumns,
	}, opts)
}

// ListAllAutoScalingGroupSchedules lists all schedules for all Auto Scaling Groups
func ListAllAutoScalingGroupSchedules(svc *asg.AutoScalingService) {
	ctx := context.TODO()
	schedules, err := svc.GetAutoScalingGroupSchedules(ctx, &ascTypes.GetAutoScalingGroupSchedulesInput{})
	if err != nil {
		log.Fatalf("Failed to get schedules for all Auto Scaling Groups: %v", err)
	}

	// Define columns for schedules
	columns := []tableformat.Column{
		{ID: "Auto Scaling Group", Visible: true, Sort: false},
		{ID: "Name", Visible: true, Sort: sortName},
		{ID: "Recurrence", Visible: true, Sort: false},
		{ID: "Start Time", Visible: true, Sort: false},
		{ID: "End Time", Visible: true, Sort: false},
		{ID: "Desired Capacity", Visible: true, Sort: false},
		{ID: "Min", Visible: true, Sort: false},
		{ID: "Max", Visible: true, Sort: false},
	}
	selectedColumns, sortBy := tableformat.BuildColumns(columns)

	opts := tableformat.RenderOptions{
		SortBy: sortBy,
		List:   list,
	}

	tableformat.Render(&asg.AutoScalingSchedulesTable{
		Schedules:       schedules,
		SelectedColumns: selectedColumns,
	}, opts)
}

func addScheduleLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Auto-Scaling Groups in list format.")
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", true, "Sort by descending ASG name.")
	cobraCmd.Flags().BoolVarP(&sortStartTime, "sort-start-time", "t", false, "Sort by descending start time (most recently started first).")
	cobraCmd.Flags().BoolVarP(&sortEndTime, "sort-end-time", "e", false, "Sort by descending end time (most recently ended first).")
	cobraCmd.Flags().BoolVarP(&sortDesiredCapacity, "sort-desired-capacity", "d", false, "Sort by descending desired capacity (most frequent first).")
	cobraCmd.Flags().BoolVarP(&sortMinSize, "sort-min-size", "m", false, "Sort by descending min size (most frequent first).")
	cobraCmd.Flags().BoolVarP(&sortMaxSize, "sort-max-size", "M", false, "Sort by descending max size (most frequent first).")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-start-time", "sort-end-time", "sort-min-size", "sort-max-size", "sort-desired-capacity")
}
