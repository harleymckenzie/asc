package asg

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/jedib0t/go-pretty/v6/table"
)

type mockASGClient struct {
	describeAutoScalingGroupsOutput *autoscaling.DescribeAutoScalingGroupsOutput
	describeScheduledActionsOutput  *autoscaling.DescribeScheduledActionsOutput
	err                             error
}

func (m *mockASGClient) DescribeAutoScalingGroups(
	_ context.Context,
	params *autoscaling.DescribeAutoScalingGroupsInput,
	_ ...func(*autoscaling.Options),
) (*autoscaling.DescribeAutoScalingGroupsOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.describeAutoScalingGroupsOutput, nil
}

func (m *mockASGClient) DescribeScheduledActions(
	_ context.Context,
	_ *autoscaling.DescribeScheduledActionsInput,
	_ ...func(*autoscaling.Options),
) (*autoscaling.DescribeScheduledActionsOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.describeScheduledActionsOutput, nil
}

func TestListInstances(t *testing.T) {
	testCases := []struct {
		name      string
		instances []types.Instance
		err       error
		wantErr   bool
	}{
		{
			name: "mixed instance types and properties",
			instances: []types.Instance{
				{
					InstanceId:     aws.String("i-1234567890abcdef0"),
					InstanceType:   aws.String("t3.micro"),
					LifecycleState: types.LifecycleStateInService,
				},
				{
					InstanceId:     aws.String("i-1234567890abcdef1"),
					InstanceType:   aws.String("t3.small"),
					LifecycleState: types.LifecycleStateInService,
				},
			},
			err:     nil,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := &AutoScalingService{
				Client: &mockASGClient{
					describeAutoScalingGroupsOutput: &autoscaling.DescribeAutoScalingGroupsOutput{
						AutoScalingGroups: []types.AutoScalingGroup{
							{
								Instances: tc.instances,
							},
						},
					},
					err: tc.err,
				},
			}

			instances, err := svc.GetInstances(context.Background(), &GetInstancesInput{
				AutoScalingGroupNames: []string{},
			})
			if (err != nil) != tc.wantErr {
				t.Errorf("ListInstances() error = %v, wantErr %v", err, tc.wantErr)
			}

			if len(instances) != len(tc.instances) {
				t.Errorf("ListInstances() returned %d instances, want %d", len(instances), len(tc.instances))
			}

			for i, instance := range instances {
				if instance.InstanceId != tc.instances[i].InstanceId {
					t.Errorf("ListInstances() returned instance %d with ID %s, want %s", i, *instance.InstanceId, *tc.instances[i].InstanceId)
				}
			}
		})
	}
}

func TestTableOutput(t *testing.T) {
	asgs := []types.AutoScalingGroup{
		{
			AutoScalingGroupName: aws.String("web-asg"),
			DesiredCapacity:      aws.Int32(2),
			MinSize:              aws.Int32(1),
			MaxSize:              aws.Int32(4),
			Instances: []types.Instance{
				{
					InstanceId:     aws.String("i-1234567890"),
					LifecycleState: "InService",
				},
				{
					InstanceId:     aws.String("i-0987654321"),
					LifecycleState: "InService",
				},
			},
		},
		{
			AutoScalingGroupName: aws.String("worker-asg"),
			DesiredCapacity:      aws.Int32(3),
			MinSize:              aws.Int32(2),
			MaxSize:              aws.Int32(6),
			Instances: []types.Instance{
				{
					InstanceId:     aws.String("i-abcdef1234"),
					LifecycleState: "InService",
				},
			},
		},
	}

	testCases := []struct {
		name            string
		selectedColumns []string
		wantHeaders     table.Row
		wantRowCount    int
	}{
		{
			name:            "full ASG details",
			selectedColumns: []string{"Name", "Instances", "Desired", "Min", "Max"},
			wantHeaders:     table.Row{"Name", "Instances", "Desired", "Min", "Max"},
			wantRowCount:    2,
		},
		{
			name:            "minimal columns",
			selectedColumns: []string{"Name", "Instances"},
			wantHeaders:     table.Row{"Name", "Instances"},
			wantRowCount:    2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			asgTable := &AutoScalingTable{
				AutoScalingGroups: asgs,
				SelectedColumns:   tc.selectedColumns,
			}

			// Test Headers
			headers := asgTable.Headers()
			if len(headers) != len(tc.wantHeaders) {
				t.Errorf("Headers() returned %d columns, want %d", len(headers), len(tc.wantHeaders))
			}

			// Test Rows
			rows := asgTable.Rows()
			if len(rows) != tc.wantRowCount {
				t.Errorf("Rows() returned %d rows, want %d", len(rows), tc.wantRowCount)
			}

			// Print the actual table output
			tw := table.NewWriter()
			tw.AppendHeader(headers)
			tw.AppendRows(rows)
			tw.SetStyle(asgTable.TableStyle())
			t.Logf("\nTable Output:\n%s", tw.Render())
		})
	}

	// Test Instance Table Output
	instances := []types.Instance{
		{
			InstanceId:     aws.String("i-1234567890"),
			LifecycleState: "InService",
			InstanceType:   aws.String("t3.micro"),
			LaunchTemplate: &types.LaunchTemplateSpecification{
				LaunchTemplateName: aws.String("web-template"),
			},
			AvailabilityZone: aws.String("us-west-2a"),
		},
	}

	instanceTestCases := []struct {
		name            string
		selectedColumns []string
		wantHeaders     table.Row
		wantRowCount    int
	}{
		{
			name:            "full instance details",
			selectedColumns: []string{"Name", "State", "Instance Type", "Launch Template/Configuration", "Availability Zone", "Health"},
			wantHeaders:     table.Row{"Name", "State", "Instance Type", "Launch Template/Configuration", "Availability Zone", "Health"},
			wantRowCount:    1,
		},
	}

	for _, tc := range instanceTestCases {
		t.Run(tc.name, func(t *testing.T) {
			instanceTable := &AutoScalingInstanceTable{
				Instances:       instances,
				SelectedColumns: tc.selectedColumns,
			}

			// Test Headers
			headers := instanceTable.Headers()
			if len(headers) != len(tc.wantHeaders) {
				t.Errorf("Headers() returned %d columns, want %d", len(headers), len(tc.wantHeaders))
			}

			// Test Rows
			rows := instanceTable.Rows()
			if len(rows) != tc.wantRowCount {
				t.Errorf("Rows() returned %d rows, want %d", len(rows), tc.wantRowCount)
			}

			// Print the actual table output
			tw := table.NewWriter()
			tw.AppendHeader(headers)
			tw.AppendRows(rows)
			tw.SetStyle(instanceTable.TableStyle())
			t.Logf("\nTable Output:\n%s", tw.Render())
		})
	}
}
