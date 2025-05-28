// The show command displays detailed information about an EC2 instance.

package ec2

import (
	"context"
	"fmt"
	"strings"

	"github.com/harleymckenzie/asc/cmd/ec2/ami"
	"github.com/harleymckenzie/asc/cmd/ec2/snapshot"
	"github.com/harleymckenzie/asc/cmd/ec2/volume"
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
	newShowFlags(showCmd)
}

// Column functions
func ec2ShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Instance Details", Header: true},
		{ID: "Instance ID", Display: true},
		{ID: "State", Display: true},
		{ID: "AMI ID", Display: true},
		{ID: "AMI Name", Display: true},
		{ID: "Launch Time", Display: true},
		{ID: "Instance Type", Display: true},
		{ID: "Placement Group", Display: true},
		{ID: "Root Device Type", Display: true},
		{ID: "Root Device Name", Display: true},
		{ID: "Virtualization Type", Display: true},
		{ID: "vCPUs", Display: true},

		{ID: "Networking", Header: true},
		{ID: "Public IP", Display: true},
		{ID: "Private IP", Display: true},
		{ID: "VPC ID", Display: true},
		{ID: "Subnet ID", Display: true},
		{ID: "Availability Zone", Display: true},

		{ID: "Security", Header: true},
		{ID: "Security Group(s)", Display: true},
		{ID: "Key Name", Display: true},
	}
}

// Flag function
func newShowFlags(cmd *cobra.Command) {}

// Command functions
// ShowEC2Resource displays detailed information for a specified EC2 resource.
// It supports instances, volumes, snapshots, and AMIs.
func ShowEC2Resource(cmd *cobra.Command, arg string) error {
	switch {
	case strings.HasPrefix(arg, "i-"):
		return ShowEC2Instance(cmd, []string{arg})
	case strings.HasPrefix(arg, "vol-"):
		return volume.ShowEC2Volume(cmd, arg)
	case strings.HasPrefix(arg, "snap-"):
		return snapshot.ShowEC2Snapshot(cmd, arg)
	case strings.HasPrefix(arg, "ami-"):
		return ami.ShowEC2AMI(cmd, arg)
	default:
		return fmt.Errorf("invalid resource type: %s", arg)
	}
}

// ShowEC2Instance is the function for showing EC2 instances
func ShowEC2Instance(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	instance, err := svc.GetInstances(ctx, &ascTypes.GetInstancesInput{
		InstanceIDs: args,
	})
	if err != nil {
		return fmt.Errorf("get instances: %w", err)
	}

	fields := ec2ShowFields()
	opts := tableformat.RenderOptions{
		Title: "Instance summary for " + *instance[0].InstanceId,
		Style: "rounded",
		Layout: tableformat.DetailTableLayout{
			Type:          "horizontal",
			ColumnsPerRow: 3,
		},
	}

	err = tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: instance[0],
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return ec2.GetAttributeValue(fieldID, instance)
		},
	}, opts)
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}
	return nil
}
