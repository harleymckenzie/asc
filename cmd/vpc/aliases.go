package vpc

import (
	"github.com/harleymckenzie/asc/cmd/vpc/igw"
	"github.com/harleymckenzie/asc/cmd/vpc/nacl"
	nat_gateway "github.com/harleymckenzie/asc/cmd/vpc/nat-gateway"
	prefix_list "github.com/harleymckenzie/asc/cmd/vpc/prefix-list"
	route_table "github.com/harleymckenzie/asc/cmd/vpc/route-table"
	subnet "github.com/harleymckenzie/asc/cmd/vpc/subnet"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// init initializes the subcommands and show commands for the VPC command.
func init() {
	// Add subcommands (ls aliases)
	lsCmd.AddCommand(igwLsCmd)
	lsCmd.AddCommand(naclLsCmd)
	lsCmd.AddCommand(natGatewayLsCmd)
	lsCmd.AddCommand(prefixListLsCmd)
	lsCmd.AddCommand(routeTableLsCmd)
	lsCmd.AddCommand(subnetLsCmd)

	// Add subcommands (show aliases)
	showCmd.AddCommand(igwShowCmd)
	showCmd.AddCommand(naclShowCmd)
	showCmd.AddCommand(natGatewayShowCmd)
	showCmd.AddCommand(prefixListShowCmd)
	showCmd.AddCommand(routeTableShowCmd)
	showCmd.AddCommand(subnetShowCmd)

	// Add flags
	nacl.NewLsFlags(naclLsCmd)
	nat_gateway.NewLsFlags(natGatewayLsCmd)
	prefix_list.NewLsFlags(prefixListLsCmd)
	route_table.NewLsFlags(routeTableLsCmd)
	subnet.NewLsFlags(subnetLsCmd)

	igw.NewShowFlags(igwShowCmd)
	nacl.NewShowFlags(naclShowCmd)
	nat_gateway.NewShowFlags(natGatewayShowCmd)
	prefix_list.NewShowFlags(prefixListShowCmd)
	route_table.NewShowFlags(routeTableShowCmd)
	subnet.NewShowFlags(subnetShowCmd)

	lsCmd.AddGroup(cmdutil.SubcommandGroups()...)
	showCmd.AddGroup(cmdutil.SubcommandGroups()...)
}

// IGW
var igwLsCmd = &cobra.Command{
	Use:     "igws",
	Short:   "List all Internet Gateways",
	Aliases: igw.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(igw.ListIGWs(cmd, args))
	},
}
var igwShowCmd = &cobra.Command{
	Use:     "igws",
	Short:   "Show detailed information about an Internet Gateway",
	Aliases: igw.CmdAliases,
	GroupID: "subcommands",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(igw.ShowVPCIGW(cmd, args[0]))
	},
}

// NACL
var naclLsCmd = &cobra.Command{
	Use:     "nacls",
	Short:   "List all Network ACLs",
	Aliases: nacl.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(nacl.ListNACLs(cmd, args))
	},
}
var naclShowCmd = &cobra.Command{
	Use:     "nacls",
	Short:   "Show detailed information about a Network ACL",
	Aliases: nacl.CmdAliases,
	GroupID: "subcommands",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(nacl.ShowNACL(cmd, args[0]))
	},
}

// NAT Gateway
var natGatewayLsCmd = &cobra.Command{
	Use:     "nat-gateways",
	Short:   "List all NAT Gateways",
	Aliases: nat_gateway.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(nat_gateway.ListNatGateways(cmd, args))
	},
}
var natGatewayShowCmd = &cobra.Command{
	Use:     "nat-gateways",
	Short:   "Show detailed information about a NAT Gateway",
	Aliases: nat_gateway.CmdAliases,
	GroupID: "subcommands",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(nat_gateway.ShowNatGateway(cmd, args[0]))
	},
}

// Prefix List
var prefixListLsCmd = &cobra.Command{
	Use:     "prefix-lists",
	Short:   "List all Prefix Lists",
	Aliases: prefix_list.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(prefix_list.ListPrefixLists(cmd, args))
	},
}
var prefixListShowCmd = &cobra.Command{
	Use:     "prefix-lists",
	Short:   "Show detailed information about a Prefix List",
	Aliases: prefix_list.CmdAliases,
	GroupID: "subcommands",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(prefix_list.ShowPrefixList(cmd, args[0]))
	},
}

// Route Table
var routeTableLsCmd = &cobra.Command{
	Use:     "route-tables",
	Short:   "List all Route Tables",
	Aliases: route_table.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(route_table.ListRouteTables(cmd, args))
	},
}
var routeTableShowCmd = &cobra.Command{
	Use:     "route-tables",
	Short:   "Show detailed information about a Route Table",
	Aliases: route_table.CmdAliases,
	GroupID: "subcommands",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(route_table.ShowRouteTable(cmd, args[0]))
	},
}

// Subnet
var subnetLsCmd = &cobra.Command{
	Use:     "subnets",
	Short:   "List all Subnets",
	Aliases: subnet.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(subnet.ListSubnets(cmd, args))
	},
}
var subnetShowCmd = &cobra.Command{
	Use:     "subnets",
	Short:   "Show detailed information about a Subnet",
	Aliases: subnet.CmdAliases,
	GroupID: "subcommands",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(subnet.ShowSubnet(cmd, args[0]))
	},
}

