package vpc

import (
	"github.com/harleymckenzie/asc/cmd/vpc/igw"
	"github.com/harleymckenzie/asc/cmd/vpc/nacl"
	nat_gateway "github.com/harleymckenzie/asc/cmd/vpc/nat-gateway"
	prefix_list "github.com/harleymckenzie/asc/cmd/vpc/prefix-list"
	route_table "github.com/harleymckenzie/asc/cmd/vpc/route-table"
	"github.com/harleymckenzie/asc/cmd/vpc/subnet"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewVPCRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vpc",
		Short:   "Perform VPC operations",
		GroupID: "service",
	}

	// Add commands
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)
	// cmd.AddCommand(modifyCmd)
	// cmd.AddCommand(rmCmd)

	// Add subcommands
	cmd.AddCommand(igw.NewIGWRootCmd())
	cmd.AddCommand(nacl.NewNACLRootCmd())
	cmd.AddCommand(nat_gateway.NewNatGatewayRootCmd())
	cmd.AddCommand(prefix_list.NewPrefixListRootCmd())
	cmd.AddCommand(route_table.NewRouteTableRootCmd())
	cmd.AddCommand(subnet.NewSubnetRootCmd())

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)

	return cmd
}
