package elasticache

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
)

// FieldValueGetter is a function that returns the value of a field for the given instance.
type FieldValueGetter func(instance any) (string, error)

// GetFieldValue returns the value of a field for the given instance.
func GetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.CacheCluster:
		return getCacheClusterFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported type: %T", instance)
	}
}

// GetTagValue returns the value of a tag for the given instance.
func GetTagValue(tagName string, instance any) (string, error) {
	// ElastiCache clusters don't have tags in the same way as EC2 instances
	// This is a placeholder for future implementation if needed
	return "", nil
}