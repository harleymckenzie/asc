// The ls command lists all EC2 instances.

package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/service/ec2"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
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

	testing bool
)

type ListInstancesInput struct {
	GetInstancesInput *ascTypes.GetInstancesInput
	SelectedColumns   []string
	TableOpts         tableformat.RenderOptions
}

// Init function
func init() {
	newLsFlags(lsCmd)
}

// Column functions
func ec2ListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Name", Display: true, Merge: false, DefaultSort: true},
		{ID: "Instance ID", Display: true, Sort: sortByID},
		{ID: "State", Display: true, Sort: false},
		{ID: "Instance Type", Display: true, Sort: sortByType},
		{ID: "Public IP", Display: true, Sort: false},
		{ID: "AMI ID", Display: showAMI, Sort: false},
		{ID: "Launch Time", Display: showLaunchTime, Sort: sortByLaunchTime, SortDirection: "desc"},
		{ID: "Private IP", Display: showPrivateIP, Sort: false},
	}
}

// getShowFields returns a list of Field objects for the given instance.
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Instance Details", Visible: true},
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

// Command variable
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
		if testing {
			return cmdutil.DefaultErrorHandler(ListEC2InstancesV2(cmd, args))
		}
		return cmdutil.DefaultErrorHandler(ListEC2Instances(cmd, args))
	},
}

// Flag function
func newLsFlags(cobraCmd *cobra.Command) {
	// Add flags - Output
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs EC2 instances in list format.")
	cobraCmd.Flags().BoolVarP(&showAMI, "ami", "A", false, "Show the AMI ID of the instance.")
	cobraCmd.Flags().BoolVarP(&showLaunchTime, "launch-time", "L", false, "Show the launch time of the instance.")
	cobraCmd.Flags().BoolVarP(&showPrivateIP, "private-ip", "P", false, "Show the private IP address of the instance.")
	cmdutil.AddTagFlag(cobraCmd)

	// Experimental flags
	cobraCmd.Flags().BoolVar(&testing, "testing", false, "Enable experimental features.")

	// Add flags - Sorting
	cobraCmd.Flags().BoolVarP(&sortByID, "sort-id", "i", false, "Sort by descending EC2 instance Id.")
	cobraCmd.Flags().BoolVarP(&sortByType, "sort-type", "T", false, "Sort by descending EC2 instance type.")
	cobraCmd.Flags().BoolVarP(&sortByLaunchTime, "sort-launch-time", "t", false, "Sort by descending launch time (most recently launched first).")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-id", "sort-type", "sort-launch-time")
}

// Command functions

// ListEC2Instances is the function for listing EC2 instances
func ListEC2Instances(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)

	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	instances, err := svc.GetInstances(ctx, &ascTypes.GetInstancesInput{})
	if err != nil {
		return fmt.Errorf("list ec2 instances: %w", err)
	}

	fields := ec2ListFields()
	opts := tableformat.RenderOptions{
		Title:  "Instances",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(instances),
		Fields:    fields,
		Tags:      cmdutil.Tags,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return ec2.GetAttributeValue(fieldID, instance)
		},
		GetTagValue: func(tag string, instance any) (string, error) {
			return ec2.GetTagValue(tag, instance)
		},
	}, opts)
	return nil
}

func ListEC2InstancesV2(cmd *cobra.Command, args []string) error {
	svc, err := createEC2Service(cmd)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	instances, err := getInstances(svc, args)
	if err != nil {
		return fmt.Errorf("get instances: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Instances",
		Style: "rounded",
	})
	fields := getListFields()
	fields = appendTagFields(fields, cmdutil.Tags, instances)

	appendHeaders(table, fields)
	err = appendRows(table, instances, fields)
	if err != nil {
		return fmt.Errorf("append rows: %w", err)
	}
	table.SortBy(fields, reverseSort)

	table.Render()
	return nil
}

// appendTagFields appends tag fields to the fields slice.
func appendTagFields(fields []tablewriter.Field, tags []string, instances []types.Instance) []tablewriter.Field {
	for _, tag := range tags {
		fields = append(fields, tablewriter.Field{Name: tag, Category: "Tags", Visible: true})
	}
	return fields
}

func appendHeaders(t tablewriter.AscWriter, fields []tablewriter.Field) {
	headerRow := tablewriter.Row{
		Values: make([]string, 0, len(fields)),
	}
	for _, field := range fields {
		if field.Visible {
			headerRow.Values = append(headerRow.Values, field.Name)
		}
	}
	t.AppendHeader(headerRow.Values)
}

func appendRows(t tablewriter.AscWriter, instances []types.Instance, fields []tablewriter.Field) error {
	for _, instance := range instances {
		instanceRow := tablewriter.Row{
			Values: make([]string, 0, len(fields)),
		}
		for _, field := range fields {
			if field.Visible {
				if field.Category == "Tags" {
					fieldValue, err := ec2.GetTagValue(field.Name, instance)
					if err != nil {
						return fmt.Errorf("get tag value: %w", err)
					}
					instanceRow.Values = append(instanceRow.Values, fieldValue)
				} else {
					fieldValue := ec2.GetFieldValue(field.Name, instance)
					instanceRow.Values = append(instanceRow.Values, fieldValue)
				}
			}
		}
		t.AppendRow(instanceRow)
	}
	return nil
}
