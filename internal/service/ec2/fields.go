package ec2

import (
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter/builder"
)

// getFields returns a list of Field objects for the given instance.
func GetFields(instance any) []builder.Field {
	return []builder.Field{
		{Name: "Instance ID", Category: "Instance Details", Visible: true},
		{Name: "State", Category: "Instance Details", Visible: true},
		{Name: "AMI ID", Category: "Instance Details", Visible: true},
		{Name: "AMI Name", Category: "Instance Details", Visible: true},
		{Name: "Launch Time", Category: "Instance Details", Visible: true},
		{Name: "Instance Type", Category: "Instance Details", Visible: true},
		{Name: "Placement Group", Category: "Instance Details", Visible: true},
		{Name: "Root Device Type", Category: "Instance Details", Visible: true},
		{Name: "Root Device Name", Category: "Instance Details", Visible: true},
		{Name: "Virtualization Type", Category: "Instance Details", Visible: true},
		{Name: "vCPUs", Category: "Instance Details", Visible: true},
		{Name: "Public IP", Category: "Network", Visible: true},
		{Name: "Private IP", Category: "Network", Visible: true},
		{Name: "Subnet ID", Category: "Network", Visible: true},
		{Name: "VPC ID", Category: "Network", Visible: true},
		{Name: "Availability Zone", Category: "Network", Visible: true},
		{Name: "Security Group(s)", Category: "Security", Visible: true},
		{Name: "Key Name", Category: "Security", Visible: true},
	}
}

// FieldValueGetter is a function that returns the value of a field for the given instance.
type FieldValueGetter func(instance any) string

// fieldValueGetters is a map of field names to their respective getter functions.
var fieldValueGetters = map[string]FieldValueGetter{
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

// getFieldValue returns the value of a field for the given instance.
func GetFieldValue(fieldName string, instance any) string {
	if getter, exists := fieldValueGetters[fieldName]; exists {
		return getter(instance)
	}
	return "-"
}

// PopulateFieldValues populates the values of the fields for the given instance.
func PopulateFieldValues(fields []builder.Field, instance any) []builder.Field {
	var populated []builder.Field
	for _, field := range fields {
		value := GetFieldValue(field.Name, instance)
		populated = append(populated, builder.Field{
			Category: field.Category,
			Name:     field.Name,
			Value:    value,
		})
	}
	return populated
}

// Individual field value getters

func getInstanceID(instance any) string {
	return string(*instance.(types.Instance).InstanceId)
}

func getState(instance any) string {
	return string(instance.(types.Instance).State.Name)
}

func getAMIID(instance any) string {
	return string(*instance.(types.Instance).ImageId)
}

func getAMIName(instance any) string {
	return string(*instance.(types.Instance).ImageId)
}

func getLaunchTime(instance any) string {
	return instance.(types.Instance).LaunchTime.Format(time.RFC3339)
}

func getInstanceType(instance any) string {
	return string(instance.(types.Instance).InstanceType)
}

func getPlacementGroup(instance any) string {
	return string(*instance.(types.Instance).Placement.GroupName)
}

func getRootDeviceType(instance any) string {
	return string(instance.(types.Instance).RootDeviceType)
}

func getRootDeviceName(instance any) string {
	return string(*instance.(types.Instance).RootDeviceName)
}

func getVirtualizationType(instance any) string {
	return string(instance.(types.Instance).VirtualizationType)
}

func getVCPUs(instance any) string {
	return strconv.Itoa(int(*instance.(types.Instance).CpuOptions.CoreCount))
}

func getPublicIP(instance any) string {
	return string(*instance.(types.Instance).PublicIpAddress)
}

func getPrivateIP(instance any) string {
	return string(*instance.(types.Instance).PrivateIpAddress)
}

func getSubnetID(instance any) string {
	return string(*instance.(types.Instance).SubnetId)
}

func getVPCID(instance any) string {
	return string(*instance.(types.Instance).VpcId)
}

func getAvailabilityZone(instance any) string {
	return string(*instance.(types.Instance).Placement.AvailabilityZone)
}

func getSecurityGroupNames(instance any) string {
	securityGroups := instance.(types.Instance).SecurityGroups
	groupNames := make([]string, len(securityGroups))
	for i, group := range securityGroups {
		groupNames[i] = *group.GroupName
	}
	return strings.Join(groupNames, "\n")
}

func getKeyName(instance any) string {
	return *instance.(types.Instance).KeyName
}
