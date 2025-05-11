// The ls command lists all EC2 instances.

package ec2

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/ec2"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/harleymckenzie/asc/pkg/shared/utils"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
)

// Variables
var (
	list           bool
	showAMI        bool
	showLaunchTime bool
	showPrivateIP  bool

	sortName       bool
	sortID         bool
	sortType       bool
	sortLaunchTime bool

	reverseSort bool
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
		{ID: "Name", Visible: true, Sort: sortName, Merge: false, DefaultSort: true},
		{ID: "Instance ID", Visible: true, Sort: sortID},
		{ID: "State", Visible: true, Sort: false},
		{ID: "Instance Type", Visible: true, Sort: sortType},
		{ID: "Public IP", Visible: true, Sort: false},
		{ID: "AMI ID", Visible: showAMI, Sort: false},
		{ID: "Launch Time", Visible: showLaunchTime, Sort: sortLaunchTime},
		{ID: "Private IP", Visible: showPrivateIP, Sort: false},
	}
}

// Command variable
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all EC2 instances",
	Aliases: []string{"list"},
	GroupID: "actions",
	Example: "asc ec2 ls -A           # List all EC2 instances with the AMI ID\n" +
		"asc ec2 ls -Lt          # List all EC2 instances, displaying and sorting by launch time\n" +
		"asc ec2 ls -l           # List all EC2 instances in list format",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ListEC2Instances(cobraCmd, args)
	},
}

// Flag function
func newLsFlags(cobraCmd *cobra.Command) {
	// Add flags - Output
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs EC2 instances in list format.")
	cobraCmd.Flags().BoolVarP(&showAMI, "ami", "A", false, "Show the AMI ID of the instance.")
	cobraCmd.Flags().
		BoolVarP(&showLaunchTime, "launch-time", "L", false, "Show the launch time of the instance.")
	cobraCmd.Flags().
		BoolVarP(&showPrivateIP, "private-ip", "P", false, "Show the private IP address of the instance.")

	// Add flags - Sorting
	cobraCmd.Flags().
		BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending EC2 instance name.")
	cobraCmd.Flags().BoolVarP(&sortID, "sort-id", "i", false, "Sort by descending EC2 instance Id.")
	cobraCmd.Flags().
		BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending EC2 instance type.")
	cobraCmd.Flags().
		BoolVarP(&sortLaunchTime, "sort-launch-time", "t", false, "Sort by descending launch time (most recently launched first).")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-id", "sort-type", "sort-launch-time")
}

// Command functions
// ListEC2Instances is the function for listing EC2 instances
func ListEC2Instances(cobraCmd *cobra.Command, args []string) {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		log.Fatalf("Failed to initialize EC2 service: %v", err)
	}

	instances, err := svc.GetInstances(ctx, &ascTypes.GetInstancesInput{})
	if err != nil {
		log.Fatalf("Failed to list EC2 instances: %v", err)
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
		GetAttribute: func(fieldID string, instance any) string {
			return ec2.GetAttributeValue(fieldID, instance)
		},
	}, opts)
}
