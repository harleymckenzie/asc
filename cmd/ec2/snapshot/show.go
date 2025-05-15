// show.go displays detailed information about a snapshot.
package snapshot

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

// ec2SnapshotShowFields returns the fields for the snapshot detail table.
func ec2SnapshotShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Details", Header: true},
		{ID: "Snapshot ID", Visible: true, Sort: sortID},
		{ID: "Owner ID", Visible: true},
		{ID: "Owner Alias", Visible: true},
		{ID: "Description", Visible: showDesc},
		{ID: "Tier", Visible: true},
		{ID: "State", Visible: true},
		{ID: "State Message", Visible: true},
		{ID: "Encryption", Visible: true},
		{ID: "Started", Visible: true, DefaultSort: true, SortDirection: "desc"},
		{ID: "Progress", Visible: true},
		{ID: "Owner ID", Visible: true},
		
		{ID: "Source Volume", Header: true},
		{ID: "Volume ID", Visible: true},
		{ID: "Volume Size", Visible: true},
		
		{ID: "Encryption", Header: true},
		{ID: "Encryption", Visible: true},
		{ID: "KMS Key ID", Visible: true},

		{ID: "Storage Tier", Header: true},
		{ID: "Storage Tier", Visible: true},
		{ID: "Restore Expiry Time", Visible: true},
	}
}

// NewShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {}

// showCmd is the cobra command for showing snapshot details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a snapshot",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowEC2Snapshot(cmd, args[0]))
	},
}

// ShowEC2Snapshot displays detailed information for a specified snapshot.
func ShowEC2Snapshot(cmd *cobra.Command, arg string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	snapshots, err := svc.GetSnapshots(ctx, &ascTypes.GetSnapshotsInput{SnapshotIDs: []string{arg}})
	if err != nil {
		return fmt.Errorf("get snapshots: %w", err)
	}
	if len(snapshots) == 0 {
		return fmt.Errorf("Snapshot not found: %s", arg)
	}

	fields := ec2SnapshotShowFields()
	opts := tableformat.RenderOptions{
		Title: fmt.Sprintf("Snapshot Details\n(%s)", arg),
		Style: "rounded",
	}

	return tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: snapshots[0],
		Fields:   fields,
		GetAttribute: func(fieldID string, snapshot any) (string, error) {
			return ec2.GetSnapshotAttributeValue(fieldID, snapshot)
		},
	}, opts)
}
