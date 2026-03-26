package cloudformation

import (
	"github.com/harleymckenzie/asc/cmd/wait"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var waitCmd = &cobra.Command{
	Use:     "wait <stack-name>",
	Short:   "Wait for a CloudFormation stack to reach a stable state",
	Example: "asc cf wait my-stack",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(WaitStack(cmd, args))
	},
}

func WaitStack(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cmd.Help()
		return nil
	}

	profile, region := cmdutil.GetPersistentFlags(cmd)
	return wait.ExecuteWait(cmd.Context(), profile, region, &awsutil.ResourceURI{
		Service:      "cf",
		ResourceType: "stack",
		Resource:     args[0],
	})
}
