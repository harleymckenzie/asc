// The ls command lists Elasticache clusters.

package elasticache

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/elasticache"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

// Variables
var (
	list         bool
	showEndpoint bool

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
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Cache Name", Category: "Cluster Details", Visible: true, DefaultSort: true},
		{Name: "Status", Category: "Cluster Details", Visible: true, SortBy: sortStatus, SortDirection: tablewriter.Asc},
		{Name: "Engine Version", Category: "Cluster Details", Visible: true, SortBy: sortEngine, SortDirection: tablewriter.Desc},
		{Name: "Configuration", Category: "Cluster Details", Visible: true, SortBy: sortType, SortDirection: tablewriter.Asc},
		{Name: "Endpoint", Category: "Network", Visible: showEndpoint},
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
func ListElasticacheClusters(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, elasticache.NewElasticacheService)
	if err != nil {
		return fmt.Errorf("create elasticache service: %w", err)
	}

	instances, err := svc.GetInstances(cmd.Context())
	if err != nil {
		return fmt.Errorf("get instances: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Elasticache Clusters",
	})
	if list {
		table.SetStyle("plain")
	}
	fields := getListFields()

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(instances), fields, elasticache.GetFieldValue, elasticache.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)

	table.Render()
	return nil
}

// Flag function
func newLsFlags(cobraCmd *cobra.Command) {
	// Add flags - Output
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Elasticache clusters in list format.")
	cobraCmd.Flags().BoolVarP(&showEndpoint, "endpoint", "e", false, "Show the endpoint of the cluster")

	// Add flags - Sorting
	cobraCmd.Flags().BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending Elasticache cluster type.")
	cobraCmd.Flags().BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by descending Elasticache cluster status.")
	cobraCmd.Flags().BoolVarP(&sortEngine, "sort-engine", "E", false, "Sort by descending Elasticache cluster engine version.")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-type", "sort-status", "sort-engine")

	// Add flags - Reverse Sort
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
}
