package ec2

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
)

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start an EC2 instance",
	Example: "asc ec2 start i-1234567890abcdef0",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(StartEC2Instance(cmd, args))
	},
}

func newStartFlags(cmd *cobra.Command) {}

func init() {
	newStartFlags(startCmd)
}

func StartEC2Instance(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	profile, region := cmdutil.GetPersistentFlags(cmd)

	if len(args) == 0 {
		cmd.Help()
		return nil
	}

	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	err = svc.StartInstance(ctx, &ascTypes.StartInstanceInput{
		InstanceID: args[0],
	})
	if err != nil {
		return fmt.Errorf("start instance: %w", err)
	}

	return ListEC2Instances(cmd, []string{args[0]})
}
