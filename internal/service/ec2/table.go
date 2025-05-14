// Package ec2 provides functions for interacting with EC2 resources.
package ec2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// Attribute is a struct that defines a field in a detailed table.
type Attribute struct {
	GetValue func(*types.Instance) string
}

type VolumeAttribute struct {
	GetValue func(*types.Volume) string
}

type SnapshotAttribute struct {
	GetValue func(*types.Snapshot) string
}

type ImageAttribute struct {
	GetValue func(*types.Image) string
}

type SecurityGroupAttribute struct {
	// GetValue returns the string value for a field from a SecurityGroup
	GetValue func(*types.SecurityGroup) string
}

type SecurityGroupRuleAttribute struct {
	GetValue func(*types.SecurityGroupRule) string
}

// GetAttributeValue is a function that returns the value of a field in a detailed table.
func GetAttributeValue(fieldID string, instance any) (string, error) {
	inst, ok := instance.(types.Instance)
	if !ok {
		return "", fmt.Errorf("instance is not a types.Instance")
	}
	attr, exists := availableAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&inst), nil
}

// availableAttributes is a function that returns a map of attributes for a detailed table.
func availableAttributes() map[string]Attribute {
	return map[string]Attribute{
		"Name": {
			GetValue: func(i *types.Instance) string {
				return getInstanceName(*i)
			},
		},
		"Instance ID": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.InstanceId)
			},
		},
		"State": {
			GetValue: func(i *types.Instance) string {
				return format.Status(string(i.State.Name))
			},
		},
		"Instance Type": {
			GetValue: func(i *types.Instance) string {
				return string(i.InstanceType)
			},
		},
		"AMI ID": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.ImageId)
			},
		},
		"AMI Name": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.ImageId)
			},
		},
		"Public IP": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.PublicIpAddress)
			},
		},
		"Private IP": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.PrivateIpAddress)
			},
		},
		"Launch Time": {
			GetValue: func(i *types.Instance) string {
				return i.LaunchTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"Subnet ID": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.SubnetId)
			},
		},
		"Security Group(s)": {
			GetValue: func(i *types.Instance) string {
				return getSecurityGroups(i.SecurityGroups)
			},
		},
		"Key Name": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.KeyName)
			},
		},
		"VPC ID": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.VpcId)
			},
		},
		"Placement Group": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.Placement.GroupName)
			},
		},
		"Availability Zone": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.Placement.AvailabilityZone)
			},
		},
		"Root Device Type": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString((*string)(&i.RootDeviceType))
			},
		},
		"Root Device Name": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.RootDeviceName)
			},
		},
		"Virtualization Type": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString((*string)(&i.VirtualizationType))
			},
		},
		"vCPUs": {
			GetValue: func(i *types.Instance) string {
				return strconv.Itoa(int(*i.CpuOptions.CoreCount))
			},
		},
	}
}

// GetImageAttributeValue is a function that returns the value of a field in a detailed table.
func GetImageAttributeValue(fieldID string, image any) (string, error) {
	img, ok := image.(types.Image)
	if !ok {
		return "", fmt.Errorf("image is not a types.Image")
	}
	attr, exists := imageAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&img), nil
}

