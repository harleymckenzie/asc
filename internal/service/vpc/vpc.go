package vpc

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
)

type VPCAPI interface {
	DescribeVpcs(
		ctx context.Context,
		params *ec2.DescribeVpcsInput,
		optFns ...func(*ec2.Options),
	) (*ec2.DescribeVpcsOutput, error)
}

type VPCService struct {
	Client VPCAPI
}

func NewVPCService(ctx context.Context, profile string, region string) (*VPCService, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)

	return &VPCService{
		Client: client,
	}, nil
}

func (svc *VPCService) GetVPCs(
	ctx context.Context,
	input *ascTypes.GetVPCsInput,
) ([]types.Vpc, error) {
	output, err := svc.Client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, err
	}

	return output.Vpcs, nil
}

// GetNACLs fetches Network ACLs from AWS.
func (svc *VPCService) GetNACLs(
	ctx context.Context,
	input *ascTypes.GetNACLsInput,
) ([]types.NetworkAcl, error) {
	ec2Input := &ec2.DescribeNetworkAclsInput{}
	if input != nil && len(input.NACLIDs) > 0 {
		ec2Input.NetworkAclIds = input.NACLIDs
	}
	output, err := svc.Client.(*ec2.Client).DescribeNetworkAcls(ctx, ec2Input)
	if err != nil {
		return nil, err
	}
	return output.NetworkAcls, nil
}

// GetNatGateways fetches NAT Gateways from AWS.
func (svc *VPCService) GetNatGateways(
	ctx context.Context,
	input *ascTypes.GetNatGatewaysInput,
) ([]types.NatGateway, error) {
	ec2Input := &ec2.DescribeNatGatewaysInput{}
	if input != nil && len(input.NatGatewayIDs) > 0 {
		ec2Input.NatGatewayIds = input.NatGatewayIDs
	}
	output, err := svc.Client.(*ec2.Client).DescribeNatGateways(ctx, ec2Input)
	if err != nil {
		return nil, err
	}
	return output.NatGateways, nil
}

// GetPrefixLists fetches Prefix Lists from AWS.
func (svc *VPCService) GetPrefixLists(
	ctx context.Context,
	input *ascTypes.GetPrefixListsInput,
) ([]types.PrefixList, error) {
	ec2Input := &ec2.DescribePrefixListsInput{}
	if input != nil && len(input.PrefixListIDs) > 0 {
		ec2Input.PrefixListIds = input.PrefixListIDs
	}
	output, err := svc.Client.(*ec2.Client).DescribePrefixLists(ctx, ec2Input)
	if err != nil {
		return nil, err
	}
	return output.PrefixLists, nil
}

// GetRouteTables fetches Route Tables from AWS.
func (svc *VPCService) GetRouteTables(
	ctx context.Context,
	input *ascTypes.GetRouteTablesInput,
) ([]types.RouteTable, error) {
	ec2Input := &ec2.DescribeRouteTablesInput{}
	if input != nil && len(input.RouteTableIDs) > 0 {
		ec2Input.RouteTableIds = input.RouteTableIDs
	}
	output, err := svc.Client.(*ec2.Client).DescribeRouteTables(ctx, ec2Input)
	if err != nil {
		return nil, err
	}
	return output.RouteTables, nil
}

// GetSubnets fetches Subnets from AWS.
func (svc *VPCService) GetSubnets(
	ctx context.Context,
	input *ascTypes.GetSubnetsInput,
) ([]types.Subnet, error) {
	ec2Input := &ec2.DescribeSubnetsInput{}
	if input != nil && len(input.SubnetIDs) > 0 {
		ec2Input.SubnetIds = input.SubnetIDs
	}
	output, err := svc.Client.(*ec2.Client).DescribeSubnets(ctx, ec2Input)
	if err != nil {
		return nil, err
	}
	return output.Subnets, nil
}

// GetIGWs fetches Internet Gateways from AWS.
func (svc *VPCService) GetIGWs(
	ctx context.Context,
	input *ascTypes.GetIGWsInput,
) ([]types.InternetGateway, error) {
	ec2Input := &ec2.DescribeInternetGatewaysInput{}
	if input != nil && len(input.IGWIDs) > 0 {
		ec2Input.InternetGatewayIds = input.IGWIDs
	}
	output, err := svc.Client.(*ec2.Client).DescribeInternetGateways(ctx, ec2Input)
	if err != nil {
		return nil, err
	}
	return output.InternetGateways, nil
}
