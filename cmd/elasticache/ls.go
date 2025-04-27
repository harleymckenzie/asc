package elasticache

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/elasticache"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

var (
	list      bool
	showEndpoint bool

	sortName   bool
	sortType   bool
	sortStatus bool
	sortEngine bool
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List Elasticache clusters",
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

		columns := []tableformat.Column{
			{ID: "Cache Name", Visible: true},
			{ID: "Status", Visible: true},
			{ID: "Engine Version", Visible: true},
			{ID: "Configuration", Visible: true},
			{ID: "Endpoint", Visible: showEndpoint},
		}
		selectedColumns, sortBy := tableformat.BuildColumns(columns)

		tableformat.Render(&elasticache.ElasticacheTable{
			Instances:       instances,
			SelectedColumns: selectedColumns,
		}, sortBy, list)
	},
}

func init() {
	// Add flags - Output
	lsCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Elasticache clusters in list format.")
	lsCmd.Flags().BoolVarP(&showEndpoint, "endpoint", "e", false, "Show the endpoint of the cluster")

	// Add flags - Sorting
	lsCmd.Flags().BoolVarP(&sortName, "sort-name", "n", true, "Sort by descending Elasticache cluster name.")
	lsCmd.Flags().BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending Elasticache cluster type.")
	lsCmd.Flags().BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by descending Elasticache cluster status.")
	lsCmd.Flags().BoolVarP(&sortEngine, "sort-engine", "E", false, "Sort by descending Elasticache cluster engine version.")
	lsCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-type", "sort-status", "sort-engine")

	lsCmd.Flags().SortFlags = false
}
