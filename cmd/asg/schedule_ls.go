package asg

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/asg"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/asg/types"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

var scheduleLsCmd = &cobra.Command{
	Use:   "schedules",
	Short: "List all schedules for an Auto Scaling Group",
	Run: func(cobraCmd *cobra.Command, args []string) {
		lsSchedule(cobraCmd, args)
	},
}

func lsSchedule(cobraCmd *cobra.Command, args []string) {
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
	schedules, err := svc.GetSchedules(ctx, &ascTypes.GetSchedulesInput{
		AutoScalingGroupName: asgName,
	})
	if err != nil {
		log.Fatalf("Failed to get schedules for Auto Scaling Group %s: %v", asgName, err)
	}

	// Define columns for schedules
	columns := []tableformat.Column{
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

// ListAllAutoScalingGroupSchedules lists all schedules for all Auto Scaling Groups
func ListAllAutoScalingGroupSchedules(svc *asg.AutoScalingService) {
	ctx := context.TODO()
	schedules, err := svc.GetSchedules(ctx, &ascTypes.GetSchedulesInput{})
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
	cobraCmd.Flags().SortFlags = false
}

func init() {
	addScheduleLsFlags(scheduleLsCmd)
}
