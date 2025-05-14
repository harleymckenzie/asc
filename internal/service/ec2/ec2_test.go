package ec2

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
)

// MockEC2Client is a mock implementation of EC2ClientAPI for unit tests.
type MockEC2Client struct {
	mock.Mock
}

func (m *MockEC2Client) DescribeInstances(
	ctx context.Context,
	params *ec2.DescribeInstancesInput,
	optFns ...func(*ec2.Options),
) (*ec2.DescribeInstancesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeInstancesOutput), args.Error(1)
}

func (m *MockEC2Client) RebootInstances(
	ctx context.Context,
	params *ec2.RebootInstancesInput,
	optFns ...func(*ec2.Options),
) (*ec2.RebootInstancesOutput, error) {
	args := m.Called(ctx, params)
	return &ec2.RebootInstancesOutput{}, args.Error(1)
}

func (m *MockEC2Client) StartInstances(
	ctx context.Context,
	params *ec2.StartInstancesInput,
	optFns ...func(*ec2.Options),
) (*ec2.StartInstancesOutput, error) {
	args := m.Called(ctx, params)
	return &ec2.StartInstancesOutput{}, args.Error(1)
}

func (m *MockEC2Client) StopInstances(
	ctx context.Context,
	params *ec2.StopInstancesInput,
	optFns ...func(*ec2.Options),
) (*ec2.StopInstancesOutput, error) {
	args := m.Called(ctx, params)
	return &ec2.StopInstancesOutput{}, args.Error(1)
}

func (m *MockEC2Client) TerminateInstances(
	ctx context.Context,
	params *ec2.TerminateInstancesInput,
	optFns ...func(*ec2.Options),
) (*ec2.TerminateInstancesOutput, error) {
	args := m.Called(ctx, params)
	return &ec2.TerminateInstancesOutput{}, args.Error(1)
}

func (m *MockEC2Client) DescribeImages(
	ctx context.Context,
	params *ec2.DescribeImagesInput,
	optFns ...func(*ec2.Options),
) (*ec2.DescribeImagesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeImagesOutput), args.Error(1)
}

func (m *MockEC2Client) DescribeSecurityGroupRules(
	ctx context.Context,
	params *ec2.DescribeSecurityGroupRulesInput,
	optFns ...func(*ec2.Options),
) (*ec2.DescribeSecurityGroupRulesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeSecurityGroupRulesOutput), args.Error(1)
}

func (m *MockEC2Client) DescribeSecurityGroups(
	ctx context.Context,
	params *ec2.DescribeSecurityGroupsInput,
	optFns ...func(*ec2.Options),
) (*ec2.DescribeSecurityGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeSecurityGroupsOutput), args.Error(1)
}

func (m *MockEC2Client) DescribeSnapshots(
	ctx context.Context,
	params *ec2.DescribeSnapshotsInput,
	optFns ...func(*ec2.Options),
) (*ec2.DescribeSnapshotsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeSnapshotsOutput), args.Error(1)
}

func (m *MockEC2Client) DescribeVolumes(
	ctx context.Context,
	params *ec2.DescribeVolumesInput,
	optFns ...func(*ec2.Options),
) (*ec2.DescribeVolumesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeVolumesOutput), args.Error(1)
}

// Unit test for GetInstances
func TestGetInstances(t *testing.T) {
	mockClient := new(MockEC2Client)
	input := &ascTypes.GetInstancesInput{InstanceIDs: []string{"i-123"}}
	mockOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{Instances: []types.Instance{{InstanceId: &input.InstanceIDs[0]}}},
		},
	}
	mockClient.On("DescribeInstances", mock.Anything, mock.Anything).Return(mockOutput, nil)

	svc := &EC2Service{Client: mockClient}
	instances, err := svc.GetInstances(context.Background(), input)
	assert.NoError(t, err)
	assert.Len(t, instances, 1)
	assert.Equal(t, "i-123", *instances[0].InstanceId)
}

// Unit test for StartInstance
func TestStartInstance(t *testing.T) {
	mockClient := new(MockEC2Client)
	input := &ascTypes.StartInstanceInput{InstanceID: "i-123"}
	mockClient.On("StartInstances", mock.Anything, mock.Anything).
		Return(&ec2.StartInstancesOutput{}, nil)

	svc := &EC2Service{Client: mockClient}
	err := svc.StartInstance(context.Background(), input)
	assert.NoError(t, err)
}

// Unit test for StopInstance
func TestStopInstance(t *testing.T) {
	mockClient := new(MockEC2Client)
	input := &ascTypes.StopInstanceInput{InstanceID: "i-123", Force: true}
	mockClient.On("StopInstances", mock.Anything, mock.Anything).
		Return(&ec2.StopInstancesOutput{}, nil)

	svc := &EC2Service{Client: mockClient}
	err := svc.StopInstance(context.Background(), input)
	assert.NoError(t, err)
}

// Unit test for RestartInstance
func TestRestartInstance(t *testing.T) {
	mockClient := new(MockEC2Client)
	input := &ascTypes.RestartInstanceInput{InstanceID: "i-123"}
	mockClient.On("RebootInstances", mock.Anything, mock.Anything).
		Return(&ec2.RebootInstancesOutput{}, nil)

	svc := &EC2Service{Client: mockClient}
	err := svc.RestartInstance(context.Background(), input)
	assert.NoError(t, err)
}

// Unit test for TerminateInstance
func TestTerminateInstance(t *testing.T) {
	mockClient := new(MockEC2Client)
	input := &ascTypes.TerminateInstanceInput{InstanceID: "i-123"}
	mockClient.On("TerminateInstances", mock.Anything, mock.Anything).
		Return(&ec2.TerminateInstancesOutput{}, nil)

	svc := &EC2Service{Client: mockClient}
	err := svc.TerminateInstance(context.Background(), input)
	assert.NoError(t, err)
}

// Integration test for NewEC2Service (skipped unless EC2_INTEGRATION=1)
func TestNewEC2Service_Integration(t *testing.T) {
	if os.Getenv("INTEGRATION") != "1" {
		t.Skip("skipping integration test; set INTEGRATION=1 to run")
	}
	svc, err := NewEC2Service(context.Background(), "", "eu-west-1")
	assert.NoError(t, err)
	assert.NotNil(t, svc)
}
