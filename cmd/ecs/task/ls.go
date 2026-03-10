package task

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
	serviceName string

	sortStatus bool
)

func init() {
	newLsFlags(lsCmd)
}

func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Task ID", Category: "Task", Visible: true, DefaultSort: true},
		{Name: "Service", Category: "Task", Visible: serviceName == ""},
		{Name: "Status", Category: "Task", Visible: true, SortBy: sortStatus, SortDirection: tablewriter.Asc},
		{Name: "Desired Status", Category: "Task", Visible: true},
		{Name: "Task Definition", Category: "Task", Visible: true},
		{Name: "Created At", Category: "Task", Visible: true},
		{Name: "Started At", Category: "Task", Visible: true},
		{Name: "vCPU", Category: "Resources", Visible: true},
		{Name: "Memory", Category: "Resources", Visible: true},
	}
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List ECS tasks",
	Aliases: []string{"list"},
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListTasks(cmd, args))
	},
}

func newLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs in list format.")
	cobraCmd.Flags().StringVarP(&cluster, "cluster", "c", "", "Filter tasks by cluster name or ARN.")
	cobraCmd.Flags().StringVarP(&serviceName, "service", "S", "", "Filter tasks by service name.")
	cmdutil.AddTagFlag(cobraCmd)

	cobraCmd.Flags().BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by task status.")

	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cobraCmd.Flags().SortFlags = false
}

func ListTasks(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, ecs.NewECSService)
	if err != nil {
		return fmt.Errorf("create ECS service: %w", err)
	}

	tasks, err := svc.GetAllTasks(cmd.Context(), cluster, serviceName)
	if err != nil {
		return fmt.Errorf("list ECS tasks: %w", err)
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         "ECS Tasks",
		PlainStyle:    list,
		Fields:        getListFields(),
		Tags:          cmdutil.Tags,
		Data:          utils.SlicesToAny(tasks),
		GetFieldValue: ecs.GetFieldValue,
		GetTagValue:   ecs.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}