// imageAttributes is a function that returns a map of attributes for a detailed table.
func imageAttributes() map[string]ImageAttribute {
	return map[string]ImageAttribute{
		"Allowed Image": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.ImageAllowed == nil {
					return ""
				}
				return strconv.FormatBool(*i.ImageAllowed)
			},
		},
		"AMI ID": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.ImageId == nil {
					return ""
				}
				return aws.ToString(i.ImageId)
			},
		},
		"AMI Name": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.Name == nil {
					return ""
				}
				return aws.ToString(i.Name)
			},
		},
		"Architecture": {
			GetValue: func(i *types.Image) string {
				if i == nil {
					return ""
				}
				return string(i.Architecture)
			},
		},
		"Block Devices": {
			GetValue: func(i *types.Image) string {
				if i == nil {
					return ""
				}
				return getBlockDevices(i.BlockDeviceMappings)
			},
		},
		"Boot Mode": {
			GetValue: func(i *types.Image) string {
				if i == nil {
					return ""
				}
				return string(i.BootMode)
			},
		},
		"Creation Date": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.CreationDate == nil {
					return ""
				}
				return aws.ToString(i.CreationDate)
			},
		},
		"Deprecation Time": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.DeprecationTime == nil {
					return ""
				}
				return aws.ToString(i.DeprecationTime)
			},
		},
		"Deregistration Protection": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.DeregistrationProtection == nil {
					return ""
				}
				return aws.ToString(i.DeregistrationProtection)
			},
		},
		"Description": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.Description == nil {
					return ""
				}
				return aws.ToString(i.Description)
			},
		},
		"Image Type": {
			GetValue: func(i *types.Image) string {
				if i == nil {
					return ""
				}
				return string(i.ImageType)
			},
		},
		"Kernel ID": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.KernelId == nil {
					return ""
				}
				return aws.ToString(i.KernelId)
			},
		},
		"Owner": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.OwnerId == nil {
					return ""
				}
				return aws.ToString(i.OwnerId)
			},
		},
		"Platform": {
			GetValue: func(i *types.Image) string {
				if i == nil {
					return ""
				}
				return string(*i.PlatformDetails)
			},
		},
		"Product Codes": {
			GetValue: func(i *types.Image) string {
				if i == nil {
					return ""
				}
				return getProductCodes(i.ProductCodes)
			},
		},
		"RAM Disk ID": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.RamdiskId == nil {
					return ""
				}
				return aws.ToString(i.RamdiskId)
			},
		},
		"Root Device Name": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.RootDeviceName == nil {
					return ""
				}
				return aws.ToString(i.RootDeviceName)
			},
		},
		"Root Device Type": {
			GetValue: func(i *types.Image) string {
				if i == nil {
					return ""
				}
				return string(i.RootDeviceType)
			},
		},
		"Source": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.ImageLocation == nil {
					return ""
				}
				return aws.ToString(i.ImageLocation)
			},
		},
		"Source AMI ID": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.SourceImageId == nil {
					return ""
				}
				return aws.ToString(i.SourceImageId)
			},
		},
		"Source AMI Region": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.SourceImageRegion == nil {
					return ""
				}
				return aws.ToString(i.SourceImageRegion)
			},
		},
		"State Reason": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.StateReason == nil || i.StateReason.Message == nil {
					return ""
				}
				return aws.ToString(i.StateReason.Message)
			},
		},
		"Status": {
			GetValue: func(i *types.Image) string {
				if i == nil {
					return ""
				}
				return format.Status(string(i.State))
			},
		},
		"Usage Operation": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.UsageOperation == nil {
					return ""
				}
				return aws.ToString(i.UsageOperation)
			},
		},
		"Virtualization": {
			GetValue: func(i *types.Image) string {
				if i == nil {
					return ""
				}
				return string(i.VirtualizationType)
			},
		},
		"Visibility": {
			GetValue: func(i *types.Image) string {
				if i == nil || i.Public == nil {
					return ""
				}
				if *i.Public {
					return "Public"
				}
				return "Private"
			},
		},
	}
}

// GetSecurityGroupAttributeValue returns the value of a field for a SecurityGroup.
func GetSecurityGroupAttributeValue(fieldID string, group any) (string, error) {
	sg, ok := group.(types.SecurityGroup)
	if !ok {
		return "", fmt.Errorf("group is not a types.SecurityGroup")
	}
	attr, exists := securityGroupAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&sg), nil
}

