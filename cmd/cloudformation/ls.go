// The ls command lists all CloudFormation stacks.

package cloudformation

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/cloudformation"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/internal/service/cloudformation/types"
)

// Variables
var (
	list            bool
	showLastUpdated bool
	showDescription bool
	sortName        bool
	sortStatus      bool
	sortLastUpdate  bool
	reverseSort     bool
)

// Init function
func init() {
	addLsFlags(lsCmd)
}

// Column functions
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Stack Name", Category: "CloudFormation", Visible: true, SortBy: sortName, SortDirection: tablewriter.Asc},
		{Name: "Status", Category: "CloudFormation", Visible: true, SortBy: sortStatus, SortDirection: tablewriter.Asc},
		{Name: "Description", Category: "CloudFormation", Visible: showDescription},
		{Name: "Last Updated", Category: "CloudFormation", Visible: true, SortBy: sortLastUpdate, SortDirection: tablewriter.Desc},
	}
}

// Command variable
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all CloudFormation stacks",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListCloudFormationStacks(cobraCmd, args))
	},
}

// Flag function
func addLsFlags(lsCmd *cobra.Command) {
	lsCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs CloudFormation stacks in list format.")
	lsCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	lsCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending CloudFormation stack name.")
	lsCmd.Flags().BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by descending CloudFormation stack status.")
	lsCmd.Flags().BoolVarP(&sortLastUpdate, "sort-last-update", "u", false, "Sort by descending CloudFormation stack last updated date.")
	lsCmd.Flags().BoolVarP(&showDescription, "show-description", "d", false, "Show the description of the CloudFormation stack.")
	lsCmd.Flags().BoolVarP(&showLastUpdated, "show-last-updated", "U", false, "Show the last updated date of the CloudFormation stack.")
}

// Command functions
func ListCloudFormationStacks(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, cloudformation.NewCloudFormationService)
	if err != nil {
		return fmt.Errorf("create new CloudFormation service: %w", err)
	}

	stacks, err := svc.GetStacks(cmd.Context(), &ascTypes.GetStacksInput{
		StackName: nil,
	})
	if err != nil {
		return fmt.Errorf("list CloudFormation stacks: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Stacks",
	})
	if list {
		table.SetRenderStyle("plain")
	}

	fields := getListFields()
	fields = cmdutil.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(stacks))

	headerRow := cmdutil.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(cmdutil.BuildRows(utils.SlicesToAny(stacks), fields, cloudformation.GetFieldValue, cloudformation.GetTagValue))
	table.SortBy(fields, reverseSort)
	table.Render()
	return nil
}
