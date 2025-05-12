// ls.go defines the 'ls' subcommand for AMI operations.
package ami

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/pkg/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/harleymckenzie/asc/pkg/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	sortID      bool
	sortName    bool
	sortState   bool
	showDesc    bool
	reverseSort bool
)

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs AMIs in list format.")
	cobraCmd.Flags().BoolVarP(&sortID, "sort-id", "i", false, "Sort by descending image ID.")
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending image name.")
	cobraCmd.Flags().
		BoolVarP(&sortState, "sort-state", "s", false, "Sort by descending image state.")
	cobraCmd.Flags().
		BoolVarP(&showDesc, "show-description", "d", false, "Show the AMI description column.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
}

// ec2AMIListFields returns the fields for the AMI list table.
func ec2AMIListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "AMI Name", Visible: true, Sort: sortName},
		{ID: "AMI ID", Visible: true, Sort: sortID},
		{ID: "Source", Visible: false},
		{ID: "Owner", Visible: true},
		{ID: "Visibility", Visible: false},
		{ID: "Status", Visible: true, Sort: sortState},
		{ID: "Creation Date", Visible: true, DefaultSort: true, SortDirection: "dsc"},
		{ID: "Platform", Visible: false},
		{ID: "Root Device Type", Visible: false},
		{ID: "Block Devices", Visible: false},
		{ID: "Virtualization", Visible: false},
	}
}

// lsCmd is the cobra command for listing AMIs.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all AMIs",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListAMIs(cmd, args))
	},
}

// ListAMIs is the handler for the ls subcommand.
func ListAMIs(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	images, err := svc.GetImages(ctx, &ascTypes.GetImagesInput{})
	if err != nil {
		return fmt.Errorf("get images: %w", err)
	}

	fields := ec2AMIListFields()
	opts := tableformat.RenderOptions{
		Title:  "AMIs",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(images),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return ec2.GetImageAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
