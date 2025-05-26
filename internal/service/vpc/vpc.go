package vpc

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
)

type VPCAPI interface {
	DescribeVpcs(ctx context.Context, params *ec2.DescribeVpcsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error)
	DescribeNetworkAcls(ctx context.Context, params *ec2.DescribeNetworkAclsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkAclsOutput, error)
	DescribeNatGateways(ctx context.Context, params *ec2.DescribeNatGatewaysInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNatGatewaysOutput, error)
	DescribePrefixLists(ctx context.Context, params *ec2.DescribePrefixListsInput, optFns ...func(*ec2.Options)) (*ec2.DescribePrefixListsOutput, error)
	DescribeManagedPrefixLists(ctx context.Context, params *ec2.DescribeManagedPrefixListsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeManagedPrefixListsOutput, error)
	DescribeRouteTables(ctx context.Context, params *ec2.DescribeRouteTablesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error)
	DescribeSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error)
	DescribeInternetGateways(ctx context.Context, params *ec2.DescribeInternetGatewaysInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInternetGatewaysOutput, error)
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

func (svc *VPCService) GetVPCs(ctx context.Context, input *ascTypes.GetVPCsInput) ([]types.Vpc, error) {
	output, err := svc.Client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, err
	}

	return output.Vpcs, nil
}

// GetNACLs fetches Network ACLs from AWS.
func (svc *VPCService) GetNACLs(ctx context.Context, input *ascTypes.GetNACLsInput) ([]types.NetworkAcl, error) {
	output, err := svc.Client.DescribeNetworkAcls(ctx, &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: input.NetworkAclIds,
	})
	if err != nil {
		return nil, err
	}
	return output.NetworkAcls, nil
}

// GetNatGateways fetches NAT Gateways from AWS.
func (svc *VPCService) GetNatGateways(ctx context.Context, input *ascTypes.GetNatGatewaysInput) ([]types.NatGateway, error) {
	output, err := svc.Client.DescribeNatGateways(ctx, &ec2.DescribeNatGatewaysInput{
		NatGatewayIds: input.NatGatewayIds,
	})
	if err != nil {
		return nil, err
	}
	return output.NatGateways, nil
}

// GetPrefixLists fetches Prefix Lists from AWS.
func (svc *VPCService) GetPrefixLists(ctx context.Context, input *ascTypes.GetPrefixListsInput) ([]types.PrefixList, error) {
	output, err := svc.Client.DescribePrefixLists(ctx, &ec2.DescribePrefixListsInput{
		PrefixListIds: input.PrefixListIds,
	})
	if err != nil {
		return nil, err
	}
	return output.PrefixLists, nil
}

// GetManagedPrefixLists fetches Managed Prefix Lists from AWS.
func (svc *VPCService) GetManagedPrefixLists(ctx context.Context, input *ascTypes.GetManagedPrefixListsInput) ([]types.ManagedPrefixList, error) {
	output, err := svc.Client.DescribeManagedPrefixLists(ctx, &ec2.DescribeManagedPrefixListsInput{
		PrefixListIds: input.PrefixListIds,
	})
	if err != nil {
		return nil, err
	}
	return output.PrefixLists, nil
}

// GetRouteTables fetches Route Tables from AWS.
func (svc *VPCService) GetRouteTables(ctx context.Context, input *ascTypes.GetRouteTablesInput) ([]types.RouteTable, error) {
	output, err := svc.Client.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{
		RouteTableIds: input.RouteTableIds,
	})
	if err != nil {
		return nil, err
	}
	return output.RouteTables, nil
}

// GetSubnets fetches Subnets from AWS.
func (svc *VPCService) GetSubnets(ctx context.Context, input *ascTypes.GetSubnetsInput) ([]types.Subnet, error) {
	filters := []types.Filter{}
	for _, vpcId := range input.VPCIds {
		filters = append(filters, types.Filter{
			Name: aws.String("vpc-id"),
			Values: []string{vpcId},
		})
	}
	output, err := svc.Client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}
	return output.Subnets, nil
}

// GetIGWs fetches Internet Gateways from AWS.
func (svc *VPCService) GetIGWs(ctx context.Context, input *ascTypes.GetIGWsInput) ([]types.InternetGateway, error) {
	output, err := svc.Client.DescribeInternetGateways(ctx, &ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: input.IGWIds,
	})
	if err != nil {
		return nil, err
	}
	return output.InternetGateways, nil
}

// FilterNACLRules fetches Network ACL Rules from AWS.
func (svc *VPCService) FilterNACLRules(rules []types.NetworkAclEntry, ingress bool) []types.NetworkAclEntry {
	filteredRules := []types.NetworkAclEntry{}
	for _, rule := range rules {
		if ingress && *rule.Egress {
			filteredRules = append(filteredRules, rule)
		} else if !ingress && !*rule.Egress {
			filteredRules = append(filteredRules, rule)
		}
	}
	return filteredRules
}
