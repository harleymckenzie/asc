// The ls command lists all RDS clusters and instances.

package rds

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/rds"
	ascTypes "github.com/harleymckenzie/asc/internal/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

// Variables
var (
	list                 bool
	showEndpoint         bool
	showEngineVersion    bool
	showModificationInfo bool

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
	newLsFlags(lsCmd)
}

// Column functions
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Cluster Identifier", Category: "RDS", Visible: true, DefaultSort: true, Merge: true, SortBy: sortCluster, SortDirection: tablewriter.Asc},
		{Name: "Identifier", Category: "RDS", Visible: true, SortBy: sortName, SortDirection: tablewriter.Asc},
		{Name: "Status", Category: "RDS", Visible: true, SortBy: sortStatus, SortDirection: tablewriter.Asc},
		{Name: "Role", Category: "RDS", Visible: true, SortBy: sortRole, SortDirection: tablewriter.Asc},
		{Name: "Engine", Category: "RDS", Visible: true, SortBy: sortEngine, SortDirection: tablewriter.Asc},
		{Name: "Engine Version", Category: "RDS", Visible: showEngineVersion},
		{Name: "Class", Category: "RDS", Visible: true},
		{Name: "Endpoint", Category: "RDS", Visible: showEndpoint},
		{Name: "Pending Modifications", Category: "RDS", Visible: showModificationInfo},
		{Name: "Maintenance Window", Category: "RDS", Visible: showModificationInfo},
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
func newLsFlags(cobraCmd *cobra.Command) {
	// Add flags - Output
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs RDS clusters and instances in list format.")
	cobraCmd.Flags().BoolVarP(&showEndpoint, "endpoint", "e", false, "Show the endpoint of the cluster")
	cobraCmd.Flags().BoolVarP(&showEngineVersion, "engine-version", "v", false, "Show the engine version of the cluster")
	cobraCmd.Flags().BoolVarP(&showModificationInfo, "modification-info", "m", false, "Show the modification info of the instance")
	cmdutil.AddTagFlag(cobraCmd)

	// Add flags - Sorting
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending RDS instance identifier.")
	cobraCmd.Flags().BoolVarP(&sortCluster, "sort-cluster", "c", false, "Sort by descending RDS cluster identifier.")
	cobraCmd.Flags().BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending RDS instance class.")
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
	svc, err := cmdutil.CreateService(cmd, rds.NewRDSService)
	if err != nil {
		return fmt.Errorf("create new RDS service: %w", err)
	}

	instances, err := svc.GetInstances(cmd.Context(), &ascTypes.GetInstancesInput{})
	if err != nil {
		return fmt.Errorf("list RDS instances: %w", err)
	}

	clusters, err := svc.GetClusters(cmd.Context(), &ascTypes.GetClustersInput{})
	if err != nil {
		return fmt.Errorf("list RDS clusters: %w", err)
	}

	// Set clusters context for role calculation
	rds.SetClustersContext(clusters)

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         "Databases",
		Style:         "rounded-separated",
		PlainStyle:    list,
		Fields:        getListFields(),
		Tags:          cmdutil.Tags,
		Data:          utils.SlicesToAny(instances),
		GetFieldValue: rds.GetFieldValue,
		GetTagValue:   rds.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}
