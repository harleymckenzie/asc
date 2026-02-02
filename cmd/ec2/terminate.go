package ec2

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
)

var terminateCmd = &cobra.Command{
	Use:     "terminate",
	Short:   "Terminate an EC2 instance",
	Aliases: []string{"rm", "delete"},
	Example: "asc ec2 terminate i-1234567890abcdef0",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(TerminateEC2Instance(cmd, args))
	},
}

func newTerminateFlags(terminateCmd *cobra.Command) {}

func init() {
	newTerminateFlags(terminateCmd)
}

func TerminateEC2Instance(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)

	if len(args) == 0 {
		cmd.Help()
		return nil
	}

	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	err = svc.TerminateInstance(ctx, &ascTypes.TerminateInstanceInput{
		InstanceID: args[0],
	})
	if err != nil {
		return fmt.Errorf("terminate instance: %w", err)
	}

	return ListEC2Instances(cmd, []string{args[0]})
}
