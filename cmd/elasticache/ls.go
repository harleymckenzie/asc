package elasticache

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/elasticache"
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

		columns := []Column{
			{ID: "name", Visible: true},
			{ID: "status", Visible: true},
			{ID: "engine_version", Visible: true},
			{ID: "instance_type", Visible: true},
			{ID: "endpoint", Visible: showEndpoint},
		}

		selectedColumns := make([]string, 0, len(columns))

		for _, col := range columns {
			if col.Visible {
				selectedColumns = append(selectedColumns, col.ID)
			}
		}

		tableformat.Render(&elasticache.ElasticacheTable{
			Instances:       instances,
			SelectedColumns: selectedColumns,
			SortOrder:       sortOrder,
		})
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
	lsCmd.Flags().SortFlags = false
}
