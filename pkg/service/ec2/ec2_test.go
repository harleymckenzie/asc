package ec2

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/jedib0t/go-pretty/v6/table"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
)

type mockEC2Client struct {
	describeInstancesOutput *ec2.DescribeInstancesOutput
	err                     error
}

func (m *mockEC2Client) DescribeInstances(
	_ context.Context,
	params *ec2.DescribeInstancesInput,
	_ ...func(*ec2.Options),
) (*ec2.DescribeInstancesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.describeInstancesOutput, nil
}

func (m *mockEC2Client) StartInstances(
	_ context.Context,
	_ *ec2.StartInstancesInput,
	_ ...func(*ec2.Options),
) (*ec2.StartInstancesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &ec2.StartInstancesOutput{}, nil
}

func (m *mockEC2Client) StopInstances(
	_ context.Context,
	_ *ec2.StopInstancesInput,
	_ ...func(*ec2.Options),
) (*ec2.StopInstancesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &ec2.StopInstancesOutput{}, nil
}

func (m *mockEC2Client) TerminateInstances(
	_ context.Context,
	_ *ec2.TerminateInstancesInput,
	_ ...func(*ec2.Options),
) (*ec2.TerminateInstancesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &ec2.TerminateInstancesOutput{}, nil
}

func (m *mockEC2Client) RebootInstances(
	_ context.Context,
	_ *ec2.RebootInstancesInput,
	_ ...func(*ec2.Options),
) (*ec2.RebootInstancesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &ec2.RebootInstancesOutput{}, nil
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
					InstanceId:   aws.String("i-1234567890abcdef0"),
					InstanceType: types.InstanceType("t3.micro"),
					State: &types.InstanceState{
						Name: types.InstanceStateName("running"),
					},
					Tags: []types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-instance"),
						},
					},
				},
				{
					InstanceId:   aws.String("i-1234567890abcdef1"),
					InstanceType: types.InstanceType("t3.small"),
					State: &types.InstanceState{
						Name: types.InstanceStateName("stopped"),
					},
					Tags: []types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-instance-2"),
						},
					},
				},
			},
			err:     nil,
			wantErr: false,
		},
		{
			name:      "empty response",
			instances: []types.Instance{},
			err:       nil,
			wantErr:   false,
		},
		{
			name:      "api error",
			instances: nil,
			err:       errors.New("Invalid instance ID"),
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockEC2Client{
				describeInstancesOutput: &ec2.DescribeInstancesOutput{
					Reservations: []types.Reservation{
						{
							Instances: tc.instances,
						},
					},
				},
				err: tc.err,
			}

			svc := &EC2Service{
				Client: mockClient,
			}

			var instanceIDs []string
			if len(tc.instances) > 0 {
				instanceIDs = []string{*tc.instances[0].InstanceId}
			}

			instances, err := svc.GetInstances(context.Background(), &ascTypes.GetInstancesInput{
				InstanceIDs: instanceIDs,
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
	instances := []types.Instance{
		{
			InstanceId:   aws.String("i-1234567890abcdef0"),
			InstanceType: types.InstanceType("t3.micro"),
			State: &types.InstanceState{
				Name: types.InstanceStateName("running"),
			},
			Tags: []types.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("test-instance"),
				},
			},
		},
		{
			InstanceId:   aws.String("i-1234567890abcdef1"),
			InstanceType: types.InstanceType("t3.small"),
			State: &types.InstanceState{
				Name: types.InstanceStateName("stopped"),
			},
			Tags: []types.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("test-instance-2"),
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
			name:            "basic instance details",
			selectedColumns: []string{"Name", "Instance ID", "State", "Instance Type"},
			wantHeaders:     table.Row{"Name", "Instance ID", "State", "Instance Type"},
			wantRowCount:    2,
		},
		{
			name:            "minimal columns",
			selectedColumns: []string{"Name", "State"},
			wantHeaders:     table.Row{"Name", "State"},
			wantRowCount:    2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Table := &EC2Table{
				Instances:       instances,
				SelectedColumns: tc.selectedColumns,
			}

			// Test Headers
			headers := ec2Table.Headers()
			if len(headers) != len(tc.wantHeaders) {
				t.Errorf("Headers() returned %d columns, want %d", len(headers), len(tc.wantHeaders))
			}
			for i, h := range headers {
				if h != tc.wantHeaders[i] {
					t.Errorf("Headers()[%d] = %v, want %v", i, h, tc.wantHeaders[i])
				}
			}

			// Test Rows
			rows := ec2Table.Rows()
			if len(rows) != tc.wantRowCount {
				t.Errorf("Rows() returned %d rows, want %d", len(rows), tc.wantRowCount)
			}

			// Print the actual table output for visual inspection
			tw := table.NewWriter()
			tw.AppendHeader(headers)
			tw.AppendRows(rows)
			tw.SetStyle(ec2Table.TableStyle())
			t.Logf("\nTable Output:\n%s", tw.Render())
		})
	}
}