// securityGroupAttributes returns a map of field IDs to SecurityGroupAttribute.
func securityGroupAttributes() map[string]SecurityGroupAttribute {
	return map[string]SecurityGroupAttribute{
		"Group ID": {
			GetValue: func(sg *types.SecurityGroup) string {
				return aws.ToString(sg.GroupId)
			},
		},
		"Group Name": {
			GetValue: func(sg *types.SecurityGroup) string {
				return aws.ToString(sg.GroupName)
			},
		},
		"Description": {
			GetValue: func(sg *types.SecurityGroup) string {
				return aws.ToString(sg.Description)
			},
		},
		"VPC ID": {
			GetValue: func(sg *types.SecurityGroup) string {
				return aws.ToString(sg.VpcId)
			},
		},
		"Owner ID": {
			GetValue: func(sg *types.SecurityGroup) string {
				return aws.ToString(sg.OwnerId)
			},
		},
		"Ingress Count": {
			GetValue: func(sg *types.SecurityGroup) string {
				return fmt.Sprintf("%d entries", len(sg.IpPermissions))
			},
		},
		"Egress Count": {
			GetValue: func(sg *types.SecurityGroup) string {
				return fmt.Sprintf("%d entries", len(sg.IpPermissionsEgress))
			},
		},
		"Tag Count": {
			GetValue: func(sg *types.SecurityGroup) string {
				return strconv.Itoa(len(sg.Tags))
			},
		},
	}
}

// GetSecurityGroupIpPermissionAttributeValue returns the value of a field for a SecurityGroupIpPermission.
func GetSecurityGroupRuleAttributeValue(fieldID string, rule any) (string, error) {
	r, ok := rule.(types.SecurityGroupRule)
	if !ok {
		return "", fmt.Errorf("rule is not a types.SecurityGroupRule")
	}
	attr, exists := securityGroupRuleAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&r), nil
}

func securityGroupRuleAttributes() map[string]SecurityGroupRuleAttribute {
	return map[string]SecurityGroupRuleAttribute{
		// Composite/Custom attributes for UI
		"Rule ID": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.SecurityGroupRuleId == nil {
					return ""
				}
				return *r.SecurityGroupRuleId
			},
		},
		"IP Version": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.CidrIpv4 != nil && *r.CidrIpv4 != "" {
					return "IPv4"
				}
				if r.CidrIpv6 != nil && *r.CidrIpv6 != "" {
					return "IPv6"
				}
				return "-"
			},
		},
		"Type": {
			GetValue: func(r *types.SecurityGroupRule) string {
				return getSecurityGroupRuleType(*r)
			},
		},
		"Protocol": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.IpProtocol == nil {
					return ""
				}
				if *r.IpProtocol == "-1" {
					return "All"
				}
				return strings.ToUpper(*r.IpProtocol)
			},
		},
		"Port Range": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if *r.FromPort == -1 {
					return "All"
				}
				if r.FromPort != nil && r.ToPort != nil {
					if *r.FromPort == *r.ToPort {
						return strconv.Itoa(int(*r.FromPort))
					}
					return fmt.Sprintf("%d-%d", *r.FromPort, *r.ToPort)
				}
				return ""
			},
		},
		"Source": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.IsEgress != nil && *r.IsEgress {
					return ""
				}
				if r.CidrIpv4 != nil && *r.CidrIpv4 != "" {
					return *r.CidrIpv4
				}
				if r.CidrIpv6 != nil && *r.CidrIpv6 != "" {
					return *r.CidrIpv6
				}
				if r.ReferencedGroupInfo != nil && r.ReferencedGroupInfo.GroupId != nil {
					return *r.ReferencedGroupInfo.GroupId
				}
				return ""
			},
		},
		"Destination": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.IsEgress != nil && *r.IsEgress {
					if r.CidrIpv4 != nil && *r.CidrIpv4 != "" {
						return *r.CidrIpv4
					}
					if r.CidrIpv6 != nil && *r.CidrIpv6 != "" {
						return *r.CidrIpv6
					}
					if r.ReferencedGroupInfo != nil && r.ReferencedGroupInfo.GroupId != nil {
						return *r.ReferencedGroupInfo.GroupId
					}
				}
				return ""
			},
		},
		"Description": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.Description == nil {
					return ""
				}
				return *r.Description
			},
		},

		// Retain all original attributes for flexibility
		"CidrIpv4": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.CidrIpv4 == nil {
					return ""
				}
				return *r.CidrIpv4
			},
		},
		"CidrIpv6": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.CidrIpv6 == nil {
					return ""
				}
				return *r.CidrIpv6
			},
		},
		"FromPort": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.FromPort == nil {
					return ""
				}
				return strconv.Itoa(int(*r.FromPort))
			},
		},
		"GroupId": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.GroupId == nil {
					return ""
				}
				return *r.GroupId
			},
		},
		"GroupOwnerId": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.GroupOwnerId == nil {
					return ""
				}
				return *r.GroupOwnerId
			},
		},
		"IsEgress": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.IsEgress == nil {
					return ""
				}
				return strconv.FormatBool(*r.IsEgress)
			},
		},
		"PrefixListId": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.PrefixListId == nil {
					return ""
				}
				return *r.PrefixListId
			},
		},
		"ReferencedGroupInfo": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.ReferencedGroupInfo == nil || r.ReferencedGroupInfo.GroupId == nil {
					return ""
				}
				return *r.ReferencedGroupInfo.GroupId
			},
		},
		"SecurityGroupRuleArn": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.SecurityGroupRuleArn == nil {
					return ""
				}
				return *r.SecurityGroupRuleArn
			},
		},
		"SecurityGroupRuleId": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.SecurityGroupRuleId == nil {
					return ""
				}
				return *r.SecurityGroupRuleId
			},
		},
		"TagCount": {
			GetValue: func(r *types.SecurityGroupRule) string {
				return strconv.Itoa(len(r.Tags))
			},
		},
		"ToPort": {
			GetValue: func(r *types.SecurityGroupRule) string {
				if r.ToPort == nil {
					return ""
				}
				return strconv.Itoa(int(*r.ToPort))
			},
		},
	}
}

