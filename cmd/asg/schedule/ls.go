// ls.go defines the 'ls' subcommand for schedule operations.
package schedule

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/asg"
	ascTypes "github.com/harleymckenzie/asc/internal/service/asg/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list                bool
	sortName            bool
	sortStartTime       bool
	sortEndTime         bool
	sortDesiredCapacity bool
	sortMinSize         bool
	sortMaxSize         bool

	reverseSort bool
)

func init() {
	NewLsFlags(lsCmd)
}

// Define columns for schedules
func asgScheduleFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Auto Scaling Group", Display: true, Sort: false, Merge: true},
		{ID: "Name", Display: true, Sort: sortName},
		{ID: "Recurrence", Display: true, Sort: false},
		{ID: "Start Time", Display: true, Sort: sortStartTime, DefaultSort: true},
		{ID: "End Time", Display: true, Sort: sortEndTime},
		{ID: "Desired Capacity", Display: true, Sort: sortDesiredCapacity},
		{ID: "Min", Display: true, Sort: sortMinSize, SortDirection: "desc"},
		{ID: "Max", Display: true, Sort: sortMaxSize, SortDirection: "desc"},
	}
}

// lsCmd is the cobra command for listing scheduled actions.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all scheduled actions for Auto Scaling Groups",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListSchedules(cobraCmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().
		BoolVarP(&list, "list", "l", false, "Outputs Auto-Scaling Groups in list format.")
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending ASG name.")
	cobraCmd.Flags().
		BoolVarP(&sortStartTime, "sort-start-time", "t", false, "Sort by descending start time (most recently started first).")
	cobraCmd.Flags().
		BoolVarP(&sortEndTime, "sort-end-time", "e", false, "Sort by descending end time (most recently ended first).")
	cobraCmd.Flags().
		BoolVarP(&sortDesiredCapacity, "sort-desired-capacity", "d", false, "Sort by descending desired capacity (most frequent first).")
	cobraCmd.Flags().
		BoolVarP(&sortMinSize, "sort-min-size", "m", false, "Sort by descending min size (most frequent first).")
	cobraCmd.Flags().
		BoolVarP(&sortMaxSize, "sort-max-size", "M", false, "Sort by descending max size (most frequent first).")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cobraCmd.MarkFlagsMutuallyExclusive(
		"sort-name",
		"sort-start-time",
		"sort-end-time",
		"sort-min-size",
		"sort-max-size",
		"sort-desired-capacity",
	)
}

//
// Command functions
//

// ListSchedules is the handler for the ls subcommand.
func ListSchedules(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)

	svc, err := asg.NewAutoScalingService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new Auto Scaling Group service: %w", err)
	}

	if len(args) > 0 {
		return ListSchedulesForGroup(svc, args[0])
	} else {
		return ListSchedulesForAllGroups(svc)
	}
}

// ListSchedulesForGroup lists all schedules for a given Auto Scaling Group.
func ListSchedulesForGroup(svc *asg.AutoScalingService, asgName string) error {
	ctx := context.TODO()
	schedules, err := svc.GetAutoScalingGroupSchedules(
		ctx,
		&ascTypes.GetAutoScalingGroupSchedulesInput{
			AutoScalingGroupName: asgName,
		},
	)
	if err != nil {
		return fmt.Errorf("get schedules for Auto Scaling Group %s: %w", asgName, err)
	}

	fields := asgScheduleFields()

	// Set "Auto Scaling Group" field Visible to false when listing for a single group
	fields[0].Display = false

	opts := tableformat.RenderOptions{
		Title:  fmt.Sprintf("Scheduled Actions\n(%s)", asgName),
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(schedules),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return asg.GetScheduleAttributeValue(fieldID, instance)
		},
	}, opts)
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}
	return nil
}

// ListSchedulesForAllGroups lists all schedules for all Auto Scaling Groups.
func ListSchedulesForAllGroups(svc *asg.AutoScalingService) error {
	ctx := context.TODO()
	schedules, err := svc.GetAutoScalingGroupSchedules(
		ctx,
		&ascTypes.GetAutoScalingGroupSchedulesInput{},
	)
	if err != nil {
		return fmt.Errorf("get schedules for all Auto Scaling Groups: %w", err)
	}

	fields := asgScheduleFields()

	opts := tableformat.RenderOptions{
		Title:  "Scheduled Actions",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(schedules),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return asg.GetScheduleAttributeValue(fieldID, instance)
		},
	}, opts)
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}
	return nil
}
