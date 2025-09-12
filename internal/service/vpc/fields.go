package vpc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type FieldValueGetter func(instance any) (string, error)

// GetFieldValue returns the value of a field for the given instance.
// This function routes field requests to the appropriate type-specific handler.
func GetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.Vpc:
		return getVPCFieldValue(fieldName, v)
	case types.NetworkAcl:
		return getNACLFieldValue(fieldName, v)
	case types.Subnet:
		return getSubnetFieldValue(fieldName, v)
	case types.RouteTable:
		return getRouteTableFieldValue(fieldName, v)
	case types.Route:
		return getRouteFieldValue(fieldName, v)
	case types.InternetGateway:
		return getIGWFieldValue(fieldName, v)
	case types.NatGateway:
		return getNATFieldValue(fieldName, v)
	case types.ManagedPrefixList:
		return getPrefixFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// GetTagValue returns the value of a tag for the given instance.
// This function handles tag retrieval for all supported VPC resource types.
func GetTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.Vpc:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.NetworkAcl:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.Subnet:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.RouteTable:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.InternetGateway:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.NatGateway:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.ManagedPrefixList:
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
