package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
)

// FieldValueGetter is a function that returns the value of a field for the given instance.
type FieldValueGetter func(instance any) (string, error)

// GetFieldValue returns the value of a field for the given instance.
// This function routes field requests to the appropriate type-specific handler.
func GetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.Instance:
		return getInstanceFieldValue(fieldName, v)
	case types.Image:
		return getImageFieldValue(fieldName, v)
	case types.Volume:
		return getVolumeFieldValue(fieldName, v)
	case types.Snapshot:
		return getSnapshotFieldValue(fieldName, v)
	case types.SecurityGroup:
		return getSecurityGroupFieldValue(fieldName, v)
	case types.SecurityGroupRule:
		return getSecurityGroupRuleFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// GetTagValue returns the value of a tag for the given instance.
// Currently supports EC2 instances only - other resource types may be added as needed.
func GetTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.Instance:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.Image:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.Volume:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.Snapshot:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.SecurityGroup:
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

// PopulateFieldValues populates the values of the fields for the given instance.
// This function is used to convert field definitions into populated field values for display.
func PopulateFieldValues(fields []tablewriter.Field, instance any) ([]tablewriter.Field, error) {
	var populated []tablewriter.Field
	for _, field := range fields {
		value, err := GetFieldValue(field.Name, instance)
		if err != nil {
			return nil, fmt.Errorf("failed to get field value for %s: %w", field.Name, err)
		}
		populated = append(populated, tablewriter.Field{
			Category: field.Category,
			Name:     field.Name,
			Value:    value,
		})
	}
	return populated, nil
}
