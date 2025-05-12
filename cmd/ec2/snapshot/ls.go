// ls.go defines the 'ls' subcommand for snapshot operations.
package snapshot

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
	sortSize    bool
	showDesc    bool
	reverseSort bool
)

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs snapshots in list format.")
	cobraCmd.Flags().BoolVarP(&sortID, "sort-id", "i", false, "Sort by descending snapshot ID.")
	cobraCmd.Flags().
		BoolVarP(&sortSize, "sort-size", "s", false, "Sort by descending snapshot size.")
	cobraCmd.Flags().
		BoolVarP(&showDesc, "show-description", "d", false, "Show the snapshot description column.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
}

// ec2SnapshotListFields returns the fields for the snapshot list table.
func ec2SnapshotListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Snapshot ID", Visible: true, Sort: sortID},
		{ID: "Size", Visible: true, Sort: sortSize},
		{ID: "Description", Visible: showDesc},
	}
}

// lsCmd is the cobra command for listing snapshots.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all snapshots",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListSnapshots(cmd, args))
	},
}

// ListSnapshots is the handler for the ls subcommand.
func ListSnapshots(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	snapshots, err := svc.GetSnapshots(ctx, &ascTypes.GetSnapshotsInput{})
	if err != nil {
		return fmt.Errorf("get snapshots: %w", err)
	}

	fields := ec2SnapshotListFields()
	opts := tableformat.RenderOptions{
		Title:  "Snapshots",
		Style:  "list",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(snapshots),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return ec2.GetSnapshotAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