// GetVolumeAttributeValue is a function that returns the value of a field in a detailed table.
func GetVolumeAttributeValue(fieldID string, volume any) (string, error) {
	vol, ok := volume.(types.Volume)
	if !ok {
		return "", fmt.Errorf("volume is not a types.Volume")
	}
	attr, exists := volumeAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&vol), nil
}

// GetSnapshotAttributeValue is a function that returns the value of a field in a detailed table.
func GetSnapshotAttributeValue(fieldID string, snapshot any) (string, error) {
	snap, ok := snapshot.(types.Snapshot)
	if !ok {
		return "", fmt.Errorf("snapshot is not a types.Snapshot")
	}
	attr, exists := snapshotAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&snap), nil
}

// snapshotAttributes is a function that returns a map of attributes for a detailed table.
func snapshotAttributes() map[string]SnapshotAttribute {
	return map[string]SnapshotAttribute{
		"Description": {
			GetValue: func(s *types.Snapshot) string {
				if s.Description == nil {
					return ""
				}
				return *s.Description
			},
		},
		"Details": {
			GetValue: func(s *types.Snapshot) string {
				if s.SnapshotId == nil {
					return ""
				}
				return aws.ToString(s.SnapshotId)
			},
		},
		"Encryption": {
			GetValue: func(s *types.Snapshot) string {
				if s.Encrypted == nil {
					return ""
				}
				return fmt.Sprintf("%t", *s.Encrypted)
			},
		},
		"KMS Key ID": {
			GetValue: func(s *types.Snapshot) string {
				if s.KmsKeyId == nil {
					return ""
				}
				return aws.ToString(s.KmsKeyId)
			},
		},
		"Owner Alias": {
			GetValue: func(s *types.Snapshot) string {
				if s.OwnerAlias == nil {
					return ""
				}
				return aws.ToString(s.OwnerAlias)
			},
		},
		"Owner ID": {
			GetValue: func(s *types.Snapshot) string {
				if s.OwnerId == nil {
					return ""
				}
				return aws.ToString(s.OwnerId)
			},
		},
		"Progress": {
			GetValue: func(s *types.Snapshot) string {
				if s.Progress == nil {
					return ""
				}
				return format.Status(aws.ToString(s.Progress))
			},
		},
		"Restore Expiry Time": {
			GetValue: func(s *types.Snapshot) string {
				if s.RestoreExpiryTime == nil {
					return ""
				}
				return s.RestoreExpiryTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"Snapshot ID": {
			GetValue: func(s *types.Snapshot) string {
				if s.SnapshotId == nil {
					return ""
				}
				return aws.ToString(s.SnapshotId)
			},
		},
		"Started": {
			GetValue: func(s *types.Snapshot) string {
				if s.StartTime == nil {
					return ""
				}
				return s.StartTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"State": {
			GetValue: func(s *types.Snapshot) string {
				return format.Status(string(s.State))
			},
		},
		"State Message": {
			GetValue: func(s *types.Snapshot) string {
				if s.StateMessage == nil {
					return ""
				}
				return aws.ToString(s.StateMessage)
			},
		},
		"Storage Tier": {
			GetValue: func(s *types.Snapshot) string {
				return string(s.StorageTier)
			},
		},
		"Tier": {
			GetValue: func(s *types.Snapshot) string {
				return string(s.StorageTier)
			},
		},
		"Volume ID": {
			GetValue: func(s *types.Snapshot) string {
				if s.VolumeId == nil {
					return ""
				}
				return aws.ToString(s.VolumeId)
			},
		},
		"Volume Size": {
			GetValue: func(s *types.Snapshot) string {
				if s.VolumeSize == nil {
					return ""
				}
				return fmt.Sprintf("%d GiB", *s.VolumeSize)
			},
		},
	}
}

// volumeAttributes is a function that returns a map of attributes for a detailed table.
func volumeAttributes() map[string]VolumeAttribute {
	return map[string]VolumeAttribute{
		"Volume ID": {
			GetValue: func(v *types.Volume) string {
				if v.VolumeId == nil {
					return ""
				}
				return aws.ToString(v.VolumeId)
			},
		},
		"Type": {
			GetValue: func(v *types.Volume) string {
				return string(v.VolumeType)
			},
		},
		"Size": {
			GetValue: func(v *types.Volume) string {
				if v.Size == nil {
					return ""
				}
				return fmt.Sprintf("%d GiB", *v.Size)
			},
		},
		"State": {
			GetValue: func(v *types.Volume) string {
				return format.Status(string(v.State))
			},
		},
		"IOPS": {
			GetValue: func(v *types.Volume) string {
				if v.Iops == nil {
					return ""
				}
				return strconv.Itoa(int(*v.Iops))
			},
		},
		"Throughput": {
			GetValue: func(v *types.Volume) string {
				if v.Throughput == nil {
					return ""
				}
				return strconv.Itoa(int(*v.Throughput))
			},
		},
		"Fast Snapshot Restored": {
			GetValue: func(v *types.Volume) string {
				if v.FastRestored == nil {
					return ""
				}
				return strconv.FormatBool(*v.FastRestored)
			},
		},
		"Availability Zone": {
			GetValue: func(v *types.Volume) string {
				return aws.ToString(v.AvailabilityZone)
			},
		},
		"Created": {
			GetValue: func(v *types.Volume) string {
				return v.CreateTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"Multi-Attach Enabled": {
			GetValue: func(v *types.Volume) string {
				return strconv.FormatBool(*v.MultiAttachEnabled)
			},
		},
		"Snapshot ID": {
			GetValue: func(v *types.Volume) string {
				return aws.ToString(v.SnapshotId)
			},
		},
		"Associated Resource": {
			GetValue: func(v *types.Volume) string {
				return getAssociatedResource(v.Attachments)
			},
		},
		"Attach Time": {
			GetValue: func(v *types.Volume) string {
				return v.Attachments[0].AttachTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"Delete on Termination": {
			GetValue: func(v *types.Volume) string {
				return strconv.FormatBool(*v.Attachments[0].DeleteOnTermination)
			},
		},
		"Device": {
			GetValue: func(v *types.Volume) string {
				return aws.ToString(v.Attachments[0].Device)
			},
		},
		"Instance ID": {
			GetValue: func(v *types.Volume) string {
				return aws.ToString(v.Attachments[0].InstanceId)
			},
		},
		"Attachment State": {
			GetValue: func(v *types.Volume) string {
				return string(v.Attachments[0].State)
			},
		},
		"Encryption": {
			GetValue: func(v *types.Volume) string {
				if v.Encrypted == nil {
					return ""
				}
				if *v.Encrypted {
					return "Encrypted"
				}
				return "Not Encrypted"
			},
		},
		"KMS Key ID": {
			GetValue: func(v *types.Volume) string {
				return aws.ToString(v.KmsKeyId)
			},
		},
	}
}

// Helper functions
// #TODO: Move to a better location

// getAssociatedResource is a function that returns the associated resource for a volume.
func getAssociatedResource(attachments []types.VolumeAttachment) string {
	if len(attachments) == 0 {
		return ""
	}
	return aws.ToString(attachments[0].InstanceId)
}

// getBlockDevices is a function that returns the block devices for an image.
func getBlockDevices(blockDevices []types.BlockDeviceMapping) string {
	if len(blockDevices) == 0 {
		return ""
	}
	devices := []string{}
	for _, bd := range blockDevices {
		deviceName := aws.ToString(bd.DeviceName)
		if bd.Ebs != nil {
			snapshotId := ""
			if bd.Ebs.SnapshotId != nil {
				snapshotId = *bd.Ebs.SnapshotId
			}
			size := ""
			if bd.Ebs.VolumeSize != nil {
				size = fmt.Sprintf("%d", *bd.Ebs.VolumeSize)
			}
			deleteOnTermination := ""
			if bd.Ebs.DeleteOnTermination != nil {
				deleteOnTermination = fmt.Sprintf("%t", *bd.Ebs.DeleteOnTermination)
			}
			volumeType := string(bd.Ebs.VolumeType)
			encrypted := ""
			if bd.Ebs.Encrypted != nil && *bd.Ebs.Encrypted {
				encrypted = "encrypted"
			}
			devices = append(
				devices,
				fmt.Sprintf(
					"%s=%s:%s:%s:%s:%s",
					deviceName,
					snapshotId,
					size,
					deleteOnTermination,
					volumeType,
					encrypted,
				),
			)
		} else if bd.VirtualName != nil {
			devices = append(devices, fmt.Sprintf("%s=%s", deviceName, *bd.VirtualName))
		}
	}
	return strings.Join(devices, "\n")
}

// getProductCodes is a function that returns the product codes for an image.
func getProductCodes(productCodes []types.ProductCode) string {
	if len(productCodes) == 0 {
		return ""
	}
	// Each product code is made up of ProductCodeId and ProductCodeType
	code := aws.ToString(productCodes[0].ProductCodeId)
	codeType := string(productCodes[0].ProductCodeType)
	return fmt.Sprintf("%s (%s)", code, codeType)
}

// getSecurityGroupRuleType is a function that returns the type of a security group rule.
// (eg,. if port range is 443-443, return "HTTPS", if port range is 22-22, return "SSH", if port range is 80-80, return "HTTP")
func getSecurityGroupRuleType(rule types.SecurityGroupRule) string {
	if rule.FromPort != nil && rule.ToPort != nil {
		if *rule.FromPort == *rule.ToPort {
			switch *rule.FromPort {
			case -1:
				return "All traffic"
			case 22:
				return "SSH"
			case 25:
				return "SMTP"
			case 53:
				return "DNS"
			case 80:
				return "HTTP"
			case 110:
				return "POP3"
			case 143:
				return "IMAP"
			case 389:
				return "LDAP"
			case 443:
				return "HTTPS"
			case 445:
				return "SMB"
			case 465:
				return "SMTPS"
			case 993:
				return "IMAPS"
			case 995:
				return "POP3S"
			case 1433:
				return "MSSQL"
			case 2049:
				return "NFS"
			case 3306:
				return "MySQL/Aurora"
			case 3389:
				return "RDP"
			case 5439:
				return "Redshift"
			case 5432:
				return "PostgreSQL"
			case 1521:
				return "Oracle RDS"
			case 5985:
				return "WinRM-HTTP"
			case 5986:
				return "WinRM-HTTPS"
			case 20049:
				return "Elastic Graphics"
			case 9042:
				return "CQLSH / Cassandra"
			default:
				// Handle custom protocol types based on IpProtocol
				if rule.IpProtocol != nil {
					return fmt.Sprintf("Custom (%s)", strings.ToUpper(*rule.IpProtocol))
				} else {
					return "Custom Protocol"
				}
			}
		}
	}
	return ""
}
