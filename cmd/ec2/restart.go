package ec2

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
)

var restartCmd = &cobra.Command{
	Use:     "restart",
	Short:   "Restart an EC2 instance",
	Aliases: []string{"reboot"},
	Example: "asc ec2 restart i-1234567890abcdef0",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(RestartEC2Instance(cmd, args))
	},
}

func RestartEC2Instance(cmd *cobra.Command, args []string) error {
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

	err = svc.RestartInstance(ctx, &ascTypes.RestartInstanceInput{
		InstanceID: args[0],
	})
	if err != nil {
		return fmt.Errorf("restart instance: %w", err)
	}

	fmt.Printf("Reboot request sent to instance %s\n", args[0])
	return nil
}

func newRestartFlags(cmd *cobra.Command) {}

func init() {
	newRestartFlags(restartCmd)
}
