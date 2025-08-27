package cloudformation

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type FieldValueGetter func(instance any) (string, error)

// CloudFormation Stack field getters
var stackFieldValueGetters = map[string]FieldValueGetter{
	"Stack Name":             getStackName,
	"Stack ID":               getStackID,
	"Description":            getStackDescription,
	"Status":                 getStackStatus,
	"Detailed Status":        getStackDetailedStatus,
	"Status Reason":          getStackStatusReason,
	"Root Stack":             getStackRootStack,
	"Parent Stack":           getStackParentStack,
	"Creation Time":          getStackCreationTime,
	"Last Updated":           getStackLastUpdated,
	"Deletion Time":          getStackDeletionTime,
	"Drift Status":           getStackDriftStatus,
	"Deletion Mode":          getStackDeletionMode,
	"Last Drift Check":       getStackLastDriftCheck,
	"Termination Protection": getStackTerminationProtection,
	"IAM Role":               getStackIAMRole,
}

// GetFieldValue returns the value of a field for the given instance.
func GetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.Stack:
		return getStackFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// getStackFieldValue returns the value of a field for a CloudFormation stack
func getStackFieldValue(fieldName string, stack types.Stack) (string, error) {
	if getter, exists := stackFieldValueGetters[fieldName]; exists {
		return getter(stack)
	}
	return "", fmt.Errorf("field %s not found in stackFieldValueGetters", fieldName)
}

// GetTagValue returns the value of a tag for the given instance.
func GetTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.Stack:
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
// CloudFormation Stack field getters
// -----------------------------------------------------------------------------

func getStackName(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).StackName), nil
}

func getStackID(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).StackId), nil
}

func getStackDescription(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).Description), nil
}

func getStackStatus(instance any) (string, error) {
	return format.Status(string(instance.(types.Stack).StackStatus)), nil
}

func getStackDetailedStatus(instance any) (string, error) {
	return format.Status(string(instance.(types.Stack).DetailedStatus)), nil
}

func getStackStatusReason(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).StackStatusReason), nil
}

func getStackRootStack(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).RootId), nil
}

func getStackParentStack(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).ParentId), nil
}

func getStackCreationTime(instance any) (string, error) {
	stack := instance.(types.Stack)
	if stack.CreationTime == nil {
		return "", nil
	}
	return stack.CreationTime.Local().Format("2006-01-02 15:04:05 MST"), nil
}

func getStackLastUpdated(instance any) (string, error) {
	stack := instance.(types.Stack)
	if stack.LastUpdatedTime == nil {
		return "", nil
	}
	return stack.LastUpdatedTime.Local().Format("2006-01-02 15:04:05 MST"), nil
}

func getStackDeletionTime(instance any) (string, error) {
	stack := instance.(types.Stack)
	if stack.DeletionTime == nil {
		return "", nil
	}
	return stack.DeletionTime.Local().Format("2006-01-02 15:04:05 MST"), nil
}

func getStackDriftStatus(instance any) (string, error) {
	stack := instance.(types.Stack)
	if stack.DriftInformation == nil {
		return "", nil
	}
	return string(stack.DriftInformation.StackDriftStatus), nil
}

func getStackDeletionMode(instance any) (string, error) {
	return string(instance.(types.Stack).DeletionMode), nil
}

func getStackLastDriftCheck(instance any) (string, error) {
	stack := instance.(types.Stack)
	if stack.DriftInformation == nil || stack.DriftInformation.LastCheckTimestamp == nil {
		return "", nil
	}
	return stack.DriftInformation.LastCheckTimestamp.Local().Format("2006-01-02 15:04:05 MST"), nil
}

func getStackTerminationProtection(instance any) (string, error) {
	stack := instance.(types.Stack)
	if stack.EnableTerminationProtection == nil {
		return "", nil
	}
	return fmt.Sprintf("%t", *stack.EnableTerminationProtection), nil
}

func getStackIAMRole(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).RoleARN), nil
}
