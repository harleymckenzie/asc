package ec2

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter/builder"
)

// FieldValueGetter is a function that returns the value of a field for the given instance.
type FieldValueGetter func(instance any) (string, error)

// fieldValueGetters is a map of field names to their respective getter functions.
var fieldValueGetters = map[string]FieldValueGetter{
	"Name":                getInstanceName,
	"Instance ID":         getInstanceID,
	"State":               getState,
	"AMI ID":              getAMIID,
	"AMI Name":            getAMIName,
	"Launch Time":         getLaunchTime,
	"Instance Type":       getInstanceType,
	"Placement Group":     getPlacementGroup,
	"Root Device Type":    getRootDeviceType,
	"Root Device Name":    getRootDeviceName,
	"Virtualization Type": getVirtualizationType,
	"vCPUs":               getVCPUs,
	"Public IP":           getPublicIP,
	"Private IP":          getPrivateIP,
	"Subnet ID":           getSubnetID,
	"VPC ID":              getVPCID,
	"Availability Zone":   getAvailabilityZone,
	"Security Group(s)":   getSecurityGroupNames,
	"Key Name":            getKeyName,
}

// PopulateFieldValues populates the values of the fields for the given instance.
func PopulateFieldValues(fields []builder.Field, instance any) ([]builder.Field, error) {
	var populated []builder.Field
	for _, field := range fields {
		value := GetFieldValue(field.Name, instance)
		populated = append(populated, builder.Field{
			Category: field.Category,
			Name:     field.Name,
			Value:    value,
		})
	}
	return populated, nil
}

// getFieldValue returns the value of a field for the given instance.
func GetFieldValue(fieldName string, instance any) string {
	if getter, exists := fieldValueGetters[fieldName]; exists {
		value, err := getter(instance)
		if err != nil {
			return fmt.Sprintf("\033[31merror: Failed to get field value for %s: %v\033[0m", fieldName, err)
		}
		return value
	}
	return fmt.Sprintf("\033[31merror: Field \"%s\" not found in fieldValueGetters\033[0m", fieldName)
}

// Individual field value getters
func getInstanceName(instance any) (string, error) {
	return GetTagValue("Name", instance)
}

func getInstanceID(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).InstanceId), nil
}

func getState(instance any) (string, error) {
	return format.Status(string(instance.(types.Instance).State.Name)), nil
}

func getAMIID(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).ImageId), nil
}

func getAMIName(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).ImageId), nil
}

func getLaunchTime(instance any) (string, error) {
	return instance.(types.Instance).LaunchTime.Format(time.RFC3339), nil
}

func getInstanceType(instance any) (string, error) {
	return string(instance.(types.Instance).InstanceType), nil
}

func getPlacementGroup(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).Placement.GroupName), nil
}

func getRootDeviceType(instance any) (string, error) {
	return string(instance.(types.Instance).RootDeviceType), nil
}

func getRootDeviceName(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).RootDeviceName), nil
}

func getVirtualizationType(instance any) (string, error) {
	return string(instance.(types.Instance).VirtualizationType), nil
}

func getVCPUs(instance any) (string, error) {
	return strconv.Itoa(int(aws.ToInt32(instance.(types.Instance).CpuOptions.CoreCount))), nil
}

func getPublicIP(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).PublicIpAddress), nil
}

func getPrivateIP(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).PrivateIpAddress), nil
}

func getSubnetID(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).SubnetId), nil
}

func getVPCID(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).VpcId), nil
}

func getAvailabilityZone(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).Placement.AvailabilityZone), nil
}

func getSecurityGroupNames(instance any) (string, error) {
	securityGroups := instance.(types.Instance).SecurityGroups
	groupNames := make([]string, len(securityGroups))
	for i, group := range securityGroups {
		groupNames[i] = aws.ToString(group.GroupName)
	}
	return strings.Join(groupNames, "\n"), nil
}

func getKeyName(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).KeyName), nil
}

// Tag value getter

func GetTagValue(tagKey string, instance any) (string, error) {
	for _, tag := range instance.(types.Instance).Tags {
		if aws.ToString(tag.Key) == tagKey {
			return aws.ToString(tag.Value), nil
		}
	}
	return "", nil
}
