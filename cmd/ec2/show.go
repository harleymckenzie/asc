// The show command displays detailed information about an EC2 instance.

package ec2

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/smithy-go"
	"github.com/harleymckenzie/asc/pkg/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

// Variables
var (
	apiErr smithy.APIError
	oe     *smithy.OperationError
)

// Init function
func init() {
	newShowFlags(showCmd)
}

func ec2ShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Instance Details", Header: true},
		{ID: "Instance ID", Visible: true},
		{ID: "State", Visible: true},
		{ID: "AMI ID", Visible: true},
		{ID: "AMI Name", Visible: true},
		{ID: "Launch Time", Visible: true},
		{ID: "Instance Type", Visible: true},
		{ID: "Placement Group", Visible: true},
		{ID: "Root Device Type", Visible: true},
		{ID: "Root Device Name", Visible: true},
		{ID: "Virtualization Type", Visible: true},
		{ID: "vCPUs", Visible: true},

		{ID: "Networking", Header: true},
		{ID: "Public IP", Visible: true},
		{ID: "Private IP", Visible: true},
		{ID: "VPC ID", Visible: true},
		{ID: "Subnet ID", Visible: true},
		{ID: "Availability Zone", Visible: true},

		{ID: "Security", Header: true},
		{ID: "Security Group(s)", Visible: true},
		{ID: "Key Name", Visible: true},
	}
}

// Command variable
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an EC2 instance",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	Run: func(cobraCmd *cobra.Command, args []string) {
		ShowEC2Instance(cobraCmd, args)
	},
}

// Flag function
func newShowFlags(cobraCmd *cobra.Command) {
	// Add flags - Output
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs EC2 instances in list format.")
}

// Command functions
// ShowEC2Instance is the function for showing EC2 instances
func ShowEC2Instance(cobraCmd *cobra.Command, args []string) {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")
	fields := ec2ShowFields()

	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			log.Fatalf("Failed to get EC2 service: %s", oe.Unwrap())
		} else {
			log.Fatalf("Failed to get EC2 service: %s", err.Error())
		}
		return
	}

	instance, err := svc.GetInstances(ctx, &ascTypes.GetInstancesInput{
		InstanceIDs: args,
	})
	if err != nil {
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			log.Fatalf("Failed to get EC2 instance: %s", oe.Unwrap())
		} else {
			log.Fatalf("Failed to get EC2 instance: %s", err.Error())
		}
		return
	}

	opts := tableformat.RenderOptions{
		Title: "EC2 Instance Details",
		Style: "rounded",
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: instance[0],
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) string {
			return ec2.GetAttributeValue(fieldID, instance)
		},
	}, opts)
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
