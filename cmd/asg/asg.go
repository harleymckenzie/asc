package asg

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/asg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	sortOrder       []string
	list            bool
	showARNs        bool
	selectedColumns []string
)

func NewASGCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asg",
		Short: "Perform Auto Scaling Group operations",
	}

	lsCmd := &cobra.Command{
		Use:   "ls [asg-name]",
		Short: "List all Auto Scaling Groups or instances within a specific ASG",
		Long: `List all Auto Scaling Groups or instances within a specific ASG.

When listing ASGs (no argument):
  - Sort by name (default), instances, desired capacity, min capacity, or max capacity
  - Use -l for list format
  - Use --arn to show ARNs

When listing instances (with ASG name):
  - Sort by name (default), instance type, or launch template/configuration
  - Use -l for list format`,
		PreRun: func(cobraCmd *cobra.Command, args []string) {
			// Clear any existing sort order
			sortOrder = []string{}

			// Set default columns based on whether we're listing ASGs or instances
			if len(args) > 0 {
				// Default columns for instances
				selectedColumns = []string{
					"name",
					"state",
					"instance_type",
					"launch_config",
					"availability_zone",
					"health",
				}
			} else {
				// Default columns for ASGs
				selectedColumns = []string{
					"name",
					"instances",
					"desired_capacity",
					"min_capacity",
					"max_capacity",
				}
				if showARNs {
					selectedColumns = append(selectedColumns, "arn")
				}
			}

			// Visit flags in the order they appear in the command line
			cobraCmd.Flags().Visit(func(f *pflag.Flag) {
				switch f.Name {
				case "sort-name":
					sortOrder = append(sortOrder, "Name")
				case "sort-instances":
					if len(args) > 0 {
						sortOrder = append(sortOrder, "Instance Type")
					} else {
						sortOrder = append(sortOrder, "Instances")
					}
				case "sort-desired-capacity":
					sortOrder = append(sortOrder, "Desired")
				case "sort-min-capacity":
					sortOrder = append(sortOrder, "Min")
				case "sort-max-capacity":
					sortOrder = append(sortOrder, "Max")
				case "sort-launch-config":
					sortOrder = append(sortOrder, "Launch Template/Configuration")
				}
			})
		},
		Run: func(cobraCmd *cobra.Command, args []string) {
			ctx := context.TODO()
			profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
			region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

			svc, err := asg.NewAutoScalingService(ctx, profile, region)
			if err != nil {
				log.Fatalf("Failed to initialize ASG service: %v", err)
			}

			if len(args) > 0 {
				err = svc.ListAutoScalingGroupInstances(ctx, args[0], sortOrder, list, selectedColumns)
			} else {
				err = svc.ListAutoScalingGroups(ctx, sortOrder, list, selectedColumns)
			}
			if err != nil {
				log.Fatalf("Failed to list ASG resources: %v", err)
			}
		},
	}
	cmd.AddCommand(lsCmd)

	// Add flags - Output
	lsCmd.Flags().BoolVarP(&list, "list", "l", false, "Output in list format")
	lsCmd.Flags().BoolVar(&showARNs, "arn", false, "Show ARNs for each Auto-Scaling Group (ASG list only)")

	// Add flags - Sorting
	lsCmd.Flags().BoolP("sort-name", "n", true, "Sort by descending name")
	lsCmd.Flags().BoolP("sort-instances", "i", false, "Sort by descending number of instances (ASG list) or instance type (instance list)")
	lsCmd.Flags().BoolP("sort-desired-capacity", "d", false, "Sort by descending desired capacity (ASG list only)")
	lsCmd.Flags().BoolP("sort-min-capacity", "m", false, "Sort by descending min capacity (ASG list only)")
	lsCmd.Flags().BoolP("sort-max-capacity", "M", false, "Sort by descending max capacity (ASG list only)")
	lsCmd.Flags().BoolP("sort-launch-config", "c", false, "Sort by descending launch template/configuration (instance list only)")
	lsCmd.Flags().SortFlags = false

	return cmd
}
