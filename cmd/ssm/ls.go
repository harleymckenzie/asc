package ssm

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ssm"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

// Variables
var (
	list        bool
	sortByName  bool
	sortByDate  bool
	reverseSort bool
)

// Init function
func init() {
	newLsFlags(lsCmd)
}

// getListFields returns a list of Field objects for displaying SSM parameters.
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Parameter Details", Visible: true, DefaultSort: true, SortBy: sortByName, SortDirection: tablewriter.Asc},
		{Name: "Type", Category: "Parameter Details", Visible: true},
		{Name: "Last Modified Date", Category: "Parameter Details", Visible: true, SortBy: sortByDate, SortDirection: tablewriter.Desc},
		{Name: "Last Modified User", Category: "Parameter Details", Visible: false},
		{Name: "Version", Category: "Parameter Details", Visible: true},
		{Name: "Tier", Category: "Parameter Details", Visible: false},
		{Name: "Description", Category: "Parameter Details", Visible: false},
	}
}

var lsCmd = &cobra.Command{
	Use:     "ls [path]",
	Short:   "List SSM parameters",
	Aliases: []string{"list"},
	GroupID: "actions",
	Example: "  asc ssm ls                    # List all parameters\n" +
		"  asc ssm ls /myapp/            # List parameters under path\n" +
		"  asc ssm ls /myapp/prod/       # List parameters under specific path",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListSSMParameters(cmd, args))
	},
}

// newLsFlags configures the flags for the ls command.
func newLsFlags(cobraCmd *cobra.Command) {
	// Output format flags
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Output parameters in list format.")

	// Sorting flags
	cobraCmd.Flags().BoolVarP(&sortByName, "sort-name", "n", false, "Sort by parameter name.")
	cobraCmd.Flags().BoolVarP(&sortByDate, "sort-date", "d", false, "Sort by last modified date (most recent first).")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-date")
}

// ListSSMParameters lists SSM parameters with optional path filtering.
func ListSSMParameters(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	// Get path from args if provided
	path := ""
	if len(args) > 0 {
		path = args[0]
	}

	// Use DescribeParameters for listing (doesn't require decryption)
	parameters, err := svc.DescribeParameters(ctx, path)
	if err != nil {
		return fmt.Errorf("describe parameters: %w", err)
	}

	if len(parameters) == 0 {
		if path != "" {
			fmt.Printf("No parameters found under path: %s\n", path)
		} else {
			fmt.Println("No parameters found.")
		}
		return nil
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Parameters",
	})
	if list {
		table.SetRenderStyle("plain")
	}

	fields := getListFields()
	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(parameters), fields, ssm.GetFieldValue, ssm.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)

	table.Render()
	return nil
}
