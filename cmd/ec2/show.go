// The show command displays detailed information about an EC2 instance.

package ec2

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/harleymckenzie/asc/pkg/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
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
	// selectedFields, sortBy, headerFields := tableformat.BuildFields(fields)

	opts := tableformat.RenderOptions{
		List:  list,
		Title: "EC2 Instance Details",
	}

	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		log.Fatalf("Failed to initialize EC2 service: %v", err)
	}

	instance, err := svc.GetInstances(ctx, &ascTypes.GetInstancesInput{
		InstanceIDs: args,
	})
	if err != nil {
		log.Fatalf("Failed to get EC2 instance: %v", err)
	}

	tableformat.RenderDetail(&ec2.EC2DetailTable{
		Instance: instance[0],
		Fields:   fields,
	}, opts)
}

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
		{Name: "Instance ID", WidthMin: 25, WidthMax: 25, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Name: "Public IP(s)", WidthMin: 25, WidthMax: 25, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Name: "Private IP(s)", WidthMin: 25, WidthMax: 25, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
	})

	t.Render()

}

func printInstanceDetailsStyle2() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Options.DrawBorder = true
	t.Style().Options.SeparateColumns = true
	t.Style().Options.SeparateHeader = true
	t.Style().Options.SeparateRows = false
	t.SetTitle("Instance details")

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Colors: text.Colors{text.Bold}},
	})

	t.AppendRow(table.Row{"Instance details", "Instance details"}, table.RowConfig{AutoMerge: true})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Instance ID", "i-0123456789abcdefg"})
	t.AppendRow(table.Row{"State", "running"})
	t.AppendRow(table.Row{"Instance Type", "t3.micro"})
	t.AppendRow(table.Row{"Public IP", "1.2.3.4"})
	t.AppendRow(table.Row{"Private IP", "1.2.3.4"})
	t.AppendRow(table.Row{"AMI ID", "ami-0123456789abcdefg"})
	t.AppendRow(table.Row{"Launch Time", "Tue Sep 05 2017 13:10:17 GMT+0100 (British Summer Time) (over 7 years)"}, table.RowConfig{AutoMerge: true})
	t.AppendRow(table.Row{"Subnet ID", "subnet-0123456789abcdefg"})
	t.AppendRow(table.Row{"Security Group(s)", "sg-0123456789abcdefg"})
	t.AppendRow(table.Row{"Key Name", "my-key-pair"})
	t.AppendRow(table.Row{"VPC ID", "vpc-0123456789abcdefg"})
	t.AppendRow(table.Row{"IAM Role", "ec2-user"})
	t.AppendRow(table.Row{"Placement Group", "default"})
	t.AppendRow(table.Row{"Availability Zone", "us-east-1a"})
	t.AppendRow(table.Row{"Root Device Type", "ebs"})
	t.AppendRow(table.Row{"Root Device Name", "/dev/sda1"})

	t.AppendSeparator()

	// Section 2: Host and placement group
	t.AppendRow(table.Row{"Host and placement group", "Host and placement group"}, table.RowConfig{AutoMerge: true})
	t.AppendSeparator()
	t.AppendRow(table.Row{"Virtualization type", "hvm"})
	t.AppendRow(table.Row{"vCPUs", "1"})

	t.Render()
}

// Brainstorming
// - Detail table should be made up of 3 columns, with no headers
// - Each row represents a key-value pair
// - The first column is the keys (eg,. Name, Instance Id, State)
// - The second column is the value (eg,. my-instance, i-0123456789abcdefg, running)
// - The third column is the key (eg,. Instance Type, Public IP, AMI ID, Launch Time, Private IP)
// - The fourth column is the value (eg,. t3.micro, 1.2.3.4, ami-0123456789abcdefg, 2021-01-01 12:00:00, 1.2.3.4)

// How this data would look:
// []table.Row{
// 	{"Name","Instance ID","State"},
// 	{"my-instance","i-0123456789abcdefg","running"},
// 	{"Instance Type","Public IP","AMI ID","Launch Time","Private IP"},
// 	{"t3.micro","1.2.3.4","ami-0123456789abcdefg","2021-01-01 12:00:00","1.2.3.4"},
// }
