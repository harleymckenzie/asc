// show.go displays detailed information about a snapshot.
package snapshot

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/pkg/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

// newShowFlags adds flags for the show subcommand.
func newShowFlags(cobraCmd *cobra.Command) {}

// ec2SnapshotShowFields returns the fields for the snapshot detail table.
func ec2SnapshotShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Snapshot ID", Visible: true},
		// Add more fields as needed
	}
}

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
		Title: "Snapshot Details",
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
