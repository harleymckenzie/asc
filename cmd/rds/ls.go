// The ls command lists all RDS clusters and instances.

package rds

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/rds"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
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
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Cluster Identifier", Category: "RDS", Visible: true, SortBy: sortCluster, SortDirection: tablewriter.Asc},
		{Name: "Identifier", Category: "RDS", Visible: true, SortBy: sortName, SortDirection: tablewriter.Asc},
		{Name: "Status", Category: "RDS", Visible: true, SortBy: sortStatus, SortDirection: tablewriter.Asc},
		{Name: "Role", Category: "RDS", Visible: true, SortBy: sortRole, SortDirection: tablewriter.Asc},
		{Name: "Engine", Category: "RDS", Visible: true, SortBy: sortEngine, SortDirection: tablewriter.Asc},
		{Name: "Engine Version", Category: "RDS", Visible: showEngineVersion},
		{Name: "Size", Category: "RDS", Visible: true},
		{Name: "Endpoint", Category: "RDS", Visible: showEndpoint},
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
	cmdutil.AddTagFlag(cobraCmd)

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
	svc, err := cmdutil.CreateService(cmd, rds.NewRDSService)
	if err != nil {
		return fmt.Errorf("create new RDS service: %w", err)
	}

	instances, err := svc.GetInstances(cmd.Context())
	if err != nil {
		return fmt.Errorf("list RDS instances: %w", err)
	}

	clusters, err := svc.GetClusters(cmd.Context())
	if err != nil {
		return fmt.Errorf("list RDS clusters: %w", err)
	}

	// Set clusters context for role calculation
	rds.SetClustersContext(clusters)

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Databases",
	})
	if list {
		table.SetRenderStyle("plain")
	}

	fields := getListFields()
	fields = cmdutil.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(instances))

	headerRow := cmdutil.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(cmdutil.BuildRows(utils.SlicesToAny(instances), fields, rds.GetFieldValue, rds.GetTagValue))
	table.SortBy(fields, reverseSort)
	table.Render()
	return nil
}
