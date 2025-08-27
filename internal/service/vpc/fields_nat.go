package vpc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type NATFieldValueGetter func(instance any) (string, error)

// NAT Gateway field getters
var natFieldValueGetters = map[string]NATFieldValueGetter{
	"NAT Gateway ID":     getNATID,
	"Connectivity":       getNATConnectivity,
	"VPC ID":             getNATVPCID,
	"Subnet ID":          getNATSubnetID,
	"State":              getNATState,
	"Primary Public IP":  getNATPrimaryPublicIP,
	"Primary Private IP": getNATPrimaryPrivateIP,
	"Created":            getNATCreated,
}

// GetNATFieldValue returns the value of a field for the given NAT Gateway instance.
func GetNATFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.NatGateway:
		return getNATFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// getNATFieldValue returns the value of a field for a NAT Gateway
func getNATFieldValue(fieldName string, nat types.NatGateway) (string, error) {
	if getter, exists := natFieldValueGetters[fieldName]; exists {
		return getter(nat)
	}
	return "", fmt.Errorf("field %s not found in natFieldValueGetters", fieldName)
}

// GetNATTagValue returns the value of a tag for the given instance.
func GetNATTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.NatGateway:
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
// NAT Gateway field getters
// -----------------------------------------------------------------------------

func getNATID(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.NatGateway).NatGatewayId), nil
}

func getNATConnectivity(instance any) (string, error) {
	nat := instance.(types.NatGateway)
	if nat.ConnectivityType == "public" {
		return "Public", nil
	}
	return "Private", nil
}

func getNATVPCID(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.NatGateway).VpcId), nil
}

func getNATSubnetID(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.NatGateway).SubnetId), nil
}

func getNATState(instance any) (string, error) {
	return format.Status(string(instance.(types.NatGateway).State)), nil
}

func getNATPrimaryPublicIP(instance any) (string, error) {
	nat := instance.(types.NatGateway)
	if len(nat.NatGatewayAddresses) == 0 {
		return "", nil
	}
	return format.StringOrEmpty(nat.NatGatewayAddresses[0].PublicIp), nil
}

func getNATPrimaryPrivateIP(instance any) (string, error) {
	nat := instance.(types.NatGateway)
	if len(nat.NatGatewayAddresses) == 0 {
		return "", nil
	}
	return format.StringOrEmpty(nat.NatGatewayAddresses[0].PrivateIp), nil
}

func getNATCreated(instance any) (string, error) {
	return format.TimeToStringOrEmpty(instance.(types.NatGateway).CreateTime), nil
}
