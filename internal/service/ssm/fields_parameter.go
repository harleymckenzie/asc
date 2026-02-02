package ssm

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// parameterFieldValueGetters maps field names to their getter functions.
var parameterFieldValueGetters = map[string]FieldValueGetter{
	"Name":               getParameterName,
	"Type":               getParameterType,
	"Value":              getParameterValue,
	"Version":            getParameterVersion,
	"Last Modified Date": getParameterLastModifiedDate,
	"ARN":                getParameterARN,
	"Data Type":          getParameterDataType,
}

// parameterMetadataFieldValueGetters maps field names to their getter functions for metadata.
var parameterMetadataFieldValueGetters = map[string]FieldValueGetter{
	"Name":               getMetadataName,
	"Type":               getMetadataType,
	"Last Modified Date": getMetadataLastModifiedDate,
	"Last Modified User": getMetadataLastModifiedUser,
	"Version":            getMetadataVersion,
	"Tier":               getMetadataTier,
	"Description":        getMetadataDescription,
}

// getParameterFieldValue returns the value of a field for a Parameter.
func getParameterFieldValue(fieldName string, param types.Parameter) (string, error) {
	if getter, exists := parameterFieldValueGetters[fieldName]; exists {
		return getter(param)
	}
	return "", fmt.Errorf("field %s not found in parameter fieldValueGetters", fieldName)
}

// getParameterMetadataFieldValue returns the value of a field for ParameterMetadata.
func getParameterMetadataFieldValue(fieldName string, param types.ParameterMetadata) (string, error) {
	if getter, exists := parameterMetadataFieldValueGetters[fieldName]; exists {
		return getter(param)
	}
	return "", fmt.Errorf("field %s not found in parameterMetadata fieldValueGetters", fieldName)
}

// Parameter field getters

func getParameterName(param any) (string, error) {
	return aws.ToString(param.(types.Parameter).Name), nil
}

func getParameterType(param any) (string, error) {
	return string(param.(types.Parameter).Type), nil
}

func getParameterValue(param any) (string, error) {
	p := param.(types.Parameter)
	// Mask SecureString values
	if p.Type == types.ParameterTypeSecureString {
		return "****", nil
	}
	return aws.ToString(p.Value), nil
}

func getParameterVersion(param any) (string, error) {
	return fmt.Sprintf("%d", param.(types.Parameter).Version), nil
}

func getParameterLastModifiedDate(param any) (string, error) {
	t := param.(types.Parameter).LastModifiedDate
	if t == nil {
		return "", nil
	}
	return t.Format(time.RFC3339), nil
}

func getParameterARN(param any) (string, error) {
	return aws.ToString(param.(types.Parameter).ARN), nil
}

func getParameterDataType(param any) (string, error) {
	return aws.ToString(param.(types.Parameter).DataType), nil
}

// ParameterMetadata field getters

func getMetadataName(param any) (string, error) {
	return aws.ToString(param.(types.ParameterMetadata).Name), nil
}

func getMetadataType(param any) (string, error) {
	return string(param.(types.ParameterMetadata).Type), nil
}

func getMetadataLastModifiedDate(param any) (string, error) {
	t := param.(types.ParameterMetadata).LastModifiedDate
	if t == nil {
		return "", nil
	}
	return t.Format(time.RFC3339), nil
}

func getMetadataLastModifiedUser(param any) (string, error) {
	return aws.ToString(param.(types.ParameterMetadata).LastModifiedUser), nil
}

func getMetadataVersion(param any) (string, error) {
	return fmt.Sprintf("%d", param.(types.ParameterMetadata).Version), nil
}

func getMetadataTier(param any) (string, error) {
	return string(param.(types.ParameterMetadata).Tier), nil
}

func getMetadataDescription(param any) (string, error) {
	return aws.ToString(param.(types.ParameterMetadata).Description), nil
}

// GetDecryptedValue returns the actual value for display when --decrypt flag is used.
func GetDecryptedValue(param types.Parameter) string {
	return aws.ToString(param.Value)
}
