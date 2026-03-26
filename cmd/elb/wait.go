package elb

import (
	"github.com/harleymckenzie/asc/cmd/wait"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var waitCmd = &cobra.Command{
	Use:     "wait <lb-name>",
	Short:   "Wait for a load balancer to reach a stable state",
	Example: "asc elb wait my-load-balancer",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(WaitLoadBalancer(cmd, args))
	},
}

func WaitLoadBalancer(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cmd.Help()
		return nil
	}

	profile, region := cmdutil.GetPersistentFlags(cmd)
	return wait.ExecuteWait(cmd.Context(), profile, region, &awsutil.ResourceURI{
		Service:      "elb",
		ResourceType: "load-balancer",
		Resource:     args[0],
	})
}
