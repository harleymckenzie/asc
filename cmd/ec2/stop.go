package ec2

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
)

var (
	force bool
)

var stopCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Stop an EC2 instance",
	Aliases: []string{"shutdown", "halt"},
	Example: "asc ec2 stop i-1234567890abcdef0\n" +
		"asc ec2 stop i-1234567890abcdef0 --force",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(StopEC2Instance(cmd, args))
	},
}

func newStopFlags(stopCmd *cobra.Command) {
	stopCmd.Flags().BoolVarP(&force, "force", "f", false, "Force stop the EC2 instance")
}

func StopEC2Instance(cmd *cobra.Command, args []string) error {
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

	err = svc.StopInstance(ctx, &ascTypes.StopInstanceInput{
		InstanceID: args[0],
		Force:      force,
	})
	if err != nil {
		return fmt.Errorf("stop instance: %w", err)
	}

	return ListEC2Instances(cmd, []string{args[0]})
}

func init() {
	newStopFlags(stopCmd)
}
