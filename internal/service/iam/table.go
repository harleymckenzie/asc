package iam

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// Attribute is a struct that defines a field in a detailed table.
type RoleAttribute struct {
	GetValue func(*types.Role) string
}

func GetRoleAttributeValue(fieldID string, instance any) (string, error) {
	role, ok := instance.(types.Role)
	if !ok {
		return "", fmt.Errorf("instance is not a types.Role")
	}
	attr, ok := availableRoleAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q", fieldID)
	}
	return attr.GetValue(&role), nil
}

func availableRoleAttributes() map[string]RoleAttribute {
	return map[string]RoleAttribute{
		"Name": {
			GetValue: func(i *types.Role) string {
				return format.StringOrEmpty(i.RoleName)
			},
		},
		"Arn": {
			GetValue: func(i *types.Role) string {
				return format.StringOrEmpty(i.Arn)
			},
		},
		"Role ID": {
			GetValue: func(i *types.Role) string {
				return format.StringOrEmpty(i.RoleId)
			},
		},
		"Path": {
			GetValue: func(i *types.Role) string {
				return format.StringOrEmpty(i.Path)
			},
		},
		"Description": {
			GetValue: func(i *types.Role) string {
				return format.StringOrEmpty(i.Description)
			},
		},
		"Creation Time": {
			GetValue: func(i *types.Role) string {
				return format.TimeToStringOrEmpty(i.CreateDate)
			},
		},
		"Assume Role Policy Document": {
			GetValue: func(i *types.Role) string {
				return format.DecodeAndFormatJSON(i.AssumeRolePolicyDocument)
			},
		},
		"Max Session Duration": {
			GetValue: func(i *types.Role) string {
				return format.Int32ToStringOrEmpty(i.MaxSessionDuration)
			},
		},
		"Permissions Boundary": {
			GetValue: func(i *types.Role) string {
				if i.PermissionsBoundary == nil {
					return ""
				}
				return format.StringOrEmpty(i.PermissionsBoundary.PermissionsBoundaryArn)
			},
		},
		"Role Last Used": {
			GetValue: func(i *types.Role) string {
				if i.RoleLastUsed == nil {
					return ""
				}
				return format.TimeToStringOrEmpty(i.RoleLastUsed.LastUsedDate)
			},
		},
		"Last Activity": {
			GetValue: func(i *types.Role) string {
				// Show the relative time since the last activity.
				// Supposedly the ListRoles API returns a subset of fields that apply to the role.
				// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/iam@v1.43.0#Client.ListRoles
				if i.RoleLastUsed == nil || i.RoleLastUsed.LastUsedDate == nil {
					return ""
				}
				return format.TimeToStringRelative(i.RoleLastUsed.LastUsedDate)
			},
		},
	}
}
