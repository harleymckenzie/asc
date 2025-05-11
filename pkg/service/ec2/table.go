package ec2

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/pkg/shared/format"
)


// Attribute is a struct that defines a field in a detailed table.
type Attribute struct {
    GetValue func(*types.Instance) string
}

// GetAttributeValue is a function that returns the value of a field in a detailed table.
func GetAttributeValue(fieldID string, instance any) string {
	inst, ok := instance.(types.Instance)
	if !ok {
		fmt.Println("Instance is not a types.Instance")
        return ""
	}
	attr := availableAttributes()[fieldID]
	return attr.GetValue(&inst)
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
