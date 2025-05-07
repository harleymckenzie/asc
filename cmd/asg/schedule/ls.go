// ls.go defines the 'ls' subcommand for schedule operations.
package schedule

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/asg"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/asg/types"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

var (
	list bool

	sortName bool
)

// lsCmd is the cobra command for listing scheduled actions.
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all scheduled actions for Auto Scaling Groups",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ListSchedules(cobraCmd, args)
	},
}

// ListSchedules is the handler for the ls subcommand.
func ListSchedules(cmd *cobra.Command, args []string) {
	ctx := context.TODO()
	profile, _ := cmd.Root().PersistentFlags().GetString("profile")
	region, _ := cmd.Root().PersistentFlags().GetString("region")

	svc, err := asg.NewAutoScalingService(ctx, profile, region)
	if err != nil {
		log.Fatalf("Failed to initialize Auto Scaling Group service: %v", err)
	}

	if len(args) > 0 {
		ListSchedulesForGroup(svc, args[0])
	} else {
		ListSchedulesForAllGroups(svc)
	}
}

// ListSchedulesForGroup lists all schedules for a given Auto Scaling Group.
func ListSchedulesForGroup(svc *asg.AutoScalingService, asgName string) {
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

// ListSchedulesForAllGroups lists all schedules for all Auto Scaling Groups.
func ListSchedulesForAllGroups(svc *asg.AutoScalingService) {
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

// addLsFlags adds flags for the ls subcommand.
func addLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Auto-Scaling Groups in list format.")
	cobraCmd.Flags().SortFlags = false
}

func init() {
	addLsFlags(lsCmd)
}
