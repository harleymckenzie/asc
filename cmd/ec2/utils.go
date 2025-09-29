package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
)

// CreateEC2Service creates a new EC2 service instance with the specified configuration
func CreateEC2Service(cmd *cobra.Command) (*ec2.EC2Service, error) {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	return ec2.NewEC2Service(ctx, profile, region)
}

// getInstances retrieves EC2 instances based on the provided arguments
// If args is empty, returns all instances. If args contains instance IDs, returns only those instances.
func getInstances(svc *ec2.EC2Service, args []string) ([]types.Instance, error) {
	ctx := context.TODO()
	return svc.GetInstances(ctx, &ascTypes.GetInstancesInput{
		InstanceIDs: args,
	})
}
