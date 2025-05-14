// The ls command lists all CloudFormation stacks.

package cloudformation

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/cloudformation"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
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
func cloudformationListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Stack Name", Visible: true, Sort: sortName},
		{ID: "Status", Visible: true, Sort: sortStatus},
		{ID: "Description", Visible: showDescription},
		{
			ID:            "Last Updated",
			Visible:       true,
			Sort:          sortLastUpdate,
			DefaultSort:   true,
			SortDirection: "desc",
		},
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
	lsCmd.Flags().
		BoolVarP(&list, "list", "l", false, "Outputs CloudFormation stacks in list format.")
	lsCmd.Flags().
		BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	lsCmd.Flags().
		BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending CloudFormation stack name.")
	lsCmd.Flags().
		BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by descending CloudFormation stack status.")
	lsCmd.Flags().
		BoolVarP(&sortLastUpdate, "sort-last-update", "u", false, "Sort by descending CloudFormation stack last updated date.")
	lsCmd.Flags().
		BoolVarP(&showDescription, "show-description", "d", false, "Show the description of the CloudFormation stack.")
	lsCmd.Flags().
		BoolVarP(&showLastUpdated, "show-last-updated", "U", false, "Show the last updated date of the CloudFormation stack.")
}

// Command functions
func ListCloudFormationStacks(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)

	svc, err := cloudformation.NewCloudFormationService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new CloudFormation service: %w", err)
	}

	stacks, err := svc.GetStacks(ctx, &ascTypes.GetStacksInput{
		StackName: nil,
	})
	if err != nil {
		return fmt.Errorf("list CloudFormation stacks: %w", err)
	}

	fields := cloudformationListFields()

	opts := tableformat.RenderOptions{
		Title:  "Stacks",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(stacks),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return cloudformation.GetAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
