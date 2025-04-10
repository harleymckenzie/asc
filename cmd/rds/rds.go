package rds

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/base"
	"github.com/harleymckenzie/asc/pkg/service/rds"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	sortOrder       []string
	list            bool
	selectedColumns []string
	showEndpoint    bool
	force           bool
	waitSync        bool
)

func NewRDSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rds",
		Short: "Perform RDS operations",
	}

	// List command
	lsCmd := &cobra.Command{
		Use:   "ls",
		Short: "List all RDS clusters and instances",
		PreRun: func(cobraCmd *cobra.Command, args []string) {
			// Clear any existing sort order
			sortOrder = []string{}

			// Set default columns
			selectedColumns = []string{
				"cluster_identifier",
				"identifier",
				"status",
				"engine",
				"size",
				"role",
			}

			if showEndpoint {
				selectedColumns = append(selectedColumns, "endpoint")
			}

			// Visit flags in the order they appear in the command line
			cobraCmd.Flags().Visit(func(f *pflag.Flag) {
				switch f.Name {
				case "sort-name":
					sortOrder = append(sortOrder, "Identifier")
				case "sort-cluster":
					sortOrder = append(sortOrder, "Cluster Identifier")
				case "sort-type":
					sortOrder = append(sortOrder, "Size")
				case "sort-engine":
					sortOrder = append(sortOrder, "Engine")
				case "sort-status":
					sortOrder = append(sortOrder, "Status")
				case "sort-role":
					sortOrder = append(sortOrder, "Role")
				}
			})
		},
		Run: func(cobraCmd *cobra.Command, args []string) {
			ctx := context.TODO()
			profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
			region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

			svc, err := rds.NewRDSService(ctx, profile, region)
			if err != nil {
				log.Fatalf("Failed to initialize RDS service: %v", err)
			}

			options := base.ListOptions{
				CommandOptions: base.CommandOptions{
					Profile: profile,
					Region:  region,
				},
				SortOrder:       sortOrder,
				List:            list,
				SelectedColumns: selectedColumns,
			}

			if err := svc.ListInstances(ctx, options); err != nil {
				log.Fatalf("Failed to list RDS instances: %v", err)
			}
		},
	}

	// Stop command
	stopCmd := &cobra.Command{
		Use:   "stop [instance-id...]",
		Short: "Stop RDS instances",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cobraCmd *cobra.Command, args []string) {
			ctx := context.TODO()
			profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
			region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

			svc, err := rds.NewRDSService(ctx, profile, region)
			if err != nil {
				log.Fatalf("Failed to initialize RDS service: %v", err)
			}

			resources := make([]base.ResourceIdentifier, len(args))
			for i, id := range args {
				resources[i] = base.ResourceIdentifier{
					Name: id,
					Type: "db-instance",
				}
			}

			options := base.StateChangeOptions{
				CommandOptions: base.CommandOptions{
					Profile:  profile,
					Region:   region,
					Force:    force,
					WaitSync: waitSync,
				},
				ResourceIDs: resources,
			}

			if err := svc.StopInstances(ctx, options); err != nil {
				log.Fatalf("Failed to stop RDS instances: %v", err)
			}
		},
	}

	// Start command
	startCmd := &cobra.Command{
		Use:   "start [instance-id...]",
		Short: "Start RDS instances",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cobraCmd *cobra.Command, args []string) {
			ctx := context.TODO()
			profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
			region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

			svc, err := rds.NewRDSService(ctx, profile, region)
			if err != nil {
				log.Fatalf("Failed to initialize RDS service: %v", err)
			}

			resources := make([]base.ResourceIdentifier, len(args))
			for i, id := range args {
				resources[i] = base.ResourceIdentifier{
					Name: id,
					Type: "db-instance",
				}
			}

			options := base.StateChangeOptions{
				CommandOptions: base.CommandOptions{
					Profile:  profile,
					Region:   region,
					Force:    force,
					WaitSync: waitSync,
				},
				ResourceIDs: resources,
			}

			if err := svc.StartInstances(ctx, options); err != nil {
				log.Fatalf("Failed to start RDS instances: %v", err)
			}
		},
	}

	cmd.AddCommand(lsCmd, stopCmd, startCmd)

	// Add flags - Output (ls command)
	lsCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs RDS clusters and instances in list format.")
	lsCmd.Flags().BoolVarP(&showEndpoint, "endpoint", "e", false, "Show the endpoint of the cluster")

	// Add flags - Sorting (ls command)
	lsCmd.Flags().BoolP("sort-name", "n", true, "Sort by descending RDS instance identifier.")
	lsCmd.Flags().BoolP("sort-cluster", "c", false, "Sort by descending RDS cluster identifier.")
	lsCmd.Flags().BoolP("sort-type", "T", false, "Sort by descending RDS instance type.")
	lsCmd.Flags().BoolP("sort-engine", "E", false, "Sort by descending database engine type.")
	lsCmd.Flags().BoolP("sort-status", "s", false, "Sort by descending RDS instance status.")
	lsCmd.Flags().BoolP("sort-role", "R", false, "Sort by descending RDS instance role.")
	lsCmd.Flags().SortFlags = false

	// Add flags - State change commands (start/stop)
	for _, subCmd := range []*cobra.Command{startCmd, stopCmd} {
		subCmd.Flags().BoolVarP(&force, "force", "f", false, "Force the operation without confirmation")
		subCmd.Flags().BoolVarP(&waitSync, "wait", "w", false, "Wait for the operation to complete")
	}

	return cmd
}
