package ec2

import (
	"github.com/harleymckenzie/asc/cmd/wait"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var waitCmd = &cobra.Command{
	Use:     "wait <instance-id>",
	Short:   "Wait for an EC2 instance to reach a stable state",
	Example: "asc ec2 wait i-1234567890abcdef0",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(WaitEC2Instance(cmd, args))
	},
}

func WaitEC2Instance(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cmd.Help()
		return nil
	}

	profile, region := cmdutil.GetPersistentFlags(cmd)
	return wait.ExecuteWait(cmd.Context(), profile, region, &awsutil.ResourceURI{
		Service:      "ec2",
		ResourceType: "instance",
		Resource:     args[0],
	})
}
