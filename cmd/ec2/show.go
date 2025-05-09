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

func ec2ShowColumns() []tableformat.Column {
	return []tableformat.Column{
		{ID: "Name", Visible: true},
		{ID: "Instance ID", Visible: true},
		{ID: "State", Visible: true},
		{ID: "Instance Type", Visible: true},
		{ID: "Public IP", Visible: true},
		{ID: "AMI ID", Visible: false},
		{ID: "Launch Time", Visible: false},
		{ID: "Private IP", Visible: true},
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

	// columns := ec2ShowColumns()
	// selectedColumns, sortBy := tableformat.BuildColumns(columns)

	// opts := tableformat.RenderOptions{
	// 	SortBy: sortBy,
	// 	List:   list,
	// 	Title:  "EC2 Instances",
	// }

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

	println(instance[0].InstanceId)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	instance_id := *instance[0].InstanceId
	instance_name := "Test Instance"
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
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateHeader = true
	t.Style().Options.SeparateRows = true

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

	// tableformat.Render(&ec2.EC2DetailTable{
	// 	Instance:        instance[0],
	// 	SelectedColumns: selectedColumns,
	// }, opts)
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
