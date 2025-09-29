// The ls command lists all VPCs.

package vpc

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

// Variables
var (
	list               bool
	sortId             bool
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
func getVPCListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "VPC ID", Category: "VPC", Visible: true, DefaultSort: true, SortBy: sortId, SortDirection: tablewriter.Asc},
		{Name: "State", Category: "VPC", Visible: true, SortBy: sortState, SortDirection: tablewriter.Asc},
		{Name: "Tenancy", Category: "VPC", Visible: showTenancy},
		{Name: "DHCP Option Set", Category: "VPC", Visible: showDHCP},
		{Name: "Main Route Table", Category: "VPC", Visible: showMainRouteTable},
		{Name: "Main Network ACL", Category: "VPC", Visible: showMainNetworkACL},
		{Name: "IPv4 CIDR", Category: "VPC", Visible: true, SortBy: sortIPv4CIDR, SortDirection: tablewriter.Asc},
		{Name: "IPv6 CIDR", Category: "VPC", Visible: true, SortBy: sortIPv6CIDR, SortDirection: tablewriter.Asc},
		{Name: "Default VPC", Category: "VPC", Visible: true},
		{Name: "Owner ID", Category: "VPC", Visible: true},
	}
}

func getSubnetListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Subnet ID", Category: "Subnet", Visible: true, SortBy: true, SortDirection: tablewriter.Asc},
		{Name: "VPC ID", Category: "Subnet", Visible: false},
		{Name: "CIDR Block", Category: "Subnet", Visible: true},
		{Name: "Availability Zone", Category: "Subnet", Visible: true},
		{Name: "State", Category: "Subnet", Visible: true},
		{Name: "Available IPs", Category: "Subnet", Visible: true},
		{Name: "Default For AZ", Category: "Subnet", Visible: true},
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
	lsCmd.Flags().BoolVarP(&sortId, "sort-id", "i", false, "Sort by descending VPC ID.")
	lsCmd.Flags().BoolVarP(&sortState, "sort-state", "s", false, "Sort by descending VPC state.")
	lsCmd.Flags().BoolVarP(&sortIPv4CIDR, "sort-ipv4-cidr", "I", false, "Sort by descending VPC IPv4 CIDR.")
	lsCmd.Flags().BoolVarP(&sortIPv6CIDR, "sort-ipv6-cidr", "6", false, "Sort by descending VPC IPv6 CIDR.")
	lsCmd.Flags().BoolVarP(&sortOwnerID, "sort-owner-id", "o", false, "Sort by descending VPC owner ID.")
	lsCmd.Flags().BoolVarP(&showDHCP, "show-dhcp", "d", false, "Show the DHCP option set for the VPC.")
	lsCmd.Flags().BoolVarP(&showMainRouteTable, "show-main-route-table", "R", false, "Show the main route table for the VPC.")
	lsCmd.Flags().BoolVarP(&showMainNetworkACL, "show-main-network-acl", "N", false, "Show the main network ACL for the VPC.")
	lsCmd.Flags().BoolVarP(&showTenancy, "show-tenancy", "T", false, "Show the tenancy for the VPC.")
	lsCmd.MarkFlagsMutuallyExclusive()
}

// List function
func ListVPCs(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	if len(args) > 0 {
		return ListVPCSubnets(cmd, args)
	} else {
		vpcList, err := svc.GetVPCs(cmd.Context(), &ascTypes.GetVPCsInput{})
		if err != nil {
			return fmt.Errorf("list VPCs: %w", err)
		}

		table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
			Title: "VPCs",
		})
		if list {
			table.SetRenderStyle("plain")
		}

		fields := getVPCListFields()
		fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(vpcList))

		headerRow := tablewriter.BuildHeaderRow(fields)
		table.AppendHeader(headerRow)
		table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(vpcList), fields, vpc.GetFieldValue, vpc.GetTagValue))
		table.SetFieldConfigs(fields, reverseSort)
		table.Render()
		return nil
	}
}

func ListVPCSubnets(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	subnets, err := svc.GetSubnets(cmd.Context(), &ascTypes.GetSubnetsInput{VPCIds: args})
	if err != nil {
		return fmt.Errorf("list subnets: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: fmt.Sprintf("%s - Subnets", args[0]),
	})
	if list {
		table.SetRenderStyle("plain")
	}

	fields := getSubnetListFields()
	fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(subnets))

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(subnets), fields, vpc.GetSubnetFieldValue, vpc.GetSubnetTagValue))
	table.SetFieldConfigs(fields, reverseSort)
	table.Render()
	return nil
}
