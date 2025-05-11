package ec2

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
	"github.com/harleymckenzie/asc/pkg/shared/awsutil"
)

// EC2ClientAPI is the interface for the EC2 client.
type EC2ClientAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	RebootInstances(ctx context.Context, params *ec2.RebootInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RebootInstancesOutput, error)
	StartInstances(ctx context.Context, params *ec2.StartInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error)
	StopInstances(ctx context.Context, params *ec2.StopInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StopInstancesOutput, error)
	TerminateInstances(ctx context.Context, params *ec2.TerminateInstancesInput, optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}

// EC2Service is a struct that holds the EC2 client.
type EC2Service struct {
	Client EC2ClientAPI
}

//
// Service functions
//

// NewEC2Service creates a new EC2 service.
func NewEC2Service(ctx context.Context, profile string, region string) (*EC2Service, error) {
	cfg, err := awsutil.LoadDefaultConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}
	client := ec2.NewFromConfig(cfg.Config)

	return &EC2Service{Client: client}, nil
}

// GetInstances fetches EC2 instances and returns them directly.
func (svc *EC2Service) GetInstances(ctx context.Context, input *ascTypes.GetInstancesInput) ([]types.Instance, error) {
	output, err := svc.Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: input.InstanceIDs,
	})
	if err != nil {
		return nil, err
	}

	var instances []types.Instance
	for _, reservation := range output.Reservations {
		instances = append(instances, reservation.Instances...)
	}
	return instances, nil
}

// getInstanceName gets the name of the instance from the tags.
func getInstanceName(instance types.Instance) string {
	// Get instance name from tags
	name := "-" // Use as default name if "Name" tag doesn't exist
	for _, tag := range instance.Tags {
		if aws.ToString(tag.Key) == "Name" {
			name = aws.ToString(tag.Value)
			break
		}
	}
	return name
}

// getSecurityGroups gets the security groups for the instance.
func getSecurityGroups(securityGroups []types.GroupIdentifier) string {
	securityGroupsList := []string{}
	for _, group := range securityGroups {
		securityGroupsList = append(securityGroupsList, aws.ToString(group.GroupId))
	}
	return strings.Join(securityGroupsList, "\n")
}

// RestartInstance restarts an instance.
func (svc *EC2Service) RestartInstance(ctx context.Context, input *ascTypes.RestartInstanceInput) error {
	_, err := svc.Client.RebootInstances(ctx, &ec2.RebootInstancesInput{
		InstanceIds: []string{input.InstanceID},
	})
	if err != nil {
		return err
	}
	return nil
}

// StartInstance starts an instance.
func (svc *EC2Service) StartInstance(ctx context.Context, input *ascTypes.StartInstanceInput) error {
	_, err := svc.Client.StartInstances(ctx, &ec2.StartInstancesInput{
		InstanceIds: []string{input.InstanceID},
	})
	if err != nil {
		return err
	}
	return nil
}

// StopInstance stops an instance.
func (svc *EC2Service) StopInstance(ctx context.Context, input *ascTypes.StopInstanceInput) error {
	_, err := svc.Client.StopInstances(ctx, &ec2.StopInstancesInput{
		InstanceIds: []string{input.InstanceID},
		Force:       &input.Force,
	})
	if err != nil {
		return err
	}
	return nil
}

// TerminateInstance terminates an instance.
func (svc *EC2Service) TerminateInstance(ctx context.Context, input *ascTypes.TerminateInstanceInput) error {
	_, err := svc.Client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{input.InstanceID},
	})
	if err != nil {
		return err
	}
	return nil
}
