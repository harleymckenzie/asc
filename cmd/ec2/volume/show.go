// The show command displays detailed information about an EC2 volume.

package volume

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
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

// Column functions
func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Volume ID", Category: "Volume Details", Visible: true},
		{Name: "Type", Category: "Volume Details", Visible: true},
		{Name: "Size", Category: "Volume Details", Visible: true},
		{Name: "State", Category: "Volume Details", Visible: true},
		{Name: "IOPS", Category: "Volume Details", Visible: true},
		{Name: "Throughput", Category: "Volume Details", Visible: true},
		{Name: "Fast Snapshot Restored", Category: "Volume Details", Visible: true},
		{Name: "Availability Zone", Category: "Volume Details", Visible: true},
		{Name: "Created", Category: "Volume Details", Visible: true},
		{Name: "Multi-Attach Enabled", Category: "Volume Details", Visible: true},

		{Name: "Snapshot ID", Category: "Associations", Visible: true},
		{Name: "Associated Resource", Category: "Associations", Visible: true},
		{Name: "Attach Time", Category: "Associations", Visible: true},
		{Name: "Delete on Termination", Category: "Associations", Visible: true},
		{Name: "Device", Category: "Associations", Visible: true},
		{Name: "Instance ID", Category: "Associations", Visible: true},
		{Name: "Attachment State", Category: "Associations", Visible: true},

		{Name: "Encryption", Category: "Encryption", Visible: true},
		{Name: "KMS Key ID", Category: "Encryption", Visible: true},
	}
}

// Command variable
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an EC2 volume",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowEC2Volume(cmd, args[0]))
	},
}

// Flag function
func NewShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "vertical")
}

// ShowEC2Volume is the function for showing EC2 volumes
func ShowEC2Volume(cmd *cobra.Command, arg string) error {
	svc, err := cmdutil.CreateService(cmd, ec2.NewEC2Service)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	volume, err := svc.GetVolumes(cmd.Context(), &ascTypes.GetVolumesInput{
		VolumeIDs: []string{arg},
	})
	if err != nil {
		return fmt.Errorf("get volumes: %w", err)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   "EC2 Volume Details (" + arg + ")",
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(volume[0], getShowFields(), ec2.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))
	table.Render()
	return nil
}
