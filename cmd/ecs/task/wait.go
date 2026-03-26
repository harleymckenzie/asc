package task

import (
	"github.com/harleymckenzie/asc/cmd/wait"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var waitCluster string

func init() {
	newWaitFlags(waitCmd)
}

func newWaitFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().StringVarP(&waitCluster, "cluster", "c", "", "Cluster name or ARN (required).")
	cobraCmd.MarkFlagRequired("cluster")
}

var waitCmd = &cobra.Command{
	Use:     "wait <task-id>",
	Short:   "Wait for an ECS task to reach a stable state",
	Example: `  asc ecs task wait abc123 --cluster my-cluster`,
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(WaitTask(cmd, args))
	},
}

func WaitTask(cmd *cobra.Command, args []string) error {
	profile, region := cmdutil.GetPersistentFlags(cmd)
	return wait.ExecuteWait(cmd.Context(), profile, region, &awsutil.ResourceURI{
		Service:      "ecs",
		ResourceType: "task",
		Resource:     args[0],
		Params:       map[string]string{"cluster": waitCluster},
	})
}
