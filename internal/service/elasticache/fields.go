package elasticache

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type FieldValueGetter func(instance any) (string, error)

// GetFieldValue returns the value of a field for the given instance.
func GetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case CacheClusterWithTags:
		return getCacheClusterFieldValue(fieldName, v.CacheCluster)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// GetTagValue returns the value of a tag for the given instance.
func GetTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case CacheClusterWithTags:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
		return "", nil
	default:
		return "", fmt.Errorf("unsupported instance type for tags: %T", instance)
	}
}
