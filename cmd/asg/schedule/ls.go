// ls.go defines the 'ls' subcommand for schedule operations.
package schedule

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/asg"
	ascTypes "github.com/harleymckenzie/asc/internal/service/asg/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
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
func getScheduleFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Auto Scaling Group", Category: "Schedule", Visible: true},
		{Name: "Name", Category: "Schedule", Visible: true, SortBy: sortName, SortDirection: tablewriter.Asc},
		{Name: "Recurrence", Category: "Schedule", Visible: true},
		{Name: "Start Time", Category: "Schedule", Visible: true, DefaultSort: true, SortBy: sortStartTime, SortDirection: tablewriter.Desc},
		{Name: "End Time", Category: "Schedule", Visible: true, SortBy: sortEndTime, SortDirection: tablewriter.Desc},
		{Name: "Desired Capacity", Category: "Schedule", Visible: true, SortBy: sortDesiredCapacity, SortDirection: tablewriter.Desc},
		{Name: "Min", Category: "Schedule", Visible: true, SortBy: sortMinSize, SortDirection: tablewriter.Desc},
		{Name: "Max", Category: "Schedule", Visible: true, SortBy: sortMaxSize, SortDirection: tablewriter.Desc},
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
	svc, err := cmdutil.CreateService(cmd, asg.NewAutoScalingService)
	if err != nil {
		return fmt.Errorf("create new Auto Scaling Group service: %w", err)
	}

	ctx := cmd.Context()
	if len(args) > 0 {
		return ListSchedulesForGroup(ctx, svc, args[0])
	} else {
		return ListSchedulesForAllGroups(ctx, svc)
	}
}

// ListSchedulesForGroup lists all schedules for a given Auto Scaling Group.
func ListSchedulesForGroup(ctx context.Context, svc *asg.AutoScalingService, asgName string) error {
	schedules, err := svc.GetAutoScalingGroupSchedules(
		ctx,
		&ascTypes.GetAutoScalingGroupSchedulesInput{
			AutoScalingGroupName: asgName,
		},
	)
	if err != nil {
		return fmt.Errorf("get schedules for Auto Scaling Group %s: %w", asgName, err)
	}

	fields := getScheduleFields()
	// Set "Auto Scaling Group" field Visible to false when listing for a single group
	fields[0].Visible = false

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         fmt.Sprintf("Scheduled Actions\n(%s)", asgName),
		PlainStyle:    list,
		Fields:        fields,
		Data:          utils.SlicesToAny(schedules),
		GetFieldValue: asg.GetFieldValue,
		GetTagValue:   asg.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}

// ListSchedulesForAllGroups lists all schedules for all Auto Scaling Groups.
func ListSchedulesForAllGroups(ctx context.Context, svc *asg.AutoScalingService) error {
	schedules, err := svc.GetAutoScalingGroupSchedules(
		ctx,
		&ascTypes.GetAutoScalingGroupSchedulesInput{},
	)
	if err != nil {
		return fmt.Errorf("get schedules for all Auto Scaling Groups: %w", err)
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         "Scheduled Actions",
		PlainStyle:    list,
		Fields:        getScheduleFields(),
		Data:          utils.SlicesToAny(schedules),
		GetFieldValue: asg.GetFieldValue,
		GetTagValue:   asg.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}
