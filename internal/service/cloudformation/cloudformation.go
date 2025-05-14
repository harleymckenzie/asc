package cloudformation

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	ascTypes "github.com/harleymckenzie/asc/internal/service/cloudformation/types"
)

// CloudFormationClientAPI is an interface that defines the methods for the CloudFormation client.
type CloudFormationClientAPI interface {
	DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error)
}

// CloudFormationService is a struct that holds the CloudFormation client.
type CloudFormationService struct {
	Client CloudFormationClientAPI
}

//
// Service functions
//

func NewCloudFormationService(ctx context.Context, profile string, region string) (*CloudFormationService, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}

	// Create a new CloudFormation client
	client := cloudformation.NewFromConfig(cfg)

	// Return a new CloudFormation service with the client
	return &CloudFormationService{
		Client: client,
	}, nil
}

func (svc *CloudFormationService) GetStacks(ctx context.Context, input *ascTypes.GetStacksInput) ([]types.Stack, error) {
	output, err := svc.Client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: input.StackName,
	})
	if err != nil {
		return nil, err
	}

	return output.Stacks, nil
}
