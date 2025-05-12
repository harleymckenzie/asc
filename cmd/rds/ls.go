// The ls command lists all RDS clusters and instances.

package rds

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/pkg/service/rds"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/harleymckenzie/asc/pkg/shared/utils"
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
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

	reverseSort bool
)

// Init function
func init() {
	addLsFlags(lsCmd)
}

// Column functions
func rdsFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Cluster Identifier", Visible: true, Sort: sortCluster, Merge: true},
		{ID: "Identifier", Visible: true, Sort: sortName},
		{ID: "Status", Visible: true, Sort: sortStatus},
		{ID: "Role", Visible: true, Sort: sortRole},
		{ID: "Engine", Visible: true, Sort: sortEngine},
		{ID: "Engine Version", Visible: showEngineVersion, Sort: false},
		{ID: "Size", Visible: true, Sort: false},
		{ID: "Endpoint", Visible: showEndpoint, Sort: false},
	}
}

// Command variable
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all RDS clusters and instances",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListRDSClusters(cobraCmd, args))
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

	// Add flags - Reverse Sort
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")

	cobraCmd.Flags().SortFlags = false
}

// ListRDSClusters is the function for listing RDS clusters and instances
func ListRDSClusters(cobraCmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	svc, err := rds.NewRDSService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new RDS service: %w", err)
	}

	instances, err := svc.GetInstances(ctx)
	if err != nil {
		return fmt.Errorf("list RDS instances: %w", err)
	}

	clusters, err := svc.GetClusters(ctx)
	if err != nil {
		return fmt.Errorf("list RDS clusters: %w", err)
	}

	fields := rdsFields()

	opts := tableformat.RenderOptions{
		Title:  "Databases",
		Style:  "rounded-separated",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(instances),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) string {
			return rds.GetAttributeValue(fieldID, instance, clusters)
		},
	}, opts)

	return nil
}
