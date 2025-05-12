package cloudformation

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/harleymckenzie/asc/pkg/shared/format"
)

// Attribute is a struct that defines a field in a detailed table.
type Attribute struct {
	GetValue func(*types.Stack) string
}

// GetAttributeValue returns the value for a given field and stack instance.
func GetAttributeValue(fieldID string, instance any) (string, error) {
	inst, ok := instance.(types.Stack)
	if !ok {
		return "", fmt.Errorf("instance is not a types.Stack")
	}
	attr, ok := availableAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q", fieldID)
	}
	return attr.GetValue(&inst), nil
}

// availableAttributes returns the fields for CloudFormation stack list tables.
func availableAttributes() map[string]Attribute {
	return map[string]Attribute{
		"Stack Name": {
			GetValue: func(i *types.Stack) string {
				return aws.ToString(i.StackName)
			},
		},
		"Stack ID": {
			GetValue: func(i *types.Stack) string {
				return aws.ToString(i.StackId)
			},
		},
		"Description": {
			GetValue: func(i *types.Stack) string {
				return aws.ToString(i.Description)
			},
		},
		"Status": {
			GetValue: func(i *types.Stack) string {
				return format.Status(string(i.StackStatus))
			},
		},
		"Detailed Status": {
			GetValue: func(i *types.Stack) string {
				return format.Status(string(i.DetailedStatus))
			},
		},
		"Status Reason": {
			GetValue: func(i *types.Stack) string {
				return aws.ToString(i.StackStatusReason)
			},
		},
		"Root Stack": {
			GetValue: func(i *types.Stack) string {
				return aws.ToString(i.RootId)
			},
		},
		"Parent Stack": {
			GetValue: func(i *types.Stack) string {
				return aws.ToString(i.ParentId)
			},
		},
		"Creation Time": {
			GetValue: func(i *types.Stack) string {
				return i.CreationTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"Last Updated": {
			GetValue: func(i *types.Stack) string {
				if i.LastUpdatedTime == nil {
					return ""
				}
				return i.LastUpdatedTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"Deletion Time": {
			GetValue: func(i *types.Stack) string {
				if i.DeletionTime == nil {
					return ""
				}
				return i.DeletionTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"Drift Status": {
			GetValue: func(i *types.Stack) string {
				return string(i.DriftInformation.StackDriftStatus)
			},
		},
		"Deletion Mode": {
			GetValue: func(i *types.Stack) string {
				return string(i.DeletionMode)
			},
		},
		"Last Drift Check": {
			GetValue: func(i *types.Stack) string {
				return i.DriftInformation.LastCheckTimestamp.Local().
					Format("2006-01-02 15:04:05 MST")
			},
		},
		"Termination Protection": {
			GetValue: func(i *types.Stack) string {
				return fmt.Sprintf("%t", *i.EnableTerminationProtection)
			},
		},
		"IAM Role": {
			GetValue: func(i *types.Stack) string {
				return aws.ToString(i.RoleARN)
			},
		},
	}
}
