package asg

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/asg"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

type Column struct {
	ID      string
	Visible bool
}

var (
	list      bool
	sortOrder []string

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
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()
		profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
		region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

		svc, err := asg.NewAutoScalingService(ctx, profile, region)
		if err != nil {
			log.Fatalf("Failed to initialize Auto Scaling Group service: %v", err)
		}

		// If a specific ASG is provided, show instance details using GetInstances,
		// Otherwise use GetAutoScalingGroups
		if len(args) > 0 {
			asgName := args[0]
			instances, err := svc.GetInstances(ctx, &asg.GetInstancesInput{
				AutoScalingGroupNames: []string{asgName},
			})
			if err != nil {
				log.Fatalf("Failed to get instances for Auto Scaling Group %s: %v", args[0], err)
			}

			// Define columns for instances
			columns := []Column{
				{ID: "instance_name", Visible: true},
				{ID: "state", Visible: true},
				{ID: "instance_type", Visible: true},
				{ID: "launch_config", Visible: true},
				{ID: "availability_zone", Visible: true},
				{ID: "health", Visible: true},
			}

			selectedColumns := make([]string, 0, len(columns))
			for _, col := range columns {
				if col.Visible {
					selectedColumns = append(selectedColumns, col.ID)
				}
			}

			tableformat.Render(&asg.AutoScalingInstanceTable{
				Instances:         instances,
				SelectedColumns:   selectedColumns,
				SortOrder:         sortOrder,
			})
		} else {
			autoScalingGroups, err := svc.GetAutoScalingGroups(ctx)
			if err != nil {
				log.Fatalf("Failed to get Auto Scaling Groups: %v", err)
			}

			// Define columns for Auto Scaling Groups
			columns := []Column{
				{ID: "name", Visible: true},
				{ID: "instances", Visible: true},
				{ID: "desired_capacity", Visible: true},
				{ID: "min_capacity", Visible: true},
				{ID: "max_capacity", Visible: true},
			}

            selectedColumns := make([]string, 0, len(columns))
            for _, col := range columns {
                if col.Visible {
                    selectedColumns = append(selectedColumns, col.ID)
                }
			}

			tableformat.Render(&asg.AutoScalingTable{
				AutoScalingGroups: autoScalingGroups,
				SelectedColumns:   selectedColumns,
				SortOrder:         sortOrder,
			})
        }
	},
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
}
