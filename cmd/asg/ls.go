package asg

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/asg"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
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
	Short: "List all Auto Scaling Groups",
	Long: "List Auto Scaling Groups, instances in an ASG, or schedules\n" +
		"  ls                      List all Auto Scaling Groups\n" +
		"  ls [asg-name]           List instances in the specified ASG\n" +
		"  ls schedules [asg-name] List schedules (optionally for specific ASG)",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()
		profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
		region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

		svc, err := asg.NewAutoScalingService(ctx, profile, region)
		if err != nil {
			log.Fatalf("Failed to initialize Auto Scaling Group service: %v", err)
		}

		// If an argument is provided that isn't a subcommand, show instance details using GetInstances
		if len(args) > 0 {
			ListAutoScalingGroupInstances(svc, args[0])
		} else {
			ListAutoScalingGroups(svc)
		}
	},
}

func ListAutoScalingGroupInstances(svc *asg.AutoScalingService, asgName string) {
	ctx := context.TODO()
	instances, err := svc.GetInstances(ctx, &asg.GetInstancesInput{
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

func ListAutoScalingGroups(svc *asg.AutoScalingService) {
	ctx := context.TODO()
	autoScalingGroups, err := svc.GetAutoScalingGroups(ctx)
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

func init() {
	// Add flags - Output
	lsCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Auto-Scaling Groups in list format.")
	lsCmd.Flags().BoolVar(&showARNs, "arn", false, "Show ARNs for each Auto-Scaling Group.")

	// Add flags - Sorting
	lsCmd.Flags().BoolP("sort-name", "n", true, "Sort by descending ASG name.")
	lsCmd.Flags().BoolP("sort-instances", "i", false, "Sort by descending number of instances.")
	lsCmd.Flags().BoolP("sort-desired-capacity", "d", false, "Sort by descending desired capacity.")
	lsCmd.Flags().BoolP("sort-min-capacity", "m", false, "Sort by descending min capacity.")
	lsCmd.Flags().BoolP("sort-max-capacity", "M", false, "Sort by descending max capacity.")
	lsCmd.Flags().SortFlags = false

	lsCmd.AddCommand(lsSchedulesCmd)
}
