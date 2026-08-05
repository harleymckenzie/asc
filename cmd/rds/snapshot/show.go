// show.go displays detailed information about an RDS snapshot.

package snapshot

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/harleymckenzie/asc/internal/service/rds"
	ascTypes "github.com/harleymckenzie/asc/internal/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

// Variables
var (
	showCluster bool
)

// Init function
func init() {
	newShowFlags(showCmd)
}

// Column functions
func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Snapshot Identifier", Category: "Snapshot Details", Visible: true},
		{Name: "Source Identifier", Category: "Snapshot Details", Visible: true},
		{Name: "Type", Category: "Snapshot Details", Visible: true},
		{Name: "Status", Category: "Snapshot Details", Visible: true},
		{Name: "Created Time", Category: "Snapshot Details", Visible: true},

		{Name: "Engine", Category: "Configuration", Visible: true},
		{Name: "Engine Version", Category: "Configuration", Visible: true},
		{Name: "Storage", Category: "Configuration", Visible: true},
		{Name: "Port", Category: "Configuration", Visible: true},

		{Name: "Encryption", Category: "Security & Network", Visible: true},
		{Name: "Availability Zone", Category: "Security & Network", Visible: true},
		{Name: "VPC ID", Category: "Security & Network", Visible: true},
	}
}

// Command variable
var showCmd = &cobra.Command{
	Use:     "show <snapshot-name>",
	Short:   "Show detailed information about an RDS snapshot",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowRDSSnapshot(cmd, args[0]))
	},
}

// Flag function
func newShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "vertical")
	cobraCmd.Flags().BoolVarP(&showCluster, "cluster", "c", false, "Show a cluster snapshot instead of an instance snapshot.")
}

// ShowRDSSnapshot displays detailed information for an RDS snapshot.
func ShowRDSSnapshot(cmd *cobra.Command, name string) error {
	svc, err := cmdutil.CreateService(cmd, rds.NewRDSService)
	if err != nil {
		return fmt.Errorf("create new RDS service: %w", err)
	}

	input := &ascTypes.GetSnapshotsInput{SnapshotIdentifier: name}

	var snapshot any
	var tagList []types.Tag
	if showCluster {
		snapshots, err := svc.GetClusterSnapshots(cmd.Context(), input)
		if err != nil {
			return fmt.Errorf("get cluster snapshot: %w", err)
		}
		if len(snapshots) == 0 {
			return fmt.Errorf("no cluster snapshot found for %s", name)
		}
		snapshot = snapshots[0]
		tagList = snapshots[0].TagList
	} else {
		snapshots, err := svc.GetSnapshots(cmd.Context(), input)
		if err != nil {
			return fmt.Errorf("get snapshot: %w", err)
		}
		if len(snapshots) == 0 {
			return fmt.Errorf("no snapshot found for %s", name)
		}
		snapshot = snapshots[0]
		tagList = snapshots[0].TagList
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("RDS Snapshot Details\n(%s)", name),
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(snapshot, getShowFields(), rds.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))
	table.AddSection(tablewriter.BuildSection("Tags", populateSnapshotTagFields(tagList), tablewriter.Horizontal))

	table.Render()
	return nil
}

// populateSnapshotTagFields converts RDS tags to tablewriter fields.
func populateSnapshotTagFields(tags []types.Tag) []tablewriter.Field {
	var fields []tablewriter.Field
	for _, tag := range tags {
		if tag.Key != nil && tag.Value != nil {
			fields = append(fields, tablewriter.Field{
				Category: "Tag",
				Name:     aws.ToString(tag.Key),
				Value:    aws.ToString(tag.Value),
			})
		}
	}
	return fields
}
