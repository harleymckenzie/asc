package nat_gateway

import (
	"github.com/harleymckenzie/asc/cmd/wait"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var waitCmd = &cobra.Command{
	Use:     "wait <nat-gateway-id>",
	Short:   "Wait for a NAT gateway to reach a stable state",
	Example: "asc vpc nat-gateway wait nat-1234567890abcdef0",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(WaitNatGateway(cmd, args))
	},
}

func WaitNatGateway(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cmd.Help()
		return nil
	}

	profile, region := cmdutil.GetPersistentFlags(cmd)
	return wait.ExecuteWait(cmd.Context(), profile, region, &awsutil.ResourceURI{
		Service:      "vpc",
		ResourceType: "nat-gateway",
		Resource:     args[0],
	})
}
