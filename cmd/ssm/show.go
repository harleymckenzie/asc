package ssm

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/service/ssm"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

// Variables
var (
	decrypt   bool
	valueOnly bool
)

// Init function
func init() {
	newShowFlags(showCmd)
}

// getShowFields returns a list of Field objects for the parameter detail view.
func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Parameter Details", Visible: true},
		{Name: "Type", Category: "Parameter Details", Visible: true},
		{Name: "Value", Category: "Parameter Details", Visible: true},
		{Name: "Version", Category: "Parameter Details", Visible: true},
		{Name: "Last Modified Date", Category: "Parameter Details", Visible: true},
		{Name: "ARN", Category: "Parameter Details", Visible: true},
		{Name: "Data Type", Category: "Parameter Details", Visible: true},
	}
}

var showCmd = &cobra.Command{
	Use:     "show <parameter-name>",
	Short:   "Show detailed information about an SSM parameter",
	Aliases: []string{"describe", "get"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	Example: "  asc ssm show /myapp/prod/db-password\n" +
		"  asc ssm show /myapp/prod/db-password --decrypt\n" +
		"  asc ssm show /myapp/prod/db-password --value-only --decrypt",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowSSMParameter(cmd, args[0]))
	},
}

// newShowFlags configures the flags for the show command.
func newShowFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&decrypt, "decrypt", "d", false, "Decrypt and show SecureString values.")
	cmd.Flags().BoolVarP(&valueOnly, "value-only", "v", false, "Print only the parameter value (useful for scripting).")
	cmdutil.AddShowFlags(cmd, "horizontal")
}

// ShowSSMParameter displays detailed information about a parameter.
func ShowSSMParameter(cmd *cobra.Command, name string) error {
	ctx := context.TODO()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	// When value-only is set, always decrypt to get the actual value
	shouldDecrypt := decrypt || valueOnly
	param, err := svc.GetParameter(ctx, &ascTypes.GetParameterInput{
		Name:    name,
		Decrypt: shouldDecrypt,
	})
	if err != nil {
		return fmt.Errorf("get parameter: %w", err)
	}

	// If value-only flag is set, just print the value and exit
	if valueOnly {
		fmt.Println(aws.ToString(param.Value))
		return nil
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:          "Parameter: " + aws.ToString(param.Name),
		Columns:        2,
		MaxColumnWidth: 100,
	})

	// Create field getter that respects decrypt flag
	fieldGetter := func(fieldName string, instance any) (string, error) {
		if fieldName == "Value" && decrypt {
			p := instance.(types.Parameter)
			return ssm.GetDecryptedValue(p), nil
		}
		return ssm.GetFieldValue(fieldName, instance)
	}

	fields, err := tablewriter.PopulateFieldValues(*param, getShowFields(), fieldGetter)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	// Layout = Horizontal or Grid
	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}
	table.AddSections(tablewriter.BuildSections(fields, layout))

	table.Render()
	return nil
}
