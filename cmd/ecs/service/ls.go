package service

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
	cluster     string

	showARN         bool
	showCreatedDate bool

	sortName   bool
	sortStatus bool
)

func init() {
	newLsFlags(lsCmd)
}

func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Service", Visible: true, DefaultSort: true, SortBy: sortName, SortDirection: tablewriter.Asc},
		{Name: "Status", Category: "Service", Visible: true, SortBy: sortStatus, SortDirection: tablewriter.Asc},
		{Name: "Launch Type", Category: "Service", Visible: true},
		{Name: "Task Definition", Category: "Service", Visible: true},
		{Name: "Desired Count", Category: "Service", Visible: true},
		{Name: "Running Count", Category: "Service", Visible: true},
		{Name: "ARN", Category: "Service", Visible: showARN},
		{Name: "Created Date", Category: "Service", Visible: showCreatedDate},
	}
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List ECS services",
	Aliases: []string{"list"},
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListServices(cmd, args))
	},
}

func newLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs in list format.")
	cobraCmd.Flags().StringVarP(&cluster, "cluster", "c", "", "Filter services by cluster name or ARN.")
	cobraCmd.Flags().BoolVarP(&showARN, "arn", "a", false, "Show service ARN.")
	cobraCmd.Flags().BoolVarP(&showCreatedDate, "created-date", "d", false, "Show created date.")
	cmdutil.AddTagFlag(cobraCmd)

	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by service name.")
	cobraCmd.Flags().BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by service status.")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-status")

	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cobraCmd.Flags().SortFlags = false
}

func ListServices(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, ecs.NewECSService)
	if err != nil {
		return fmt.Errorf("create ECS service: %w", err)
	}

	services, err := svc.GetAllServices(cmd.Context(), cluster)
	if err != nil {
		return fmt.Errorf("list ECS services: %w", err)
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         "ECS Services",
		PlainStyle:    list,
		Fields:        getListFields(),
		Tags:          cmdutil.Tags,
		Data:          utils.SlicesToAny(services),
		GetFieldValue: ecs.GetFieldValue,
		GetTagValue:   ecs.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}
