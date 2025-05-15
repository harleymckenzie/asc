// The show command displays detailed information about an EC2 volume.

package volume

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {
	NewShowFlags(showCmd)
}

// Column functions
func ec2VolumeShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Volume ID", Visible: true},
		{ID: "Type", Visible: true},
		{ID: "Size", Visible: true},
		{ID: "State", Visible: true},
		{ID: "IOPS", Visible: true},
		{ID: "Throughput", Visible: true},
		{ID: "Fast Snapshot Restored", Visible: true},
		{ID: "Availability Zone", Visible: true},
		{ID: "Created", Visible: true},
		{ID: "Multi-Attach Enabled", Visible: true},

		{ID: "Associations", Header: true},
		{ID: "Snapshot ID", Visible: true},
		{ID: "Associated Resource", Visible: true},
		{ID: "Attach Time", Visible: true},
		{ID: "Delete on Termination", Visible: true},
		{ID: "Device", Visible: true},
		{ID: "Instance ID", Visible: true},
		{ID: "Attachment State", Visible: true},

		{ID: "Encryption", Header: true},
		{ID: "Encryption", Visible: true},
		{ID: "KMS Key ID", Visible: true},
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
func NewShowFlags(cobraCmd *cobra.Command) {}

// ShowEC2Volume is the function for showing EC2 volumes
func ShowEC2Volume(cobraCmd *cobra.Command, args string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cobraCmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	volume, err := svc.GetVolumes(ctx, &ascTypes.GetVolumesInput{
		VolumeIDs: []string{args},
	})
	if err != nil {
		return fmt.Errorf("get volumes: %w", err)
	}

	fields := ec2VolumeShowFields()
	opts := tableformat.RenderOptions{
		Title: "EC2 Volume Details",
		Style: "rounded",
	}

	err = tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: volume[0],
		Fields:   fields,
		GetAttribute: func(fieldID string, volume any) (string, error) {
			return ec2.GetVolumeAttributeValue(fieldID, volume)
		},
	}, opts)
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}
	return nil
}
