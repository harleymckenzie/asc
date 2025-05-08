// The ls command lists Elasticache clusters.

package elasticache

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/elasticache"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

// Variables
var (
	list         bool
	showEndpoint bool

	sortName   bool
	sortType   bool
	sortStatus bool
	sortEngine bool
)

// Init function
func init() {
	newLsFlags(lsCmd)
}

// Column functions
func elasticacheColumns() []tableformat.Column {
	return []tableformat.Column{
		{ID: "Cache Name", Visible: true, Sort: sortName},
		{ID: "Status", Visible: true},
		{ID: "Engine Version", Visible: true},
		{ID: "Configuration", Visible: true},
		{ID: "Endpoint", Visible: showEndpoint},
	}
}

// Command variable
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List Elasticache clusters",
	GroupID: "actions",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()
		profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
		region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

		svc, err := elasticache.NewElasticacheService(ctx, profile, region)
		if err != nil {
			log.Fatalf("Failed to initialize Elasticache service: %v", err)
		}

		instances, err := svc.GetInstances(ctx)
		if err != nil {
			log.Fatalf("Failed to list Elasticache instances: %v", err)
		}

		columns := elasticacheColumns()
		selectedColumns, sortBy := tableformat.BuildColumns(columns)

		opts := tableformat.RenderOptions{
			SortBy: sortBy,
			List:   list,
			Title:  "Elasticache Clusters",
		}

		tableformat.Render(&elasticache.ElasticacheTable{
			Instances:       instances,
			SelectedColumns: selectedColumns,
		}, opts)
	},
}

// Flag function
func newLsFlags(cobraCmd *cobra.Command) {
	// Add flags - Output
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Elasticache clusters in list format.")
	cobraCmd.Flags().BoolVarP(&showEndpoint, "endpoint", "e", false, "Show the endpoint of the cluster")

	// Add flags - Sorting
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", true, "Sort by descending Elasticache cluster name.")
	cobraCmd.Flags().BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending Elasticache cluster type.")
	cobraCmd.Flags().BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by descending Elasticache cluster status.")
	cobraCmd.Flags().BoolVarP(&sortEngine, "sort-engine", "E", false, "Sort by descending Elasticache cluster engine version.")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-type", "sort-status", "sort-engine")
}
