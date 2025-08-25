package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
)

func createEC2Service(cmd *cobra.Command) (*ec2.EC2Service, error) {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	return ec2.NewEC2Service(ctx, profile, region)
}

func getInstances(svc *ec2.EC2Service, args []string) ([]types.Instance, error) {
	ctx := context.TODO()
	return svc.GetInstances(ctx, &ascTypes.GetInstancesInput{
		InstanceIDs: args,
	})
}


