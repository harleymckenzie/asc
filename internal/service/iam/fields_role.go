package iam

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

var roleFieldValueGetters = map[string]FieldValueGetter{
	"Role Name":            getRoleName,
	"Role ID":              getRoleID,
	"ARN":                  getRoleARN,
	"Path":                 getRolePath,
	"Created Date":         getRoleCreatedDate,
	"Max Session Duration": getRoleMaxSessionDuration,
	"Description":          getRoleDescription,
	"Last Used Date":       getRoleLastUsedDate,
	"Last Used Region":     getRoleLastUsedRegion,
}

func getRoleName(instance any) (string, error) {
	return aws.ToString(instance.(types.Role).RoleName), nil
}

func getRoleID(instance any) (string, error) {
	return aws.ToString(instance.(types.Role).RoleId), nil
}

func getRoleARN(instance any) (string, error) {
	return aws.ToString(instance.(types.Role).Arn), nil
}

func getRolePath(instance any) (string, error) {
	return aws.ToString(instance.(types.Role).Path), nil
}

func getRoleCreatedDate(instance any) (string, error) {
	r := instance.(types.Role)
	if r.CreateDate != nil {
		return r.CreateDate.Format(time.RFC3339), nil
	}
	return "", nil
}

func getRoleMaxSessionDuration(instance any) (string, error) {
	r := instance.(types.Role)
	if r.MaxSessionDuration != nil {
		return fmt.Sprintf("%d seconds", *r.MaxSessionDuration), nil
	}
	return "", nil
}

func getRoleDescription(instance any) (string, error) {
	return aws.ToString(instance.(types.Role).Description), nil
}

func getRoleLastUsedDate(instance any) (string, error) {
	r := instance.(types.Role)
	if r.RoleLastUsed != nil && r.RoleLastUsed.LastUsedDate != nil {
		return r.RoleLastUsed.LastUsedDate.Format(time.RFC3339), nil
	}
	return "", nil
}

func getRoleLastUsedRegion(instance any) (string, error) {
	r := instance.(types.Role)
	if r.RoleLastUsed != nil {
		return aws.ToString(r.RoleLastUsed.Region), nil
	}
	return "", nil
}
