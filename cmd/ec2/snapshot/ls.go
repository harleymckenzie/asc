// ls.go defines the 'ls' subcommand for snapshot operations.
package snapshot

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	sortID      bool
	sortSize    bool
	showDesc    bool
	reverseSort bool
	owner       string
)

// Init function
func init() {
	NewLsFlags(lsCmd)
}

// ec2SnapshotListFields returns the fields for the snapshot list table.
func ec2SnapshotListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Snapshot ID", Display: true, Sort: sortID},
		{ID: "Volume Size", Display: true},
		{ID: "Volume Size Raw", Hidden: true, Sort: sortSize, SortDirection: "desc"}, // This is used for sorting, as Size is a combination of numbers and letters
		{ID: "Description", Display: showDesc},
		{ID: "Tier", Display: true},
		{ID: "State", Display: true},
		{ID: "Started", Display: true, DefaultSort: true, SortDirection: "desc"},
		{ID: "Progress", Display: true, SortDirection: "desc"},
		{ID: "Encryption", Display: true},
		{ID: "Data Transfer Progress", Display: false, SortDirection: "desc"},
		{ID: "KMS Key ID", Display: false},
		{ID: "Owner ID", Display: true},
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

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs snapshots in list format.")
	cobraCmd.Flags().BoolVarP(&sortID, "sort-id", "i", false, "Sort by descending snapshot ID.")
	cobraCmd.Flags().BoolVarP(&sortSize, "sort-size", "s", false, "Sort by descending snapshot size.")
	cobraCmd.Flags().BoolVarP(&showDesc, "show-description", "d", false, "Show the snapshot description column.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
	cobraCmd.Flags().StringVar(&owner, "owner", "", "Accepts a single AWS account ID or 'all' to show all snapshots. If not provided, only your own snapshots are shown.")
}

// ListSnapshots is the handler for the ls subcommand.
func ListSnapshots(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	// If 'all' is provided, dont use a filter
	// If a specific owner is provided, use the owner-id filter
	// Otherwise, use the self filter
	ownerIds := []string{}
	if owner == "all" {
		// Do nothing
	} else if owner != "" {
		ownerIds = append(ownerIds, owner)
	} else {
		ownerIds = append(ownerIds, "self")
	}

	snapshots, err := svc.GetSnapshots(ctx, &ascTypes.GetSnapshotsInput{
		OwnerIds: ownerIds,
	})
	if err != nil {
		return fmt.Errorf("get snapshots: %w", err)
	}

	fields := ec2SnapshotListFields()
	style := "rounded"
	if list {
		style = "list"
	}
	opts := tableformat.RenderOptions{
		Title:  "Snapshots",
		Style:  style,
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
