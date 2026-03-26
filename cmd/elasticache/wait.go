package elasticache

import (
	"github.com/harleymckenzie/asc/cmd/wait"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var waitCmd = &cobra.Command{
	Use:     "wait <cluster-id>",
	Short:   "Wait for an ElastiCache cluster to reach a stable state",
	Example: "asc elasticache wait my-cluster",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(WaitCluster(cmd, args))
	},
}

func WaitCluster(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cmd.Help()
		return nil
	}

	profile, region := cmdutil.GetPersistentFlags(cmd)
	return wait.ExecuteWait(cmd.Context(), profile, region, &awsutil.ResourceURI{
		Service:      "elasticache",
		ResourceType: "cluster",
		Resource:     args[0],
	})
}
