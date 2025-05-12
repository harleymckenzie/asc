package ec2

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/pkg/service/ec2"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
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

func newStartFlags(cobraCmd *cobra.Command) {}

func init() {
	newStartFlags(startCmd)
}

func StartEC2Instance(cobraCmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	if len(args) == 0 {
		cobraCmd.Help()
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

	return ListEC2Instances(cobraCmd, []string{args[0]})
}
