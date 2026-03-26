package ec2

import (
	"context"
	"fmt"
	"strings"

	ec2sdk "github.com/aws/aws-sdk-go-v2/service/ec2"
)

var ec2TerminalStates = map[string]bool{
	"running":    true,
	"stopped":    true,
	"terminated": true,
}

var volumeTerminalStates = map[string]bool{
	"available": true,
	"in-use":    true,
	"deleted":   true,
	"error":     true,
}

var snapshotTerminalStates = map[string]bool{
	"completed": true,
	"error":     true,
}

var imageTerminalStates = map[string]bool{
	"available":    true,
	"failed":       true,
	"deregistered": true,
}

// IsTerminalInstanceState returns true if the EC2 instance state is a stable,
// non-processing state.
func IsTerminalInstanceState(status string) bool {
	return ec2TerminalStates[strings.ToLower(status)]
}

// IsTerminalVolumeState returns true if the EBS volume state is stable.
func IsTerminalVolumeState(status string) bool {
	return volumeTerminalStates[strings.ToLower(status)]
}

// IsTerminalSnapshotState returns true if the EBS snapshot state is stable.
func IsTerminalSnapshotState(status string) bool {
	return snapshotTerminalStates[strings.ToLower(status)]
}

// IsTerminalImageState returns true if the AMI state is stable.
func IsTerminalImageState(status string) bool {
	return imageTerminalStates[strings.ToLower(status)]
}

// GetInstanceStatus returns the current state of an EC2 instance.
func (svc *EC2Service) GetInstanceStatus(ctx context.Context, instanceID string) (string, error) {
	output, err := svc.Client.DescribeInstances(ctx, &ec2sdk.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return "", fmt.Errorf("describe instance: %w", err)
	}
	if len(output.Reservations) == 0 || len(output.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("instance %s not found", instanceID)
	}
	return string(output.Reservations[0].Instances[0].State.Name), nil
}

// GetVolumeStatus returns the current state of an EBS volume.
func (svc *EC2Service) GetVolumeStatus(ctx context.Context, volumeID string) (string, error) {
	output, err := svc.Client.DescribeVolumes(ctx, &ec2sdk.DescribeVolumesInput{
		VolumeIds: []string{volumeID},
	})
	if err != nil {
		return "", fmt.Errorf("describe volume: %w", err)
	}
	if len(output.Volumes) == 0 {
		return "", fmt.Errorf("volume %s not found", volumeID)
	}
	return string(output.Volumes[0].State), nil
}

// GetSnapshotStatus returns the current state of an EBS snapshot.
func (svc *EC2Service) GetSnapshotStatus(ctx context.Context, snapshotID string) (string, error) {
	output, err := svc.Client.DescribeSnapshots(ctx, &ec2sdk.DescribeSnapshotsInput{
		SnapshotIds: []string{snapshotID},
	})
	if err != nil {
		return "", fmt.Errorf("describe snapshot: %w", err)
	}
	if len(output.Snapshots) == 0 {
		return "", fmt.Errorf("snapshot %s not found", snapshotID)
	}
	return string(output.Snapshots[0].State), nil
}

// GetImageStatus returns the current state of an AMI.
func (svc *EC2Service) GetImageStatus(ctx context.Context, imageID string) (string, error) {
	output, err := svc.Client.DescribeImages(ctx, &ec2sdk.DescribeImagesInput{
		ImageIds: []string{imageID},
	})
	if err != nil {
		return "", fmt.Errorf("describe image: %w", err)
	}
	if len(output.Images) == 0 {
		return "", fmt.Errorf("image %s not found", imageID)
	}
	return string(output.Images[0].State), nil
}
