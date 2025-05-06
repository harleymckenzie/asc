package elb

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/jedib0t/go-pretty/v6/table"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/elb/types"
)

type mockELBClient struct {
	describeLoadBalancersOutput *elbv2.DescribeLoadBalancersOutput
	describeTargetGroupsOutput  *elbv2.DescribeTargetGroupsOutput
	err                         error
}

func (m *mockELBClient) DescribeLoadBalancers(
	_ context.Context,
	_ *elbv2.DescribeLoadBalancersInput,
	_ ...func(*elbv2.Options),
) (*elbv2.DescribeLoadBalancersOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.describeLoadBalancersOutput, nil
}

func (m *mockELBClient) DescribeTargetGroups(
	_ context.Context,
	_ *elbv2.DescribeTargetGroupsInput,
	_ ...func(*elbv2.Options),
) (*elbv2.DescribeTargetGroupsOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.describeTargetGroupsOutput, nil
}

func TestListLoadBalancers(t *testing.T) {
	testCases := []struct {
		name          string
		loadBalancers []types.LoadBalancer
		err           error
		wantErr       bool
	}{
		{
			name: "mixed load balancer types",
			loadBalancers: []types.LoadBalancer{
				{
					LoadBalancerName: aws.String("test-lb-1"),
					LoadBalancerArn:  aws.String("arn:aws:elasticloadbalancing:region:account:loadbalancer/app/test-lb-1/1234567890"),
					DNSName:          aws.String("test-lb-1.region.elb.amazonaws.com"),
					Type:             types.LoadBalancerTypeEnumApplication,
					Scheme:           types.LoadBalancerSchemeEnumInternal,
					State: &types.LoadBalancerState{
						Code: types.LoadBalancerStateEnumActive,
					},
					VpcId: aws.String("vpc-1234567890"),
				},
				{
					LoadBalancerName: aws.String("test-lb-2"),
					LoadBalancerArn:  aws.String("arn:aws:elasticloadbalancing:region:account:loadbalancer/net/test-lb-2/0987654321"),
					DNSName:          aws.String("test-lb-2.region.elb.amazonaws.com"),
					Type:             types.LoadBalancerTypeEnumNetwork,
					Scheme:           types.LoadBalancerSchemeEnumInternetFacing,
					State: &types.LoadBalancerState{
						Code: types.LoadBalancerStateEnumActive,
					},
					VpcId: aws.String("vpc-0987654321"),
				},
			},
			err:     nil,
			wantErr: false,
		},
		{
			name:          "empty response",
			loadBalancers: []types.LoadBalancer{},
			err:           nil,
			wantErr:       false,
		},
		{
			name:          "api error",
			loadBalancers: nil,
			err:           errors.New("Invalid load balancer"),
			wantErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockELBClient{
				describeLoadBalancersOutput: &elbv2.DescribeLoadBalancersOutput{
					LoadBalancers: tc.loadBalancers,
				},
				err: tc.err,
			}

			svc := &ELBService{
				Client: mockClient,
			}

			loadBalancers, err := svc.GetLoadBalancers(context.Background(), &ascTypes.GetLoadBalancersInput{})
			if (err != nil) != tc.wantErr {
				t.Errorf("GetLoadBalancers() error = %v, wantErr %v", err, tc.wantErr)
			}

			if len(loadBalancers) != len(tc.loadBalancers) {
				t.Errorf("GetLoadBalancers() returned %d load balancers, want %d", len(loadBalancers), len(tc.loadBalancers))
			}
		})
	}
}

