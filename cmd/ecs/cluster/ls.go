package cluster

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ecs"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	reverseSort bool

	sortName   bool
	sortStatus bool
)

func init() {
	newLsFlags(lsCmd)
}

func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Cluster", Visible: true, DefaultSort: true, SortBy: sortName, SortDirection: tablewriter.Asc},
		{Name: "Status", Category: "Cluster", Visible: true, SortBy: sortStatus, SortDirection: tablewriter.Asc},
		{Name: "Active Services", Category: "Cluster", Visible: true},
		{Name: "Running Tasks", Category: "Cluster", Visible: true},
		{Name: "Pending Tasks", Category: "Cluster", Visible: true},
		{Name: "Registered Instances", Category: "Cluster", Visible: true},
	}
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List ECS clusters",
	Aliases: []string{"list"},
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListClusters(cmd, args))
	},
}

func newLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs in list format.")
	cmdutil.AddTagFlag(cobraCmd)

	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by cluster name.")
	cobraCmd.Flags().BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by cluster status.")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-status")

	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cobraCmd.Flags().SortFlags = false
}

func ListClusters(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, ecs.NewECSService)
	if err != nil {
		return fmt.Errorf("create ECS service: %w", err)
	}

	clusters, err := svc.GetAllClusters(cmd.Context())
	if err != nil {
		return fmt.Errorf("list ECS clusters: %w", err)
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         "ECS Clusters",
		PlainStyle:    list,
		Fields:        getListFields(),
		Tags:          cmdutil.Tags,
		Data:          utils.SlicesToAny(clusters),
		GetFieldValue: ecs.GetFieldValue,
		GetTagValue:   ecs.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}
