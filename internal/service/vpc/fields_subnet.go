package vpc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type SubnetFieldValueGetter func(instance any) (string, error)

// Subnet field getters
var subnetFieldValueGetters = map[string]SubnetFieldValueGetter{
	"Subnet ID":         getSubnetID,
	"VPC ID":            getSubnetVPCID,
	"CIDR Block":        getSubnetCIDRBlock,
	"Availability Zone": getSubnetAvailabilityZone,
	"State":             getSubnetState,
	"Available IPs":     getSubnetAvailableIPs,
	"Default For AZ":    getSubnetDefaultForAZ,
}

// GetSubnetFieldValue returns the value of a field for the given Subnet instance.
func GetSubnetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.Subnet:
		return getSubnetFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// getSubnetFieldValue returns the value of a field for a Subnet
func getSubnetFieldValue(fieldName string, subnet types.Subnet) (string, error) {
	if getter, exists := subnetFieldValueGetters[fieldName]; exists {
		return getter(subnet)
	}
	return "", fmt.Errorf("field %s not found in subnetFieldValueGetters", fieldName)
}

// GetSubnetTagValue returns the value of a tag for the given instance.
func GetSubnetTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.Subnet:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	default:
		return "", fmt.Errorf("unsupported instance type for tags: %T", instance)
	}
	return "", nil
}

// -----------------------------------------------------------------------------
// Subnet field getters
// -----------------------------------------------------------------------------

func getSubnetID(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.Subnet).SubnetId), nil
}

func getSubnetVPCID(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.Subnet).VpcId), nil
}

func getSubnetCIDRBlock(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.Subnet).CidrBlock), nil
}

func getSubnetAvailabilityZone(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.Subnet).AvailabilityZone), nil
}

func getSubnetState(instance any) (string, error) {
	return format.Status(string(instance.(types.Subnet).State)), nil
}

func getSubnetAvailableIPs(instance any) (string, error) {
	subnet := instance.(types.Subnet)
	return format.Int32ToStringOrDefault(subnet.AvailableIpAddressCount, "-"), nil
}

func getSubnetDefaultForAZ(instance any) (string, error) {
	subnet := instance.(types.Subnet)
	return format.BoolToLabel(subnet.DefaultForAz, "Yes", "No"), nil
}
