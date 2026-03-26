package efs

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs/types"
)

type FieldValueGetter func(instance any) (string, error)

func GetFieldValue(fieldName string, instance any) (string, error) {
	switch instance.(type) {
	case types.FileSystemDescription:
		if getter, ok := fileSystemFieldValueGetters[fieldName]; ok {
			return getter(instance)
		}
		return "", fmt.Errorf("unknown field: %s", fieldName)
	default:
		return "", fmt.Errorf("unsupported type: %T", instance)
	}
}

func GetTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.FileSystemDescription:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	default:
		return "", fmt.Errorf("unsupported type for tags: %T", instance)
	}
	return "", nil
}
