package ec2

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
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
				ctx:    context.Background(),
			}

			err := svc.ListInstances(context.Background())
			if (err != nil) != tc.wantErr {
				t.Errorf("ListInstances() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
