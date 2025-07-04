// The ls command lists all RDS clusters and instances.

package rds

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/rds"
	ascTypes "github.com/harleymckenzie/asc/internal/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

// Variables
var (
	list              bool
	showEndpoint      bool
	showEngineVersion bool
	showMaintenanceWindow bool

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
		{ID: "Cluster Identifier", Display: true, Sort: sortCluster, Merge: true},
		{ID: "Identifier", Display: true, Sort: sortName},
		{ID: "Status", Display: true, Sort: sortStatus},
		{ID: "Role", Display: true, Sort: sortRole},
		{ID: "Engine", Display: true, Sort: sortEngine},
		{ID: "Engine Version", Display: showEngineVersion, Sort: false, SortDirection: "desc"},
		{ID: "Size", Display: true, Sort: false},
		{ID: "Endpoint", Display: showEndpoint, Sort: false},
		{ID: "Maintenance Window", Display: showMaintenanceWindow, Sort: false},
	}
}

// Command variable
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all RDS clusters and instances",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListRDSClusters(cmd, args))
	},
}

// Flag function
func addLsFlags(cobraCmd *cobra.Command) {
	// Add flags - Output
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs RDS clusters and instances in list format.")
	cobraCmd.Flags().BoolVarP(&showEndpoint, "endpoint", "e", false, "Show the endpoint of the cluster")
	cobraCmd.Flags().BoolVarP(&showEngineVersion, "engine-version", "v", false, "Show the engine version of the cluster")
	cobraCmd.Flags().BoolVarP(&showMaintenanceWindow, "maintenance-window", "P", false, "Show the maintenance window of the cluster")

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
func ListRDSClusters(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)

	svc, err := rds.NewRDSService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new RDS service: %w", err)
	}

	instances, err := svc.GetInstances(ctx, &ascTypes.GetInstancesInput{})
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
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return rds.GetAttributeValue(fieldID, instance, clusters)
		},
	}, opts)

	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}

	return nil
}
