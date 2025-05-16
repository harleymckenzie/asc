// The ls command lists all EC2 instances.

package ec2

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/cmd/ec2/ami"
	"github.com/harleymckenzie/asc/cmd/ec2/security_group"
	"github.com/harleymckenzie/asc/cmd/ec2/snapshot"
	"github.com/harleymckenzie/asc/cmd/ec2/volume"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
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

	// Add subcommands
	lsCmd.AddCommand(amiLsCmd)
	lsCmd.AddCommand(volumeLsCmd)
	lsCmd.AddCommand(securityGroupLsCmd)
	lsCmd.AddCommand(snapshotLsCmd)

	// Add flags to subcommands
	ami.NewLsFlags(amiLsCmd)
	snapshot.NewLsFlags(snapshotLsCmd)
	volume.NewLsFlags(volumeLsCmd)
	security_group.NewLsFlags(securityGroupLsCmd)

	// Add groups
	lsCmd.AddGroup(cmdutil.SubcommandGroups()...)

}

// Column functions
func ec2ListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Name", Display: true, Merge: false, DefaultSort: true},
		{ID: "Instance ID", Display: true, Sort: sortID},
		{ID: "State", Display: true, Sort: false},
		{ID: "Instance Type", Display: true, Sort: sortType},
		{ID: "Public IP", Display: true, Sort: false},
		{ID: "AMI ID", Display: showAMI, Sort: false},
		{ID: "Launch Time", Display: showLaunchTime, Sort: sortLaunchTime, SortDirection: "desc"},
		{ID: "Private IP", Display: showPrivateIP, Sort: false},
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
		return cmdutil.DefaultErrorHandler(ListEC2Instances(cmd, args))
	},
}

// Subcommands

var amiLsCmd = &cobra.Command{
	Use:     "amis",
	Short:   "List AMIs",
	Aliases: ami.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ami.ListAMIs(cmd, args))
	},
}

var securityGroupLsCmd = &cobra.Command{
	Use:     "security-groups",
	Short:   "List all security groups",
	Aliases: security_group.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(security_group.ListSecurityGroups(cmd, args))
	},
}

var snapshotLsCmd = &cobra.Command{
	Use:     "snapshots",
	Short:   "List all snapshots",
	Aliases: snapshot.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(snapshot.ListSnapshots(cmd, args))
	},
}

var volumeLsCmd = &cobra.Command{
	Use:     "volumes",
	Short:   "List all volumes",
	Aliases: volume.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(volume.ListVolumes(cmd, args))
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
	cobraCmd.Flags().BoolVarP(&sortID, "sort-id", "i", false, "Sort by descending EC2 instance Id.")
	cobraCmd.Flags().
		BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending EC2 instance type.")
	cobraCmd.Flags().
		BoolVarP(&sortLaunchTime, "sort-launch-time", "t", false, "Sort by descending launch time (most recently launched first).")
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
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return ec2.GetAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
