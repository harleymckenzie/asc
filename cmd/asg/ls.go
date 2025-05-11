// The ls command list Auto Scaling Groups, as well as an alias for the relevant subcommand.
// It re-uses existing functions and flags from the relevant commands.

package asg

import (
	"context"
	"fmt"
	"log"

	"github.com/harleymckenzie/asc/cmd/asg/schedule"
	"github.com/harleymckenzie/asc/pkg/service/asg"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/asg/types"
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/harleymckenzie/asc/pkg/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list bool

	showARNs bool

	sortName            bool
	sortInstances       bool
	sortDesiredCapacity bool
	sortMinCapacity     bool
	sortMaxCapacity     bool

	reverseSort bool
)

func init() {
	newLsFlags(lsCmd)

	// Add the lsSchedulesCmd to the lsCmd
	lsCmd.AddCommand(scheduleLsCmd)
	lsCmd.AddGroup(cmdutil.SubcommandGroups()...)

	// Add the lsSchedulesCmd to the lsCmd
	schedule.NewLsFlags(scheduleLsCmd)
}

//
// Column functions
//

func asgFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Name", Visible: true, Sort: sortName},
		{ID: "Instances", Visible: true, Sort: sortInstances},
		{ID: "Desired", Visible: true, Sort: sortDesiredCapacity},
		{ID: "Min", Visible: true, Sort: sortMinCapacity},
		{ID: "Max", Visible: true, Sort: sortMaxCapacity},
		{ID: "ARN", Visible: showARNs, Sort: false},
	}
}

func asgInstanceFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Name", Visible: true, Sort: sortName, DefaultSort: true},
		{ID: "State", Visible: true, Sort: false},
		{ID: "Instance Type", Visible: true, Sort: false},
		{ID: "Launch Template/Configuration", Visible: true, Sort: false},
		{ID: "Availability Zone", Visible: true, Sort: false},
		{ID: "Health", Visible: true, Sort: false},
	}
}

// lsCmd is the command for listing Auto Scaling Groups, instances in an ASG, or schedules
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List Auto Scaling Groups, instances in an ASG, or schedules",
	Example: "  ls                      List all Auto Scaling Groups\n" +
		"  ls [asg-name]           List instances in the specified ASG\n" +
		"  ls schedules [asg-name] List schedules (optionally for specific ASG)",
	GroupID: "actions",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ListAutoScalingGroups(cobraCmd, args)
	},
}

// newLsFlags is the function for adding flags to the ls command
func newLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().
		BoolVarP(&list, "list", "l", false, "Outputs Auto-Scaling Groups in list format.")
	cobraCmd.Flags().
		BoolVarP(&showARNs, "arn", "a", false, "Show ARNs for each Auto-Scaling Group.")
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending ASG name.")
	cobraCmd.Flags().
		BoolVarP(&sortInstances, "sort-instances", "i", false, "Sort by descending number of instances. (ASG output only)")
	cobraCmd.Flags().
		BoolVarP(&sortDesiredCapacity, "sort-desired-capacity", "d", false, "Sort by descending desired capacity. (ASG output only)")
	cobraCmd.Flags().
		BoolVarP(&sortMinCapacity, "sort-min-capacity", "m", false, "Sort by descending min capacity. (ASG output only)")
	cobraCmd.Flags().
		BoolVarP(&sortMaxCapacity, "sort-max-capacity", "M", false, "Sort by descending max capacity. (ASG output only)")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
}

// scheduleLsCmd is the command for listing schedules for an Auto Scaling Group
var scheduleLsCmd = &cobra.Command{
	Use:     "schedules",
	Short:   "List schedules for an Auto Scaling Group",
	GroupID: "subcommands",
	Run: func(cobraCmd *cobra.Command, args []string) {
		schedule.ListSchedules(cobraCmd, args)
	},
}

//
// Command functions
//

func ListAutoScalingGroups(cobraCmd *cobra.Command, args []string) {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	svc, err := asg.NewAutoScalingService(ctx, profile, region)
	if err != nil {
		log.Fatalf("Failed to initialize Auto Scaling Group service: %v", err)
	}

	if len(args) > 0 {
		fmt.Printf("Listing instances for Auto Scaling Group %s\n", args[0])
		ListAutoScalingGroupInstances(svc, args[0])
	} else {
		autoScalingGroups, err := svc.GetAutoScalingGroups(ctx, &ascTypes.GetAutoScalingGroupsInput{})
		if err != nil {
			log.Fatalf("Failed to get Auto Scaling Groups: %v", err)
		}

		fields := asgFields()

		opts := tableformat.RenderOptions{
			Title:  "Auto Scaling Groups",
			Style:  "rounded",
			SortBy: tableformat.GetSortByField(fields, reverseSort),
		}

		if list {
			opts.Style = "list"
		}

		tableformat.RenderTableList(&tableformat.ListTable{
			Instances: utils.SlicesToAny(autoScalingGroups),
			Fields:    fields,
			GetAttribute: func(fieldID string, instance any) string {
				return asg.GetAttributeValue(fieldID, instance)
			},
		}, opts)
	}
}

// ListAutoScalingGroupInstances is the function for listing instances in an Auto Scaling Group
func ListAutoScalingGroupInstances(svc *asg.AutoScalingService, asgName string) {
	ctx := context.TODO()
	instances, err := svc.GetAutoScalingGroupInstances(
		ctx,
		&ascTypes.GetAutoScalingGroupInstancesInput{
			AutoScalingGroupNames: []string{asgName},
		},
	)
	if err != nil {
		log.Fatalf("Failed to get instances for Auto Scaling Group %s: %v", asgName, err)
	}

	// Define columns for instances
	fields := asgInstanceFields()

	opts := tableformat.RenderOptions{
		Title:  "Auto Scaling Group Instances",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(instances),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) string {
			return asg.GetInstanceAttributeValue(fieldID, instance)
		},
	}, opts)
}
