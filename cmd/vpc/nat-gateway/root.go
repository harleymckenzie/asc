package nat_gateway

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"nat-gateways", "natgateway", "natgateways"}

// NewNatGatewayRootCmd returns the root command for NAT Gateway operations.
func NewNatGatewayRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nat-gateway",
		Short:   "Perform NAT Gateway operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)

	return cmd
}
