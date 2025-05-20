package ec2

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)

// EC2ClientAPI is the interface for the EC2 client.
type EC2ClientAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	DescribeVolumes(ctx context.Context, params *ec2.DescribeVolumesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error)
	DescribeSecurityGroupRules(ctx context.Context, params *ec2.DescribeSecurityGroupRulesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupRulesOutput, error)
	DescribeSnapshots(ctx context.Context, params *ec2.DescribeSnapshotsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSnapshotsOutput, error)
	DescribeImages(ctx context.Context, params *ec2.DescribeImagesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error)
	DescribeSecurityGroups(ctx context.Context, params *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
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

// GetSecurityGroupRules gets the security group rules for the security group.
func (svc *EC2Service) GetSecurityGroupRules(ctx context.Context, input *ascTypes.GetSecurityGroupRulesInput) ([]types.SecurityGroupRule, error) {
	output, err := svc.Client.DescribeSecurityGroupRules(ctx, &ec2.DescribeSecurityGroupRulesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("group-id"),
				Values: []string{input.SecurityGroupID},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return output.SecurityGroupRules, nil
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

// GetVolumes fetches EC2 volumes and returns them directly.
func (svc *EC2Service) GetVolumes(ctx context.Context, input *ascTypes.GetVolumesInput) ([]types.Volume, error) {
	output, err := svc.Client.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{
		VolumeIds: input.VolumeIDs,
	})
	if err != nil {
		return nil, err
	}

	volumes := append([]types.Volume{}, output.Volumes...)
	return volumes, nil
}

// GetSnapshots fetches EC2 snapshots and returns them directly.
func (svc *EC2Service) GetSnapshots(ctx context.Context, input *ascTypes.GetSnapshotsInput) ([]types.Snapshot, error) {
	output, err := svc.Client.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
		SnapshotIds: input.SnapshotIDs,
		Filters:     input.Filters,
		OwnerIds:    input.OwnerIds,
	})
	if err != nil {
		return nil, err
	}

	snapshots := append([]types.Snapshot{}, output.Snapshots...)
	return snapshots, nil
}

// GetImages fetches EC2 images and returns them directly.
func (svc *EC2Service) GetImages(ctx context.Context, input *ascTypes.GetImagesInput) ([]types.Image, error) {
	output, err := svc.Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		ImageIds: input.ImageIDs,
	})
	if err != nil {
		return nil, err
	}

	images := append([]types.Image{}, output.Images...)
	return images, nil
}

// GetSecurityGroups fetches EC2 security groups and returns them directly.
func (svc *EC2Service) GetSecurityGroups(ctx context.Context, input *ascTypes.GetSecurityGroupsInput) ([]types.SecurityGroup, error) {
	output, err := svc.Client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupIds: input.GroupIDs,
	})
	if err != nil {
		return nil, err
	}

	groups := append([]types.SecurityGroup{}, output.SecurityGroups...)
	return groups, nil
}

// GetImagesWithFilters fetches EC2 images with custom filters and owners.
func (svc *EC2Service) GetImagesWithFilters(ctx context.Context, input *ascTypes.GetImagesInput, filters []types.Filter, owners []string) ([]types.Image, error) {
	output, err := svc.Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		ImageIds: input.ImageIDs,
		Filters:  filters,
		Owners:   owners,
	})
	if err != nil {
		return nil, err
	}
	images := append([]types.Image{}, output.Images...)
	return images, nil
}

// FilterSecurityGroupRules will filter the rules by inbound or outbound
func FilterSecurityGroupRules(rules []types.SecurityGroupRule, egress bool) []types.SecurityGroupRule {
	filteredRules := []types.SecurityGroupRule{}
	for _, rule := range rules {
		if egress && *rule.IsEgress {
			filteredRules = append(filteredRules, rule)
		} else if !egress && !*rule.IsEgress {
			filteredRules = append(filteredRules, rule)
		}
	}
	return filteredRules
}
