// show.go displays detailed information about a snapshot.
package snapshot

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {
	NewShowFlags(showCmd)
}

func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Snapshot ID", Category: "Snapshot Details", Visible: true, SortBy: sortID, SortDirection: tablewriter.Asc},
		{Name: "Owner ID", Category: "Snapshot Details", Visible: true},
		{Name: "Owner Alias", Category: "Snapshot Details", Visible: true},
		{Name: "Description", Category: "Snapshot Details", Visible: showDesc},
		{Name: "Tier", Category: "Snapshot Details", Visible: true},
		{Name: "State", Category: "Snapshot Details", Visible: true},
		{Name: "Encryption", Category: "Snapshot Details", Visible: true},
		{Name: "Started", Category: "Snapshot Details", Visible: true, SortBy: sortID, SortDirection: tablewriter.Desc},
		{Name: "Progress", Category: "Snapshot Details", Visible: true},
		{Name: "Owner ID", Category: "Snapshot Details", Visible: true},

		{Name: "Source Volume", Category: "Snapshot Details", Visible: true},
		{Name: "Volume ID", Category: "Snapshot Details", Visible: true},
		{Name: "Volume Size", Category: "Snapshot Details", Visible: true},
		
		{Name: "Encryption", Category: "Snapshot Details", Visible: true},
		{Name: "Encryption", Category: "Snapshot Details", Visible: true},
		{Name: "KMS Key ID", Category: "Snapshot Details", Visible: true},

		{Name: "Storage Tier", Category: "Snapshot Details", Visible: true},
		{Name: "Restore Expiry Time", Category: "Snapshot Details", Visible: true},
	}
}

// NewShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "vertical")
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
	svc, err := cmdutil.CreateService(cmd, ec2.NewEC2Service)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	snapshots, err := svc.GetSnapshots(context.TODO(), &ascTypes.GetSnapshotsInput{SnapshotIDs: []string{arg}})
	if err != nil {
		return fmt.Errorf("get snapshots: %w", err)
	}
	if len(snapshots) == 0 {
		return fmt.Errorf("snapshot not found: %s", arg)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:          "Snapshot summary for " + *snapshots[0].SnapshotId,
		Columns:        3,
		MaxColumnWidth: 80,
	})
	fields, err := cmdutil.PopulateFieldValues(snapshots[0], getShowFields(), ec2.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	// Layout = Horizontal or Grid
	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}
	table.AddSections(tablewriter.BuildSections(fields, layout))
	tags, err := awsutil.PopulateTagFields(snapshots[0].Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags from snapshot: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}
