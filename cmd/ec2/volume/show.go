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
		{ID: "Volume ID", Display: true},
		{ID: "Type", Display: true},
		{ID: "Size", Display: true},
		{ID: "State", Display: true},
		{ID: "IOPS", Display: true},
		{ID: "Throughput", Display: true},
		{ID: "Fast Snapshot Restored", Display: true},
		{ID: "Availability Zone", Display: true},
		{ID: "Created", Display: true},
		{ID: "Multi-Attach Enabled", Display: true},

		{ID: "Associations", Header: true},
		{ID: "Snapshot ID", Display: true},
		{ID: "Associated Resource", Display: true},
		{ID: "Attach Time", Display: true},
		{ID: "Delete on Termination", Display: true},
		{ID: "Device", Display: true},
		{ID: "Instance ID", Display: true},
		{ID: "Attachment State", Display: true},

		{ID: "Encryption", Header: true},
		{ID: "Encryption", Display: true},
		{ID: "KMS Key ID", Display: true},
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
func ShowEC2Volume(cobraCmd *cobra.Command, args string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cobraCmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	if cobraCmd.Flags().Changed("output") {
		if err := cmdutil.ValidateFlagChoice(cobraCmd, "output", cmdutil.ValidLayouts); err != nil {
			return err
		}
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
		Layout: tableformat.DetailTableLayout{
			Type: cmdutil.GetLayout(cobraCmd),
		},
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
