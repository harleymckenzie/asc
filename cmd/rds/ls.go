// The ls command lists all RDS clusters and instances.

package rds

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/rds"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

// Variables
var (
	list              bool
	showEndpoint      bool
	showEngineVersion bool

	sortName    bool
	sortCluster bool
	sortType    bool
	sortEngine  bool
	sortStatus  bool
	sortRole    bool
)

// Init function
func init() {
	addLsFlags(lsCmd)
}

// Column functions
func rdsColumns() []tableformat.Column {
	return []tableformat.Column{
		{ID: "Cluster Identifier", Visible: true, Sort: sortCluster},
		{ID: "Identifier", Visible: true, Sort: sortName},
		{ID: "Status", Visible: true, Sort: sortStatus},
		{ID: "Engine", Visible: true, Sort: sortEngine},
		{ID: "Engine Version", Visible: showEngineVersion, Sort: false},
		{ID: "Size", Visible: true, Sort: false},
		{ID: "Role", Visible: true, Sort: sortRole},
		{ID: "Endpoint", Visible: showEndpoint, Sort: false},
	}
}

// Command variable
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all RDS clusters and instances",
	GroupID: "actions",
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

		clusters, err := svc.GetClusters(ctx)
		if err != nil {
			log.Fatalf("Failed to list RDS clusters: %v", err)
		}

		columns := rdsColumns()
		selectedColumns, sortBy := tableformat.BuildColumns(columns)

		opts := tableformat.RenderOptions{
			SortBy: sortBy,
			List:   list,
			Title:  "RDS Clusters and Instances",
		}

		tableformat.Render(&rds.RDSTable{
			Instances:       instances,
			Clusters:        clusters,
			SelectedColumns: selectedColumns,
		}, opts)
	},
}

// Flag function
func addLsFlags(cobraCmd *cobra.Command) {
	// Add flags - Output
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs RDS clusters and instances in list format.")
	cobraCmd.Flags().BoolVarP(&showEndpoint, "endpoint", "e", false, "Show the endpoint of the cluster")
	cobraCmd.Flags().BoolVarP(&showEngineVersion, "engine-version", "v", false, "Show the engine version of the cluster")

	// Add flags - Sorting
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending RDS instance identifier.")
	cobraCmd.Flags().BoolVarP(&sortCluster, "sort-cluster", "c", false, "Sort by descending RDS cluster identifier.")
	cobraCmd.Flags().BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending RDS instance type.")
	cobraCmd.Flags().BoolVarP(&sortEngine, "sort-engine", "E", false, "Sort by descending database engine type.")
	cobraCmd.Flags().BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by descending RDS instance status.")
	cobraCmd.Flags().BoolVarP(&sortRole, "sort-role", "R", false, "Sort by descending RDS instance role.")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-cluster", "sort-type", "sort-engine", "sort-status", "sort-role")

	cobraCmd.Flags().SortFlags = false
}
