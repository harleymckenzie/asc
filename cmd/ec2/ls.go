// The ls command lists all EC2 instances.

package ec2

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
)

// Variables
var (
	list           bool
	showAMI        bool
	showLaunchTime bool
	showPrivateIP  bool

	sortByID         bool
	sortByType       bool
	sortByLaunchTime bool

	reverseSort bool
)

type ListInstancesInput struct {
	GetInstancesInput *ascTypes.GetInstancesInput
	SelectedColumns   []string
	TableOpts         tablewriter.AscTableRenderOptions
}

// Init function
func init() {
	newLsFlags(lsCmd)
}

// Column functions
// getListFields returns a list of Field objects for displaying EC2 instance information
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Instance Details", Visible: true, DefaultSort: true},
		{Name: "Instance ID", Category: "Instance Details", Visible: true, SortBy: sortByID, SortDirection: tablewriter.Asc},
		{Name: "State", Category: "Instance Details", Visible: true},
		{Name: "AMI ID", Category: "Instance Details", Visible: false},
		{Name: "AMI Name", Category: "Instance Details", Visible: false},
		{Name: "Launch Time", Category: "Instance Details", Visible: false, SortBy: sortByLaunchTime, SortDirection: tablewriter.Desc},
		{Name: "Instance Type", Category: "Instance Details", Visible: true, SortBy: sortByType, SortDirection: tablewriter.Asc},
		{Name: "Placement Group", Category: "Instance Details", Visible: false},
		{Name: "Root Device Type", Category: "Instance Details", Visible: false},
		{Name: "Root Device Name", Category: "Instance Details", Visible: false},
		{Name: "Virtualization Type", Category: "Instance Details", Visible: false},
		{Name: "vCPUs", Category: "Instance Details", Visible: false},
		{Name: "Public IP", Category: "Network", Visible: true},
		{Name: "Private IP", Category: "Network", Visible: false},
		{Name: "Subnet ID", Category: "Network", Visible: false},
		{Name: "VPC ID", Category: "Network", Visible: false},
		{Name: "Availability Zone", Category: "Network", Visible: false},
		{Name: "Security Group(s)", Category: "Security", Visible: false},
		{Name: "Key Name", Category: "Security", Visible: false},
	}
}

// lsCmd is the main command for listing EC2 instances and related resources
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List EC2 instances, AMIs, snapshots, and volumes",
	Aliases: []string{"list"},
	GroupID: "actions",
	Example: "  asc ec2 ls                   # List all EC2 instances\n" +
		"  asc ec2 ls amis              # List all AMIs\n" +
		"  asc ec2 ls security-groups   # List all security groups\n" +
		"  asc ec2 ls snapshots         # List all snapshots\n" +
		"  asc ec2 ls volumes           # List all volumes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListEC2Instances(cmd, args))
	},
}

// newLsFlags configures the flags for the ls command
func newLsFlags(cobraCmd *cobra.Command) {
	// Output format flags
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs EC2 instances in list format.")
	cobraCmd.Flags().BoolVarP(&showAMI, "ami", "A", false, "Show the AMI ID of the instance.")
	cobraCmd.Flags().BoolVarP(&showLaunchTime, "launch-time", "L", false, "Show the launch time of the instance.")
	cobraCmd.Flags().BoolVarP(&showPrivateIP, "private-ip", "P", false, "Show the private IP address of the instance.")
	cmdutil.AddTagFlag(cobraCmd)

	// Sorting flags
	cobraCmd.Flags().BoolVarP(&sortByID, "sort-id", "i", false, "Sort by descending EC2 instance Id.")
	cobraCmd.Flags().BoolVarP(&sortByType, "sort-type", "T", false, "Sort by descending EC2 instance type.")
	cobraCmd.Flags().BoolVarP(&sortByLaunchTime, "sort-launch-time", "t", false, "Sort by descending launch time (most recently launched first).")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-id", "sort-type", "sort-launch-time")
}

// ListEC2Instances handles the listing of EC2 instances and related resources
func ListEC2Instances(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, ec2.NewEC2Service)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	instances, err := getInstances(svc, args)
	if err != nil {
		return fmt.Errorf("get instances: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Instances",
	})
	if list {
		table.SetRenderStyle("plain")
	}
	fields := getListFields()
	fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(instances))

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(instances), fields, ec2.GetFieldValue, ec2.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)

	table.Render()
	return nil
}
