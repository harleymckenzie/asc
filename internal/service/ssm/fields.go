package ssm

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// FieldValueGetter is a function that returns the value of a field for the given parameter.
type FieldValueGetter func(param any) (string, error)

// GetFieldValue returns the value of a field for the given parameter.
// This function routes field requests to the appropriate type-specific handler.
func GetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.Parameter:
		return getParameterFieldValue(fieldName, v)
	case types.ParameterMetadata:
		return getParameterMetadataFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// GetTagValue returns the value of a tag for the given parameter.
// Note: SSM Parameters accessed via GetParameter don't include tags directly.
// Tags must be fetched separately via ListTagsForResource.
func GetTagValue(tagKey string, instance any) (string, error) {
	// SSM Parameters don't have tags in the Parameter struct
	// Tags would need to be fetched separately
	return "", nil
}