func TestListTargetGroups(t *testing.T) {
	testCases := []struct {
		name         string
		targetGroups []types.TargetGroup
		err          error
		wantErr      bool
	}{
		{
			name: "mixed target group types",
			targetGroups: []types.TargetGroup{
				{
					TargetGroupName: aws.String("test-tg-1"),
					TargetGroupArn:  aws.String("arn:aws:elasticloadbalancing:region:account:targetgroup/test-tg-1/1234567890"),
					Protocol:        types.ProtocolEnumHttp,
					Port:            aws.Int32(80),
					TargetType:      types.TargetTypeEnumInstance,
					VpcId:           aws.String("vpc-1234567890"),
				},
				{
					TargetGroupName: aws.String("test-tg-2"),
					TargetGroupArn:  aws.String("arn:aws:elasticloadbalancing:region:account:targetgroup/test-tg-2/0987654321"),
					Protocol:        types.ProtocolEnumTcp,
					Port:            aws.Int32(443),
					TargetType:      types.TargetTypeEnumIp,
					VpcId:           aws.String("vpc-0987654321"),
				},
			},
			err:     nil,
			wantErr: false,
		},
		{
			name:         "empty response",
			targetGroups: []types.TargetGroup{},
			err:          nil,
			wantErr:      false,
		},
		{
			name:         "api error",
			targetGroups: nil,
			err:          errors.New("Invalid target group"),
			wantErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockELBClient{
				describeTargetGroupsOutput: &elbv2.DescribeTargetGroupsOutput{
					TargetGroups: tc.targetGroups,
				},
				err: tc.err,
			}

			svc := &ELBService{
				Client: mockClient,
			}

			targetGroups, err := svc.GetTargetGroups(context.Background(), &ascTypes.GetTargetGroupsInput{})
			if (err != nil) != tc.wantErr {
				t.Errorf("GetTargetGroups() error = %v, wantErr %v", err, tc.wantErr)
			}

			if len(targetGroups) != len(tc.targetGroups) {
				t.Errorf("GetTargetGroups() returned %d target groups, want %d", len(targetGroups), len(tc.targetGroups))
			}
		})
	}
}

func TestTableOutput(t *testing.T) {
	loadBalancers := []types.LoadBalancer{
		{
			LoadBalancerName: aws.String("test-lb-1"),
			DNSName:          aws.String("test-lb-1.region.elb.amazonaws.com"),
			Type:             types.LoadBalancerTypeEnumApplication,
			State: &types.LoadBalancerState{
				Code: types.LoadBalancerStateEnumActive,
			},
		},
	}

	targetGroups := []types.TargetGroup{
		{
			TargetGroupName: aws.String("test-tg-1"),
			Protocol:        types.ProtocolEnumHttp,
			Port:            aws.Int32(80),
			TargetType:      types.TargetTypeEnumInstance,
		},
	}

	testCases := []struct {
		name            string
		selectedColumns []string
		wantHeaders     table.Row
		wantRowCount    int
		testTargetGroup bool
	}{
		{
			name:            "load balancer basic details",
			selectedColumns: []string{"Name", "DNS Name", "Type", "State"},
			wantHeaders:     table.Row{"Name", "DNS Name", "Type", "State"},
			wantRowCount:    1,
			testTargetGroup: false,
		},
		{
			name:            "target group basic details",
			selectedColumns: []string{"Name", "Protocol", "Port", "Target Type"},
			wantHeaders:     table.Row{"Name", "Protocol", "Port", "Target Type"},
			wantRowCount:    1,
			testTargetGroup: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var headers table.Row
			var rows []table.Row

			if tc.testTargetGroup {
				tgTable := &ELBTargetGroupTable{
					TargetGroups:    targetGroups,
					SelectedColumns: tc.selectedColumns,
				}
				headers = tgTable.Headers()
				rows = tgTable.Rows()
			} else {
				lbTable := &ELBTable{
					LoadBalancers:   loadBalancers,
					SelectedColumns: tc.selectedColumns,
				}
				headers = lbTable.Headers()
				rows = lbTable.Rows()
			}

			if len(headers) != len(tc.wantHeaders) {
				t.Errorf("Headers() returned %d columns, want %d", len(headers), len(tc.wantHeaders))
			}

			if len(rows) != tc.wantRowCount {
				t.Errorf("Rows() returned %d rows, want %d", len(rows), tc.wantRowCount)
			}

			// Print table for visual inspection
			tw := table.NewWriter()
			tw.AppendHeader(headers)
			tw.AppendRows(rows)
			t.Logf("\nTable Output:\n%s", tw.Render())
		})
	}
}
