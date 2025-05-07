// The ls command list Auto Scaling Groups, as well as an alias for the relevant subcommand.
// It re-uses existing functions and flags from the relevant commands.

package asg

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/asg"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/asg/types"
	"github.com/harleymckenzie/asc/cmd/asg/schedule"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
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
)

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

func newLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Auto-Scaling Groups in list format.")
	cobraCmd.Flags().BoolVar(&showARNs, "arn", false, "Show ARNs for each Auto-Scaling Group.")
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", true, "Sort by descending ASG name.")
	cobraCmd.Flags().BoolVarP(&sortInstances, "sort-instances", "i", false, "Sort by descending number of instances.")
	cobraCmd.Flags().BoolVarP(&sortDesiredCapacity, "sort-desired-capacity", "d", false, "Sort by descending desired capacity.")
	cobraCmd.Flags().BoolVarP(&sortMinCapacity, "sort-min-capacity", "m", false, "Sort by descending min capacity.")
	cobraCmd.Flags().BoolVarP(&sortMaxCapacity, "sort-max-capacity", "M", false, "Sort by descending max capacity.")
}

var scheduleLsCmd = &cobra.Command{
	Use:   "schedules",
	Short: "List schedules for an Auto Scaling Group",
	GroupID: "subcommands",
	Run: func(cobraCmd *cobra.Command, args []string) {
		schedule.ListSchedules(cobraCmd, args)
	},
}

func init() {
	newLsFlags(lsCmd)

	// Add the lsSchedulesCmd to the lsCmd
	lsCmd.AddCommand(scheduleLsCmd)
	lsCmd.AddGroup(cmdutil.SubcommandGroups()...)

	// Add the lsSchedulesCmd to the lsCmd
	schedule.NewLsFlags(scheduleLsCmd)
}

func ListAutoScalingGroups(cobraCmd *cobra.Command, args []string) {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	svc, err := asg.NewAutoScalingService(ctx, profile, region)
	if err != nil {
		log.Fatalf("Failed to initialize Auto Scaling Group service: %v", err)
	}

	if len(args) > 0 {
		ListAutoScalingGroupInstances(svc, args[0])
	} else {
		autoScalingGroups, err := svc.GetAutoScalingGroups(ctx, &ascTypes.GetAutoScalingGroupsInput{
			AutoScalingGroupNames: []string{},
		})
		if err != nil {
			log.Fatalf("Failed to get Auto Scaling Groups: %v", err)
		}

		// Define columns for Auto Scaling Groups
		columns := []tableformat.Column{
			{ID: "Name", Visible: true, Sort: sortName},
			{ID: "Instances", Visible: true, Sort: sortInstances},
			{ID: "Desired", Visible: true, Sort: sortDesiredCapacity},
			{ID: "Min", Visible: true, Sort: sortMinCapacity},
			{ID: "Max", Visible: true, Sort: sortMaxCapacity},
		}
		selectedColumns, sortBy := tableformat.BuildColumns(columns)

		opts := tableformat.RenderOptions{
			SortBy: sortBy,
			List:   list,
			Title:  "Auto Scaling Groups",
		}

		tableformat.Render(&asg.AutoScalingTable{
			AutoScalingGroups: autoScalingGroups,
			SelectedColumns:   selectedColumns,
		}, opts)
	}
}

func ListAutoScalingGroupInstances(svc *asg.AutoScalingService, asgName string) {
	ctx := context.TODO()
	instances, err := svc.GetAutoScalingGroupInstances(ctx, &ascTypes.GetAutoScalingGroupInstancesInput{
		AutoScalingGroupNames: []string{asgName},
	})
	if err != nil {
		log.Fatalf("Failed to get instances for Auto Scaling Group %s: %v", asgName, err)
	}

	// Define columns for instances
	columns := []tableformat.Column{
		{ID: "Name", Visible: true, Sort: sortName},
		{ID: "State", Visible: true, Sort: false},
		{ID: "Instance Type", Visible: true, Sort: false},
		{ID: "Launch Template/Configuration", Visible: true, Sort: false},
		{ID: "Availability Zone", Visible: true, Sort: false},
		{ID: "Health", Visible: true, Sort: false},
	}
	selectedColumns, sortBy := tableformat.BuildColumns(columns)

	opts := tableformat.RenderOptions{
		SortBy: sortBy,
		List:   list,
	}

	tableformat.Render(&asg.AutoScalingInstanceTable{
		Instances:       instances,
		SelectedColumns: selectedColumns,
	}, opts)
}
