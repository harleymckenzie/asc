package task

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/harleymckenzie/asc/internal/service/ecs"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ecs/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

var showCluster string

func init() {
	NewShowFlags(showCmd)
}

func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Task ID", Category: "General", Visible: true},
		{Name: "ARN", Category: "General", Visible: true},
		{Name: "Status", Category: "General", Visible: true},
		{Name: "Desired Status", Category: "General", Visible: true},
		{Name: "Health Status", Category: "General", Visible: true},
		{Name: "Cluster", Category: "General", Visible: true},
		{Name: "Group", Category: "General", Visible: true},

		{Name: "Task Definition", Category: "Task Definition", Visible: true},
		{Name: "Launch Type", Category: "Task Definition", Visible: true},
		{Name: "Platform Version", Category: "Task Definition", Visible: true},
		{Name: "vCPU", Category: "Task Definition", Visible: true},
		{Name: "Memory", Category: "Task Definition", Visible: true},

		{Name: "Connectivity", Category: "Network", Visible: true},

		{Name: "Created At", Category: "Timestamps", Visible: true},
		{Name: "Started At", Category: "Timestamps", Visible: true},
		{Name: "Stopped At", Category: "Timestamps", Visible: true},
		{Name: "Stop Code", Category: "Timestamps", Visible: true},
		{Name: "Stopped Reason", Category: "Timestamps", Visible: true},

		{Name: "Containers", Category: "Containers", Visible: true},
	}
}

var showCmd = &cobra.Command{
	Use:     "show [task-id]",
	Short:   "Show detailed information about an ECS task",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowTask(cmd, args[0]))
	},
}

func NewShowFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&showCluster, "cluster", "c", "", "Cluster name or ARN (required).")
	cmd.MarkFlagRequired("cluster")
	cmdutil.AddShowFlags(cmd, "vertical")
}

func ShowTask(cmd *cobra.Command, taskID string) error {
	svc, err := cmdutil.CreateService(cmd, ecs.NewECSService)
	if err != nil {
		return fmt.Errorf("create ECS service: %w", err)
	}

	tasks, err := svc.DescribeTasks(cmd.Context(), &ascTypes.DescribeTasksInput{
		Cluster: showCluster,
		Tasks:   []string{taskID},
	})
	if err != nil {
		return fmt.Errorf("describe task: %w", err)
	}

	if len(tasks) == 0 {
		return fmt.Errorf("task %s not found", taskID)
	}

	task := tasks[0]

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("ECS Task Details\n(%s)", ecs.ShortARN(aws.ToString(task.TaskArn))),
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(task, getShowFields(), ecs.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))

	tags, err := populateECSTagFields(task.Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}

// populateECSTagFields converts ECS tags to tablewriter fields
func populateECSTagFields(tags []types.Tag) ([]tablewriter.Field, error) {
	var fields []tablewriter.Field
	for _, tag := range tags {
		if tag.Key != nil && tag.Value != nil {
			fields = append(fields, tablewriter.Field{
				Category: "Tag",
				Name:     aws.ToString(tag.Key),
				Value:    aws.ToString(tag.Value),
			})
		}
	}
	return fields, nil
}
