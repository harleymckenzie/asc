package rds

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/rds"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

type Column struct {
	ID      string
	Visible bool
	Sort    bool
}

var (
	list   bool
	sortBy string

	showEndpoint      bool
	showEngineVersion bool

	sortName    bool
	sortCluster bool
	sortType    bool
	sortEngine  bool
	sortStatus  bool
	sortRole    bool
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all RDS clusters and instances",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()
		profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
		region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

		svc, err := rds.NewRDSService(ctx, profile, region)
		if err != nil {
			log.Fatalf("Failed to initialize RDS service: %v", err)
		}

		instances, err := svc.GetInstances(ctx)
		if err != nil {
			log.Fatalf("Failed to list EC2 instances: %v", err)
		}

		// Define available columns and associated flags
		columns := []Column{
			{ID: "Cluster Identifier", Visible: true, Sort: sortCluster},
			{ID: "Identifier", Visible: true, Sort: sortName},
			{ID: "Status", Visible: true, Sort: sortStatus},
			{ID: "Engine", Visible: true, Sort: sortEngine},
			{ID: "Engine Version", Visible: showEngineVersion, Sort: false},
			{ID: "Size", Visible: true},
			{ID: "Role", Visible: true, Sort: sortRole},
			{ID: "Endpoint", Visible: showEndpoint},
		}

		selectedColumns := make([]string, 0, len(columns))

		// Dynamically build the list of columns
		for _, col := range columns {
			if col.Visible {
				selectedColumns = append(selectedColumns, col.ID)
			}
		}

		tableformat.Render(&rds.RDSTable{
			Instances:       instances,
			SelectedColumns: selectedColumns,
		}, sortBy)
	},
}

func init() {
	// Add flags - Output
	lsCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs RDS clusters and instances in list format.")
	lsCmd.Flags().BoolVarP(&showEndpoint, "endpoint", "e", false, "Show the endpoint of the cluster")
	lsCmd.Flags().BoolVarP(&showEngineVersion, "engine-version", "v", false, "Show the engine version of the cluster")

	// Add flags - Sorting
	lsCmd.Flags().BoolVarP(&sortName, "sort-name", "n", true, "Sort by descending RDS instance identifier.")
	lsCmd.Flags().BoolVarP(&sortCluster, "sort-cluster", "c", false, "Sort by descending RDS cluster identifier.")
	lsCmd.Flags().BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending RDS instance type.")
	lsCmd.Flags().BoolVarP(&sortEngine, "sort-engine", "E", false, "Sort by descending database engine type.")
	lsCmd.Flags().BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by descending RDS instance status.")
	lsCmd.Flags().BoolVarP(&sortRole, "sort-role", "R", false, "Sort by descending RDS instance role.")
	lsCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-cluster", "sort-type", "sort-engine", "sort-status", "sort-role")

	lsCmd.Flags().SortFlags = false
}
