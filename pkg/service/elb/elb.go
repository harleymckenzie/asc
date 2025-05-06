package elb

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/elb/types"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/jedib0t/go-pretty/v6/table"
)

type ELBTable struct {
	LoadBalancers   []types.LoadBalancer
	SelectedColumns []string
}

type ELBTargetGroupTable struct {
	TargetGroups    []types.TargetGroup
	SelectedColumns []string
}

type ELBClientAPI interface {
	DescribeLoadBalancers(ctx context.Context, params *elbv2.DescribeLoadBalancersInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeLoadBalancersOutput, error)
	DescribeTargetGroups(ctx context.Context, params *elbv2.DescribeTargetGroupsInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeTargetGroupsOutput, error)
}

type ELBService struct {
	Client ELBClientAPI
}

func availableColumns() map[string]ascTypes.LoadBalancerColumnDef {
	return map[string]ascTypes.LoadBalancerColumnDef{
		"Name": {
			GetValue: func(i *types.LoadBalancer) string {
				return aws.ToString(i.LoadBalancerName)
			},
		},
	}
}

func availableTargetGroupColumns() map[string]ascTypes.TargetGroupColumnDef {
	return map[string]ascTypes.TargetGroupColumnDef{
		"Name": {
			GetValue: func(i *types.TargetGroup) string {
				return aws.ToString(i.TargetGroupName)
			},
		},
		"Health Check Enabled": {
			GetValue: func(i *types.TargetGroup) string {
				if i.HealthCheckEnabled != nil {
					return strconv.FormatBool(*i.HealthCheckEnabled)
				}
				return "N/A"
			},
		},
		"Health Check Path": {
			GetValue: func(i *types.TargetGroup) string {
				return aws.ToString(i.HealthCheckPath)
			},
		},
		"Health Check Port": {
			GetValue: func(i *types.TargetGroup) string {
				if i.HealthCheckPort != nil {
					return aws.ToString(i.HealthCheckPort)
				}
				return "N/A"
			},
		},
		"Health Check Protocol": {
			GetValue: func(i *types.TargetGroup) string {
				return string(i.HealthCheckProtocol)
			},
		},
		"Health Check Timeout": {
			GetValue: func(i *types.TargetGroup) string {
				return strconv.Itoa(int(*i.HealthCheckTimeoutSeconds))
			},
		},
		"Healthy Threshold": {
			GetValue: func(i *types.TargetGroup) string {
				return strconv.Itoa(int(*i.HealthyThresholdCount))
			},
		},
		"Unhealthy Threshold": {
			GetValue: func(i *types.TargetGroup) string {
				return strconv.Itoa(int(*i.UnhealthyThresholdCount))
			},
		},
		"Target Type": {
			GetValue: func(i *types.TargetGroup) string {
				return string(i.TargetType)
			},
		},
		"Target Group ARN": {
			GetValue: func(i *types.TargetGroup) string {
				return aws.ToString(i.TargetGroupArn)
			},
		},
		"VPC ID": {
			GetValue: func(i *types.TargetGroup) string {
				return aws.ToString(i.VpcId)
			},
		},
	}
}

//
// Table functions
//

// Header and Row functions for Load Balancers
func (et *ELBTable) Headers() table.Row {
	return tableformat.BuildHeaders(et.SelectedColumns)
}

func (et *ELBTable) Rows() []table.Row {
	rows := []table.Row{}
	for _, lb := range et.LoadBalancers {
		row := table.Row{}
		for _, colID := range et.SelectedColumns {
			row = append(row, availableColumns()[colID].GetValue(&lb))
		}
		rows = append(rows, row)
	}
	return rows
}

func (et *ELBTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{
		{Name: "Name", WidthMin: 10, WidthMax: 10},
	}
}

func (et *ELBTable) TableStyle() table.Style {
	return table.StyleRounded
}

// Header and Row functions for Target Groups
func (et *ELBTargetGroupTable) Headers() table.Row {
	return tableformat.BuildHeaders(et.SelectedColumns)
}

func (et *ELBTargetGroupTable) Rows() []table.Row {
	rows := []table.Row{}
	for _, tg := range et.TargetGroups {
		row := table.Row{}
		for _, colID := range et.SelectedColumns {
			row = append(row, availableTargetGroupColumns()[colID].GetValue(&tg))
		}
		rows = append(rows, row)
	}
	return rows
}

func (et *ELBTargetGroupTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{
		// {Name: "Name", WidthMin: 10, WidthMax: 10},
	}
}

func (et *ELBTargetGroupTable) TableStyle() table.Style {
	return table.StyleRounded
}

//
// Service functions
//

func NewELBService(ctx context.Context, profile string, region string) (*ELBService, error) {
	var cfg aws.Config
	var err error

	opts := []func(*config.LoadOptions) error{}

	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	cfg, err = config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &ELBService{Client: elbv2.NewFromConfig(cfg)}, nil
}

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
