package rds

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/rds"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	sortOrder []string
	list      bool
)

func NewRDSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rds",
		Short: "Perform RDS operations",
	}

	// ls sub command
	lsCmd := &cobra.Command{
		Use:    "ls",
		Short:  "List all RDS clusters and instances",
		PreRun: preRunLs,
		Run:    runLs,
	}

	cmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolVarP(&list, "list", "l", false,
		"Outputs RDS clusters and instances in list format.")
	lsCmd.Flags().BoolP("sort-name", "n", true,
		"Sort by descending RDS instance identifier.")
	lsCmd.Flags().BoolP("sort-cluster", "c", false,
		"Sort by descending RDS cluster identifier.")
	lsCmd.Flags().BoolP("sort-type", "T", false,
		"Sort by descending RDS instance type.")
	lsCmd.Flags().BoolP("sort-engine", "E", false,
		"Sort by descending database engine type.")
	lsCmd.Flags().SortFlags = false

	return cmd
}

func preRunLs(cobraCmd *cobra.Command, args []string) {
	sortOrder = []string{}
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
		}
	})
}

func runLs(cobraCmd *cobra.Command, args []string) {
	ctx := context.TODO()

	svc, err := rds.NewRDSService(ctx, "akoovadev")
	if err != nil {
		log.Fatalf("Failed to initialize RDS service: %v", err)
	}

	err = svc.ListInstances(ctx, sortOrder, list)
	if err != nil {
		log.Fatalf("Error describing database clusters and instances: %v", err)
	}
}
