// The ls command lists Elasticache clusters.

package elasticache

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/elasticache"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/utils"
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

	reverseSort bool
)

// Init function
func init() {
	newLsFlags(lsCmd)
}

// Column functions
func elasticacheFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Cache Name", Visible: true, Sort: sortName, DefaultSort: true},
		{ID: "Status", Visible: true, Sort: sortStatus},
		{ID: "Engine Version", Visible: true, Sort: sortEngine},
		{ID: "Configuration", Visible: true, Sort: sortType},
		{ID: "Endpoint", Visible: showEndpoint},
	}
}

// Command variable
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List Elasticache clusters",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListElasticacheClusters(cobraCmd, args))
	},
}

// ListElasticacheClusters is the function for listing Elasticache clusters
func ListElasticacheClusters(cobraCmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cobraCmd)

	svc, err := elasticache.NewElasticacheService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new Elasticache service: %w", err)
	}

	instances, err := svc.GetInstances(ctx)
	if err != nil {
		return fmt.Errorf("list Elasticache instances: %w", err)
	}

	fields := elasticacheFields()

	opts := tableformat.RenderOptions{
		Title:  "Elasticache Clusters",
		Style:  "rounded-separated",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(instances),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return elasticache.GetAttributeValue(fieldID, instance)
		},
	}, opts)
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}
	return nil
}

// Flag function
func newLsFlags(cobraCmd *cobra.Command) {
	// Add flags - Output
	cobraCmd.Flags().
		BoolVarP(&list, "list", "l", false, "Outputs Elasticache clusters in list format.")
	cobraCmd.Flags().
		BoolVarP(&showEndpoint, "endpoint", "e", false, "Show the endpoint of the cluster")

	// Add flags - Sorting
	cobraCmd.Flags().
		BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending Elasticache cluster name.")
	cobraCmd.Flags().
		BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending Elasticache cluster type.")
	cobraCmd.Flags().
		BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by descending Elasticache cluster status.")
	cobraCmd.Flags().
		BoolVarP(&sortEngine, "sort-engine", "E", false, "Sort by descending Elasticache cluster engine version.")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-type", "sort-status", "sort-engine")

	// Add flags - Reverse Sort
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
}
