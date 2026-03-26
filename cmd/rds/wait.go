package rds

import (
	"github.com/harleymckenzie/asc/cmd/wait"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var (
	waitCluster bool
)

func init() {
	newWaitFlags(waitCmd)
}

func newWaitFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().SortFlags = false
	cobraCmd.Flags().BoolVarP(&waitCluster, "cluster", "c", false, "Wait for a cluster instead of an instance")
}

var waitCmd = &cobra.Command{
	Use:   "wait <identifier>",
	Short: "Wait for an RDS instance or cluster to reach a stable state",
	Example: `  asc rds wait my-database
  asc rds wait my-cluster --cluster`,
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(WaitRDSResource(cmd, args))
	},
}

func WaitRDSResource(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cmd.Help()
		return nil
	}

	resourceType := "instance"
	if waitCluster {
		resourceType = "cluster"
	}

	profile, region := cmdutil.GetPersistentFlags(cmd)
	return wait.ExecuteWait(cmd.Context(), profile, region, &awsutil.ResourceURI{
		Service:      "rds",
		ResourceType: resourceType,
		Resource:     args[0],
	})
}
