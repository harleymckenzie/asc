package vpc

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type IGWFieldValueGetter func(instance any) (string, error)

// Internet Gateway field getters
var igwFieldValueGetters = map[string]IGWFieldValueGetter{
	"Internet Gateway ID": getIGWID,
	"VPC ID":              getIGWVPCID,
	"State":               getIGWState,
	"Owner":               getIGWOwner,
}

// GetIGWFieldValue returns the value of a field for the given Internet Gateway instance.
func GetIGWFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.InternetGateway:
		return getIGWFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// getIGWFieldValue returns the value of a field for an Internet Gateway
func getIGWFieldValue(fieldName string, igw types.InternetGateway) (string, error) {
	if getter, exists := igwFieldValueGetters[fieldName]; exists {
		return getter(igw)
	}
	return "", fmt.Errorf("field %s not found in igwFieldValueGetters", fieldName)
}

// GetIGWTagValue returns the value of a tag for the given instance.
func GetIGWTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.InternetGateway:
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
// Internet Gateway field getters
// -----------------------------------------------------------------------------

func getIGWID(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.InternetGateway).InternetGatewayId), nil
}

func getIGWVPCID(instance any) (string, error) {
	igw := instance.(types.InternetGateway)
	if len(igw.Attachments) == 0 {
		return "", nil
	}
	var vpcs []string
	for _, a := range igw.Attachments {
		vpcs = append(vpcs, format.StringOrEmpty(a.VpcId))
	}
	return strings.Join(vpcs, ", "), nil
}

func getIGWState(instance any) (string, error) {
	igw := instance.(types.InternetGateway)
	if len(igw.Attachments) == 0 {
		return "-", nil
	}
	if igw.Attachments[0].State == "available" {
		return format.Status("Attached"), nil
	}
	return format.Status(string(igw.Attachments[0].State)), nil
}

func getIGWOwner(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.InternetGateway).OwnerId), nil
}
