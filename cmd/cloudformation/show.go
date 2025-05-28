// The show command displays detailed information about a CloudFormation stack.

package cloudformation

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/cloudformation"
	ascTypes "github.com/harleymckenzie/asc/internal/service/cloudformation/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {

}

// Column functions
func cloudformationShowFields() []tableformat.Field {
	return []tableformat.Field{
        {ID: "Overview", Display: true, Header: true},
        {ID: "Stack ID", Display: true},
		{ID: "Description", Display: true},
		{ID: "Status", Display: true},
		{ID: "Detailed Status", Display: true},
		{ID: "Status Reason", Display: true},
		{ID: "Creation Time", Display: true},
        {ID: "Last Updated", Display: true},
        {ID: "Deletion Time", Display: true},
        {ID: "Deletion Mode", Display: true},
        {ID: "Drift Status", Display: true},
        {ID: "Root Stack", Display: true},
        {ID: "Parent Stack", Display: true},
        {ID: "Termination Protection", Display: true},
        {ID: "IAM Role", Display: true},
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
func addShowFlags(showCmd *cobra.Command) {}

// Show function
func ShowCloudFormationStack(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	fields := cloudformationShowFields()

	svc, err := cloudformation.NewCloudFormationService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new CloudFormation service: %w", err)
	}

	stack, err := svc.GetStacks(ctx, &ascTypes.GetStacksInput{
		StackName: &args[0],
	})
	if err != nil {
		return fmt.Errorf("get stack: %w", err)
	}

	opts := tableformat.RenderOptions{
		Title: fmt.Sprintf("Stack Details\n(%s)", args[0]),
		Style: "rounded",
		Layout: tableformat.DetailTableLayout{
			Type: "vertical",
			ColumnsPerRow: 3,
			// ColumnMinWidth: 20,
			// ColumnMaxWidth: 80,
		},
	}

	err = tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: stack[0],
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return cloudformation.GetAttributeValue(fieldID, instance)
		},
	}, opts)
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}
	return nil
}
