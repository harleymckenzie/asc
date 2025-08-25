// The show command displays detailed information about an EC2 instance.

package ec2

import (
	"fmt"
	"strings"

	"github.com/harleymckenzie/asc/cmd/ec2/ami"
	"github.com/harleymckenzie/asc/cmd/ec2/snapshot"
	"github.com/harleymckenzie/asc/cmd/ec2/volume"
	"github.com/harleymckenzie/asc/internal/service/ec2"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter/builder"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {
	newShowFlags(showCmd)
}

// getShowFields returns a list of Field objects for the given instance.
func getShowFields() []builder.Field {
	return []builder.Field{
		{Name: "Instance ID", Category: "Instance Details", Visible: true},
		{Name: "State", Category: "Instance Details", Visible: true},
		{Name: "AMI ID", Category: "Instance Details", Visible: true},
		{Name: "AMI Name", Category: "Instance Details", Visible: true},
		{Name: "Launch Time", Category: "Instance Details", Visible: true},
		{Name: "Instance Type", Category: "Instance Details", Visible: true},
		{Name: "Placement Group", Category: "Instance Details", Visible: true},
		{Name: "Root Device Type", Category: "Instance Details", Visible: true},
		{Name: "Root Device Name", Category: "Instance Details", Visible: true},
		{Name: "Virtualization Type", Category: "Instance Details", Visible: true},
		{Name: "vCPUs", Category: "Instance Details", Visible: true},
		{Name: "Public IP", Category: "Network", Visible: true},
		{Name: "Private IP", Category: "Network", Visible: true},
		{Name: "Subnet ID", Category: "Network", Visible: true},
		{Name: "VPC ID", Category: "Network", Visible: true},
		{Name: "Availability Zone", Category: "Network", Visible: true},
		{Name: "Security Group(s)", Category: "Security", Visible: true},
		{Name: "Key Name", Category: "Security", Visible: true},
	}
}

// Flag function
func newShowFlags(cmd *cobra.Command) {
	cmdutil.AddShowFlags(cmd, "vertical")
}

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

func ShowEC2Instance(cmd *cobra.Command, args []string) error {
	svc, err := createEC2Service(cmd)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	instance, err := getInstances(svc, args)
	if err != nil {
		return fmt.Errorf("get instances: %w", err)
	}

	fields := getShowFields()
	fields, err = ec2.PopulateFieldValues(fields, instance[0])
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}
	tags, err := awsutil.PopulateTagFields(instance[0].Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags from instance: %w", err)
	}

	renderOptions := tablewriter.AscTableRenderOptions{
		Title:          "Instance summary for " + *instance[0].InstanceId,
		Style:          "rounded",
		Columns:        3,
		MinColumnWidth: 0,
		MaxColumnWidth: 70,
	}
	t := tablewriter.NewAscWriter(renderOptions)

	switch cmdutil.GetLayout(cmd) {
	case "grid":
		appendGridSection(t, fields, tags)
	case "vertical":
		appendVerticalSection(t, fields, tags)
	}

	t.Render()
	return nil
}

func appendGridSection(t tablewriter.AscWriter, fields []builder.Field, tags []builder.Field) {
	builder.AddSections(t, builder.BuildSections(fields, builder.Grid))
	builder.AddSection(t, builder.BuildSection("Tags", tags, builder.Horizontal))
}

func appendVerticalSection(t tablewriter.AscWriter, fields []builder.Field, tags []builder.Field) {
	builder.AddSections(t, builder.BuildSections(fields, builder.Horizontal))
	builder.AddSection(t, builder.BuildSection("Tags", tags, builder.Horizontal))
}