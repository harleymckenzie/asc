package taskdefinition

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/harleymckenzie/asc/internal/service/ecs"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ecs/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	NewShowFlags(showCmd)
}

func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Family", Category: "General", Visible: true},
		{Name: "Revision", Category: "General", Visible: true},
		{Name: "ARN", Category: "General", Visible: true},
		{Name: "Status", Category: "General", Visible: true},
		{Name: "Registered At", Category: "General", Visible: true},

		{Name: "Network Mode", Category: "Configuration", Visible: true},
		{Name: "Requires Compatibilities", Category: "Configuration", Visible: true},
		{Name: "vCPU", Category: "Configuration", Visible: true},
		{Name: "Memory", Category: "Configuration", Visible: true},

		{Name: "Task Role ARN", Category: "IAM", Visible: true},
		{Name: "Execution Role ARN", Category: "IAM", Visible: true},

		{Name: "Containers", Category: "Containers", Visible: true},
	}
}

var showCmd = &cobra.Command{
	Use:     "show [task-definition]",
	Short:   "Show detailed information about an ECS task definition",
	Long:    "Show detailed information about an ECS task definition. Accepts family:revision or full ARN.",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowTaskDefinition(cmd, args[0]))
	},
}

func NewShowFlags(cmd *cobra.Command) {
	cmdutil.AddShowFlags(cmd, "vertical")
}

func ShowTaskDefinition(cmd *cobra.Command, taskDef string) error {
	svc, err := cmdutil.CreateService(cmd, ecs.NewECSService)
	if err != nil {
		return fmt.Errorf("create ECS service: %w", err)
	}

	td, err := svc.DescribeTaskDefinition(cmd.Context(), &ascTypes.DescribeTaskDefinitionInput{
		TaskDefinition: taskDef,
	})
	if err != nil {
		return fmt.Errorf("describe task definition: %w", err)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("ECS Task Definition Details\n(%s:%d)", aws.ToString(td.Family), td.Revision),
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(*td, getShowFields(), ecs.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))
	table.Render()
	return nil
}
