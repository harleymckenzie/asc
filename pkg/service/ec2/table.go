package ec2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/pkg/shared/format"
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

// GetAttributeValue is a function that returns the value of a field in a detailed table.
func GetAttributeValue(fieldID string, instance any) (string, error) {
	inst, ok := instance.(types.Instance)
	if !ok {
		return "", fmt.Errorf("instance is not a types.Instance")
	}
	attr, ok := availableAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q", fieldID)
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

// GetVolumeAttributeValue is a function that returns the value of a field in a detailed table.
func GetVolumeAttributeValue(fieldID string, volume any) (string, error) {
	vol, ok := volume.(types.Volume)
	if !ok {
		return "", fmt.Errorf("volume is not a types.Volume")
	}
	attr, ok := volumeAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q", fieldID)
	}
	return attr.GetValue(&vol), nil
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
		"Volume Type": {
			GetValue: func(v *types.Volume) string {
				return string(v.VolumeType)
			},
		},
		"Size": {
			GetValue: func(v *types.Volume) string {
				if v.Size == nil {
					return ""
				}
				return strconv.Itoa(int(*v.Size))
			},
		},
		"State": {
			GetValue: func(v *types.Volume) string {
				return string(v.State)
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
		"Fast Restored": {
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
		"Encrypted": {
			GetValue: func(v *types.Volume) string {
				if v.Encrypted == nil {
					return ""
				}
				return strconv.FormatBool(*v.Encrypted)
			},
		},
		"KMS Key ID": {
			GetValue: func(v *types.Volume) string {
				return aws.ToString(v.KmsKeyId)
			},
		},
	}
}

// GetSnapshotAttributeValue is a function that returns the value of a field in a detailed table.
func GetSnapshotAttributeValue(fieldID string, snapshot any) (string, error) {
	snap, ok := snapshot.(types.Snapshot)
	if !ok {
		return "", fmt.Errorf("snapshot is not a types.Snapshot")
	}
	attr, ok := snapshotAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q", fieldID)
	}
	return attr.GetValue(&snap), nil
}

// snapshotAttributes is a function that returns a map of attributes for a detailed table.
func snapshotAttributes() map[string]SnapshotAttribute {
	return map[string]SnapshotAttribute{
		"Snapshot ID": {
			GetValue: func(s *types.Snapshot) string {
				return aws.ToString(s.SnapshotId)
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
	attr, ok := imageAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q", fieldID)
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

// GetSecurityGroupAttributeValue returns the value of a field for a SecurityGroup.
func GetSecurityGroupAttributeValue(fieldID string, group any) (string, error) {
	g, ok := group.(types.SecurityGroup)
	if !ok {
		return "", fmt.Errorf("group is not a types.SecurityGroup")
	}
	attr, ok := securityGroupAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q", fieldID)
	}
	return attr.GetValue(&g), nil
}

// securityGroupAttributes returns a map of field IDs to SecurityGroupAttribute.
func securityGroupAttributes() map[string]SecurityGroupAttribute {
	return map[string]SecurityGroupAttribute{
		"Group ID": {
			GetValue: func(g *types.SecurityGroup) string {
				return aws.ToString(g.GroupId)
			},
		},
		"Group Name": {
			GetValue: func(g *types.SecurityGroup) string {
				return aws.ToString(g.GroupName)
			},
		},
		"Description": {
			GetValue: func(g *types.SecurityGroup) string {
				return aws.ToString(g.Description)
			},
		},
		"VPC ID": {
			GetValue: func(g *types.SecurityGroup) string {
				return aws.ToString(g.VpcId)
			},
		},
		"Owner ID": {
			GetValue: func(g *types.SecurityGroup) string {
				return aws.ToString(g.OwnerId)
			},
		},
		"Ingress Count": {
			GetValue: func(g *types.SecurityGroup) string {
				return strconv.Itoa(len(g.IpPermissions))
			},
		},
		"Egress Count": {
			GetValue: func(g *types.SecurityGroup) string {
				return strconv.Itoa(len(g.IpPermissionsEgress))
			},
		},
		"Tag Count": {
			GetValue: func(g *types.SecurityGroup) string {
				return strconv.Itoa(len(g.Tags))
			},
		},
	}
}
