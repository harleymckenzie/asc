package elasticache

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type FieldValueGetter func(instance any) (string, error)

// GetFieldValue returns the value of a field for the given instance.
func GetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.CacheCluster:
		return getCacheClusterFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// GetTagValue returns the value of a tag for the given instance.
func GetTagValue(tagKey string, instance any) (string, error) {
	switch instance.(type) {
	case types.CacheCluster:
		// ElastiCache clusters don't have tags in the same way as EC2 instances
		// This is a placeholder for future implementation if needed
		// Tags may require separate API calls to retrieve
		return "", nil
	default:
		return "", fmt.Errorf("unsupported instance type for tags: %T", instance)
	}
}
