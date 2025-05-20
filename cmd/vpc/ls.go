// The ls command lists all VPCs.

package vpc

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/cmd/vpc/subnet"
	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

// Variables
var (
	list               bool
	reverseSort        bool
	sortName           bool
	sortState          bool
	sortIPv4CIDR       bool
	sortIPv6CIDR       bool
	sortOwnerID        bool
	showDHCP           bool
	showMainRouteTable bool
	showMainNetworkACL bool
	showTenancy        bool
)

// Init function
func init() {
	addLsFlags(lsCmd)
}

// Column functions
func vpcListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "VPC ID", Display: true, DefaultSort: true},
		{ID: "State", Display: true, Sort: sortState},
		{ID: "Tenancy", Display: showTenancy},
		{ID: "DHCP Option Set", Display: showDHCP},
		{ID: "Main Route Table", Display: showMainRouteTable},
		{ID: "Main Network ACL", Display: showMainNetworkACL},
		{ID: "IPv4 CIDR", Display: true, Sort: sortIPv4CIDR},
		{ID: "IPv6 CIDR", Display: true, Sort: sortIPv6CIDR},
		{ID: "Default VPC", Display: true},
		{ID: "Owner ID", Display: true},
	}
}

// Command variable
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all VPCs",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListVPCs(cobraCmd, args))
	},
}

// Flag function
func addLsFlags(lsCmd *cobra.Command) {
	lsCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs VPCs in list format.")
	lsCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	lsCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending VPC name.")
	lsCmd.Flags().BoolVarP(&sortState, "sort-state", "s", false, "Sort by descending VPC state.")
	lsCmd.Flags().BoolVarP(&sortIPv4CIDR, "sort-ipv4-cidr", "i", false, "Sort by descending VPC IPv4 CIDR.")
	lsCmd.Flags().BoolVarP(&sortIPv6CIDR, "sort-ipv6-cidr", "I", false, "Sort by descending VPC IPv6 CIDR.")
	lsCmd.Flags().BoolVarP(&sortOwnerID, "sort-owner-id", "o", false, "Sort by descending VPC owner ID.")
	lsCmd.Flags().BoolVarP(&showDHCP, "show-dhcp", "d", false, "Show the DHCP option set for the VPC.")
	lsCmd.Flags().BoolVarP(&showMainRouteTable, "show-main-route-table", "R", false, "Show the main route table for the VPC.")
	lsCmd.Flags().BoolVarP(&showMainNetworkACL, "show-main-network-acl", "N", false, "Show the main network ACL for the VPC.")
	lsCmd.Flags().BoolVarP(&showTenancy, "show-tenancy", "T", false, "Show the tenancy for the VPC.")
	lsCmd.MarkFlagsMutuallyExclusive()
}

// List function
func ListVPCs(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	if len(args) > 0 {
		return ListVPCSubnets(cmd, args)
	} else {
		vpcList, err := svc.GetVPCs(ctx, &ascTypes.GetVPCsInput{})
		if err != nil {
			return fmt.Errorf("list VPCs: %w", err)
		}

		fields := vpcListFields()
		opts := tableformat.RenderOptions{
			Title:  "VPCs",
			Style:  "rounded",
			SortBy: tableformat.GetSortByField(fields, reverseSort),
		}

		if list {
			opts.Style = "list"
		}

		tableformat.RenderTableList(&tableformat.ListTable{
			Instances: utils.SlicesToAny(vpcList),
			Fields:    fields,
			GetAttribute: func(fieldID string, instance any) (string, error) {
				return vpc.GetVPCAttributeValue(fieldID, instance)
			},
		}, opts)

		return nil
	}
}

func ListVPCSubnets(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	subnets, err := svc.GetSubnets(ctx, &ascTypes.GetSubnetsInput{VPCIds: args})
	if err != nil {
		return fmt.Errorf("list subnets: %w", err)
	}

	fields := subnet.SubnetListFields()
	opts := tableformat.RenderOptions{
		Title:  fmt.Sprintf("%s - Subnets", args[0]),
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(subnets),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return vpc.GetSubnetAttributeValue(fieldID, instance)
		},
	}, opts)

	return nil
}
