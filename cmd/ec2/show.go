// The show command displays detailed information about an EC2 instance.

package ec2

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/harleymckenzie/asc/cmd/ec2/ami"
	"github.com/harleymckenzie/asc/cmd/ec2/snapshot"
	"github.com/harleymckenzie/asc/cmd/ec2/volume"
	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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
		Title: "EC2 Instance Details",
		Style: "rounded",
		Layout: tableformat.DetailTableLayout{
			Type:          tableformat.DetailTableLayoutAlt,
			ColumnsPerRow: 4,
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

// printInstanceDetailsStyle1 is a function that was created to preview a concept for the table format
// It is not currently used, but is kept here for reference
func printInstanceDetailsStyle1() {
	instance_id := "i-0123456789abcdefg"
	instance_name := "Test Instance"
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight) // Clean look
	t.SetTitle(fmt.Sprintf("Instance details for: %s (%s)", instance_id, instance_name))
	t.AppendHeader(table.Row{"Instance ID", "Public IP(s)", "Private IP(s)"})
	t.AppendRow(table.Row{"i-0123456789abcdefg", "1.2.3.4", "1.2.3.4"})
	t.AppendRow(table.Row{"IPv6 Address", "Instance State", "Instance Type"})
	t.AppendRow(table.Row{"2001:db8:1234:5678:90ab:cdef:1234:5678", "running", "t3.micro"})
	t.AppendRow(table.Row{"Launch Time", "VPC ID", "IAM Role"})
	t.AppendRow(table.Row{"2021-01-01 12:00:00", "vpc-0123456789abcdefg", "ec2-user"})
	t.AppendRow(table.Row{"Subnet ID", "Security Group(s)", "Key Name"})
	t.AppendRow(table.Row{"subnet-0123456789abcdefg", "sg-0123456789abcdefg", "my-key-pair"})

	t.SetStyle(table.StyleRounded)
	t.Style().Options.DrawBorder = true
	t.Style().Options.SeparateColumns = true
	t.Style().Options.SeparateHeader = true
	t.Style().Options.SeparateRows = true

	t.Style().Format.Header = text.FormatDefault

	t.Style().Size.WidthMin = 70
	t.Style().Color.Header = text.Colors{text.Bold}
	t.Style().Color.RowAlternate = text.Colors{text.Bold}
	t.Style().Format.Header = text.FormatDefault
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:        "Instance ID",
			WidthMin:    25,
			WidthMax:    25,
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:        "Public IP(s)",
			WidthMin:    25,
			WidthMax:    25,
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
		{
			Name:        "Private IP(s)",
			WidthMin:    25,
			WidthMax:    25,
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
		},
	})

	t.Render()

}
