// The ls command list Auto Scaling Groups, as well as an alias for the relevant subcommand.
// It re-uses existing functions and flags from the relevant commands.

package asg

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
}

//
// Column functions
//

func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Auto Scaling Group", Visible: true, SortBy: true, SortDirection: tablewriter.Asc},
		{Name: "Instances", Category: "Auto Scaling Group", Visible: true, SortBy: sortInstances, SortDirection: tablewriter.Desc},
		{Name: "Desired", Category: "Auto Scaling Group", Visible: true, SortBy: sortDesiredCapacity, SortDirection: tablewriter.Desc},
		{Name: "Min", Category: "Auto Scaling Group", Visible: true, SortBy: sortMinCapacity, SortDirection: tablewriter.Desc},
		{Name: "Max", Category: "Auto Scaling Group", Visible: true, SortBy: sortMaxCapacity, SortDirection: tablewriter.Desc},
		{Name: "ARN", Category: "Auto Scaling Group", Visible: showARNs},
	}
}

func getInstanceFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Instance", Visible: true, SortBy: sortName, SortDirection: tablewriter.Asc},
		{Name: "State", Category: "Instance", Visible: true},
		{Name: "Instance Type", Category: "Instance", Visible: true},
		{Name: "Launch Template/Configuration", Category: "Instance", Visible: true},
		{Name: "Availability Zone", Category: "Instance", Visible: true},
		{Name: "Health", Category: "Instance", Visible: true},
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
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListAutoScalingGroups(cobraCmd, args))
	},
}

// newLsFlags is the function for adding flags to the ls command
func newLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Auto-Scaling Groups in list format.")
	cobraCmd.Flags().BoolVarP(&showARNs, "arn", "a", false, "Show ARNs for each Auto-Scaling Group.")
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending ASG name.")
	cobraCmd.Flags().BoolVarP(&sortInstances, "sort-instances", "i", false, "Sort by descending number of instances. (ASG output only)")
	cobraCmd.Flags().BoolVarP(&sortDesiredCapacity, "sort-desired-capacity", "d", false, "Sort by descending desired capacity. (ASG output only)")
	cobraCmd.Flags().BoolVarP(&sortMinCapacity, "sort-min-capacity", "m", false, "Sort by descending min capacity. (ASG output only)")
	cobraCmd.Flags().BoolVarP(&sortMaxCapacity, "sort-max-capacity", "M", false, "Sort by descending max capacity. (ASG output only)")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
}

//
// Command functions
//

func ListAutoScalingGroups(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, asg.NewAutoScalingService)
	if err != nil {
		return fmt.Errorf("create new Auto Scaling Group service: %w", err)
	}

	if len(args) > 0 {
		fmt.Printf("Listing instances for Auto Scaling Group %s\n", args[0])
		return ListAutoScalingGroupInstances(svc, args[0])
	} else {
		autoScalingGroups, err := svc.GetAutoScalingGroups(cmd.Context(), &ascTypes.GetAutoScalingGroupsInput{})
		if err != nil {
			return fmt.Errorf("get Auto Scaling Groups: %w", err)
		}

		table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
			Title: "Auto Scaling Groups",
		})
		if list {
			table.SetRenderStyle("plain")
		}

		fields := getListFields()
		fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(autoScalingGroups))

		headerRow := tablewriter.BuildHeaderRow(fields)
		table.AppendHeader(headerRow)
		table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(autoScalingGroups), fields, asg.GetFieldValue, asg.GetTagValue))
		table.SetFieldConfigs(fields, reverseSort)
		table.Render()
		return nil
	}
}

// ListAutoScalingGroupInstances is the function for listing instances in an Auto Scaling Group
func ListAutoScalingGroupInstances(svc *asg.AutoScalingService, asgName string) error {
	instances, err := svc.GetAutoScalingGroupInstances(
		context.TODO(),
		&ascTypes.GetAutoScalingGroupInstancesInput{
			AutoScalingGroupNames: []string{asgName},
		},
	)
	if err != nil {
		return fmt.Errorf("get instances for Auto Scaling Group %s: %w", asgName, err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Auto Scaling Group Instances",
	})
	if list {
		table.SetRenderStyle("plain")
	}

	fields := getInstanceFields()
	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(instances), fields, asg.GetFieldValue, asg.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)
	table.Render()
	return nil
}
