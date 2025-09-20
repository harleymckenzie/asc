// The show command displays detailed information about a CloudFormation stack.

package cloudformation

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/cloudformation"
	ascTypes "github.com/harleymckenzie/asc/internal/service/cloudformation/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {
	NewShowFlags(showCmd)
}

// Column functions
func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Stack ID", Category: "Overview", Visible: true},
		{Name: "Description", Category: "Overview", Visible: true},
		{Name: "Status", Category: "Overview", Visible: true},
		{Name: "Detailed Status", Category: "Overview", Visible: true},
		{Name: "Status Reason", Category: "Overview", Visible: true},
		{Name: "Creation Time", Category: "Timeline", Visible: true},
		{Name: "Last Updated", Category: "Timeline", Visible: true},
		{Name: "Deletion Time", Category: "Timeline", Visible: true},
		{Name: "Deletion Mode", Category: "Timeline", Visible: true},
		{Name: "Drift Status", Category: "Management", Visible: true},
		{Name: "Root Stack", Category: "Hierarchy", Visible: true},
		{Name: "Parent Stack", Category: "Hierarchy", Visible: true},
		{Name: "Termination Protection", Category: "Security", Visible: true},
		{Name: "IAM Role", Category: "Security", Visible: true},
	}
}

// Command variable
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a CloudFormation stack",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowCloudFormationStack(cobraCmd, args))
	},
}

// Flag function
func NewShowFlags(showCmd *cobra.Command) {
	cmdutil.AddShowFlags(showCmd, "horizontal")
}

// Show function
func ShowCloudFormationStack(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, cloudformation.NewCloudFormationService)
	if err != nil {
		return fmt.Errorf("create new CloudFormation service: %w", err)
	}

	stack, err := svc.GetStacks(cmd.Context(), &ascTypes.GetStacksInput{
		StackName: &args[0],
	})
	if err != nil {
		return fmt.Errorf("get stack: %w", err)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("Stack Details\n(%s)", args[0]),
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(stack[0], getShowFields(), cloudformation.GetFieldValue)
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
