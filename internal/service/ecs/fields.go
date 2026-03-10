package ecs

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type FieldValueGetter func(instance any) (string, error)

// GetFieldValue returns the value of a field for the given instance.
func GetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.Cluster:
		return getClusterFieldValue(fieldName, v)
	case types.Service:
		return getServiceFieldValue(fieldName, v)
	case types.Task:
		return getTaskFieldValue(fieldName, v)
	case types.TaskDefinition:
		return getTaskDefinitionFieldValue(fieldName, v)
	case TaskDefinitionFamily:
		return getTaskDefinitionFamilyFieldValue(fieldName, v)
	case TaskDefinitionRevision:
		return getTaskDefinitionRevisionFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// GetTagValue returns the value of a tag for the given instance.
func GetTagValue(tagKey string, instance any) (string, error) {
	var tags []types.Tag

	switch v := instance.(type) {
	case types.Cluster:
		tags = v.Tags
	case types.Service:
		tags = v.Tags
	case types.Task:
		tags = v.Tags
	case types.TaskDefinition:
		// Task definitions don't typically have tags in this structure
		return "", nil
	case TaskDefinitionFamily:
		return "", nil
	case TaskDefinitionRevision:
		return "", nil
	default:
		return "", fmt.Errorf("unsupported instance type for tags: %T", instance)
	}

	for _, tag := range tags {
		if aws.ToString(tag.Key) == tagKey {
			return aws.ToString(tag.Value), nil
		}
	}
	return "", nil
}

// TaskDefinitionFamily represents a task definition family for listing.
type TaskDefinitionFamily struct {
	Name string
}

// TaskDefinitionRevision represents a task definition revision for listing.
type TaskDefinitionRevision struct {
	ARN      string
	Family   string
	Revision string
}
