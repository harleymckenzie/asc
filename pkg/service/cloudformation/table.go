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
func GetAttributeValue(fieldID string, instance any) string {
	inst, ok := instance.(types.Stack)
	if !ok {
		fmt.Println("Instance is not a types.Stack")
		return ""
	}
	attr := availableAttributes()[fieldID]
	return attr.GetValue(&inst)
}

// availableAttributes returns the fields for CloudFormation stack list tables.
func availableAttributes() map[string]Attribute {
	return map[string]Attribute{
		"Stack Name": {
			GetValue: func(i *types.Stack) string {
				return aws.ToString(i.StackName)
			},
		},
		"Status": {
			GetValue: func(i *types.Stack) string {
				return format.Status(string(i.StackStatus))
			},
		},
		"Description": {
			GetValue: func(i *types.Stack) string {
				return aws.ToString(i.Description)
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
	}
}
