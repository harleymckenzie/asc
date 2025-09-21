package ec2

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// fieldValueGetters is a map of field names to their respective getter functions.
var ec2FieldValueGetters = map[string]FieldValueGetter{
	"Name":                getInstanceName,
	"Instance ID":         getInstanceID,
	"State":               getInstanceState,
	"AMI ID":              getInstanceAMIID,
	"Launch Time":         getInstanceLaunchTime,
	"Instance Type":       getInstanceType,
	"Placement Group":     getInstancePlacementGroup,
	"Root Device Type":    getInstanceRootDeviceType,
	"Root Device Name":    getInstanceRootDeviceName,
	"Virtualization Type": getInstanceVirtualizationType,
	"vCPUs":               getInstanceVCPUs,
	"Public IP":           getInstancePublicIP,
	"Private IP":          getInstancePrivateIP,
	"Subnet ID":           getInstanceSubnetID,
	"VPC ID":              getInstanceVPCID,
	"Availability Zone":   getInstanceAvailabilityZone,
	"Security Group(s)":   getInstanceSecurityGroupNames,
	"Key Name":            getInstanceKeyName,
}

// getInstanceFieldValue returns the value of a field for an EC2 instance
func getInstanceFieldValue(fieldName string, instance types.Instance, svc *EC2Service) (string, error) {
	// Special handling for fields that need service context
	if fieldName == "AMI Name" && svc != nil {
		return getInstanceAMIName(instance, svc)
	}
	

	if getter, exists := ec2FieldValueGetters[fieldName]; exists {
		value, err := getter(instance)
		if err != nil {
			return "", fmt.Errorf("failed to get field value for %s: %w", fieldName, err)
		}
		return value, nil
	}
	return "", fmt.Errorf("field %s not found in instance fieldValueGetters", fieldName)
}

// Individual field value getters
func getInstanceName(instance any) (string, error) {
	return GetTagValue("Name", instance)
}

func getInstanceID(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).InstanceId), nil
}

func getInstanceState(instance any) (string, error) {
	return format.Status(string(instance.(types.Instance).State.Name)), nil
}

func getInstanceAMIID(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).ImageId), nil
}

// getInstanceAMINameWithService resolves the actual AMI name using the service
func getInstanceAMIName(instance types.Instance, svc *EC2Service) (string, error) {
	imageID := aws.ToString(instance.ImageId)
	if imageID == "" {
		return "", nil
	}

	images, err := svc.GetImages(context.TODO(), &ascTypes.GetImagesInput{
		ImageIds: []string{imageID},
	})
	if err != nil {
		return "", err
	}
	if len(images) == 0 {
		return "", nil
	}
	return aws.ToString(images[0].Name), nil
}

func getInstanceLaunchTime(instance any) (string, error) {
	return instance.(types.Instance).LaunchTime.Format(time.RFC3339), nil
}

func getInstanceType(instance any) (string, error) {
	return string(instance.(types.Instance).InstanceType), nil
}

func getInstancePlacementGroup(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).Placement.GroupName), nil
}

func getInstanceRootDeviceType(instance any) (string, error) {
	return string(instance.(types.Instance).RootDeviceType), nil
}

func getInstanceRootDeviceName(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).RootDeviceName), nil
}

func getInstanceVirtualizationType(instance any) (string, error) {
	return string(instance.(types.Instance).VirtualizationType), nil
}

func getInstanceVCPUs(instance any) (string, error) {
	return strconv.Itoa(int(aws.ToInt32(instance.(types.Instance).CpuOptions.CoreCount))), nil
}

func getInstancePublicIP(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).PublicIpAddress), nil
}

func getInstancePrivateIP(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).PrivateIpAddress), nil
}

func getInstanceSubnetID(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).SubnetId), nil
}

func getInstanceVPCID(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).VpcId), nil
}

func getInstanceAvailabilityZone(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).Placement.AvailabilityZone), nil
}

func getInstanceSecurityGroupNames(instance any) (string, error) {
	securityGroups := instance.(types.Instance).SecurityGroups
	groupNames := make([]string, len(securityGroups))
	for i, group := range securityGroups {
		groupNames[i] = aws.ToString(group.GroupName)
	}
	return strings.Join(groupNames, "\n"), nil
}

func getInstanceKeyName(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).KeyName), nil
}
