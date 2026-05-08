package iam

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type FieldValueGetter func(instance any) (string, error)

func GetFieldValue(fieldName string, instance any) (string, error) {
	switch instance.(type) {
	case types.Role:
		if getter, ok := roleFieldValueGetters[fieldName]; ok {
			return getter(instance)
		}
		return "", fmt.Errorf("unknown field: %s for type Role", fieldName)
	case types.Policy:
		if getter, ok := policyFieldValueGetters[fieldName]; ok {
			return getter(instance)
		}
		return "", fmt.Errorf("unknown field: %s for type Policy", fieldName)
	case types.User:
		if getter, ok := userFieldValueGetters[fieldName]; ok {
			return getter(instance)
		}
		return "", fmt.Errorf("unknown field: %s for type User", fieldName)
	default:
		return "", fmt.Errorf("unsupported type: %T", instance)
	}
}

func GetTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.Role:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.Policy:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.User:
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
