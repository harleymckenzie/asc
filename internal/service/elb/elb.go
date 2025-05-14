package elb

import (
	"context"
	"strings"

	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	ascTypes "github.com/harleymckenzie/asc/internal/service/elb/types"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)

type ELBClientAPI interface {
	DescribeLoadBalancers(ctx context.Context, params *elbv2.DescribeLoadBalancersInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeLoadBalancersOutput, error)
	DescribeTargetGroups(ctx context.Context, params *elbv2.DescribeTargetGroupsInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeTargetGroupsOutput, error)
}

type ELBService struct {
	Client ELBClientAPI
}

//
// Service functions
//

// NewELBService creates a new ELB service.
func NewELBService(ctx context.Context, profile string, region string) (*ELBService, error) {
	cfg, err := awsutil.LoadDefaultConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	return &ELBService{Client: elbv2.NewFromConfig(cfg.Config)}, nil
}

// GetLoadBalancers gets all the load balancers.
func (svc *ELBService) GetLoadBalancers(ctx context.Context, input *ascTypes.GetLoadBalancersInput) ([]types.LoadBalancer, error) {
	output, err := svc.Client.DescribeLoadBalancers(ctx, &elbv2.DescribeLoadBalancersInput{
		Names: input.ListLoadBalancersInput.Names,
	})
	if err != nil {
		return nil, err
	}

	var loadBalancers []types.LoadBalancer
	loadBalancers = append(loadBalancers, output.LoadBalancers...)
	return loadBalancers, nil
}

// GetTargetGroups gets all the target groups.
func (svc *ELBService) GetTargetGroups(ctx context.Context, input *ascTypes.GetTargetGroupsInput) ([]types.TargetGroup, error) {
	output, err := svc.Client.DescribeTargetGroups(ctx, &elbv2.DescribeTargetGroupsInput{
		Names: input.ListTargetGroupsInput.Names,
	})
	if err != nil {
		return nil, err
	}

	var targetGroups []types.TargetGroup
	targetGroups = append(targetGroups, output.TargetGroups...)
	return targetGroups, nil
}

// getTargetGroupLoadBalancer gets the load balancer name from the target group.
func getTargetGroupLoadBalancer(targetGroup types.TargetGroup) string {
	if len(targetGroup.LoadBalancerArns) > 0 {
		name := strings.Split(string(targetGroup.LoadBalancerArns[0]), "/")
		return name[len(name)-2]
	}
	return ""
}
