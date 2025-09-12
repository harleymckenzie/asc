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

// getStackName returns the name of the CloudFormation stack
func getStackName(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).StackName), nil
}

// getStackID returns the unique identifier of the CloudFormation stack
func getStackID(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).StackId), nil
}

// getStackDescription returns the description of the CloudFormation stack
func getStackDescription(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).Description), nil
}

// getStackStatus returns the current status of the CloudFormation stack
func getStackStatus(instance any) (string, error) {
	return format.Status(string(instance.(types.Stack).StackStatus)), nil
}

// getStackDetailedStatus returns the detailed status of the CloudFormation stack
func getStackDetailedStatus(instance any) (string, error) {
	return format.Status(string(instance.(types.Stack).DetailedStatus)), nil
}

// getStackStatusReason returns the reason for the current stack status
func getStackStatusReason(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).StackStatusReason), nil
}

// getStackRootStack returns the root stack ID if this is a nested stack
func getStackRootStack(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).RootId), nil
}

// getStackParentStack returns the parent stack ID if this is a nested stack
func getStackParentStack(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).ParentId), nil
}

// getStackCreationTime returns the timestamp when the stack was created
func getStackCreationTime(instance any) (string, error) {
	stack := instance.(types.Stack)
	if stack.CreationTime == nil {
		return "", nil
	}
	return stack.CreationTime.Local().Format("2006-01-02 15:04:05 MST"), nil
}

// getStackLastUpdated returns the timestamp when the stack was last updated
func getStackLastUpdated(instance any) (string, error) {
	stack := instance.(types.Stack)
	if stack.LastUpdatedTime == nil {
		return "", nil
	}
	return stack.LastUpdatedTime.Local().Format("2006-01-02 15:04:05 MST"), nil
}

// getStackDeletionTime returns the timestamp when the stack deletion was initiated
func getStackDeletionTime(instance any) (string, error) {
	stack := instance.(types.Stack)
	if stack.DeletionTime == nil {
		return "", nil
	}
	return stack.DeletionTime.Local().Format("2006-01-02 15:04:05 MST"), nil
}

// getStackDriftStatus returns the current drift status of the stack
func getStackDriftStatus(instance any) (string, error) {
	stack := instance.(types.Stack)
	if stack.DriftInformation == nil {
		return "", nil
	}
	return string(stack.DriftInformation.StackDriftStatus), nil
}

// getStackDeletionMode returns the deletion mode of the stack
func getStackDeletionMode(instance any) (string, error) {
	return string(instance.(types.Stack).DeletionMode), nil
}

// getStackLastDriftCheck returns the timestamp of the last drift check
func getStackLastDriftCheck(instance any) (string, error) {
	stack := instance.(types.Stack)
	if stack.DriftInformation == nil || stack.DriftInformation.LastCheckTimestamp == nil {
		return "", nil
	}
	return stack.DriftInformation.LastCheckTimestamp.Local().Format("2006-01-02 15:04:05 MST"), nil
}

// getStackTerminationProtection returns whether termination protection is enabled
func getStackTerminationProtection(instance any) (string, error) {
	stack := instance.(types.Stack)
	if stack.EnableTerminationProtection == nil {
		return "", nil
	}
	return fmt.Sprintf("%t", *stack.EnableTerminationProtection), nil
}

// getStackIAMRole returns the IAM role ARN used by the stack
func getStackIAMRole(instance any) (string, error) {
	return aws.ToString(instance.(types.Stack).RoleARN), nil
}
