package iam

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

var policyFieldValueGetters = map[string]FieldValueGetter{
	"Policy Name":               getPolicyName,
	"Policy ID":                 getPolicyID,
	"ARN":                       getPolicyARN,
	"Path":                      getPolicyPath,
	"Attachment Count":          getPolicyAttachmentCount,
	"Default Version ID":        getPolicyDefaultVersionID,
	"Created Date":              getPolicyCreatedDate,
	"Updated Date":              getPolicyUpdatedDate,
	"Description":               getPolicyDescription,
	"Permissions Boundary Count": getPolicyPermissionsBoundaryCount,
}

func getPolicyName(instance any) (string, error) {
	return aws.ToString(instance.(types.Policy).PolicyName), nil
}

func getPolicyID(instance any) (string, error) {
	return aws.ToString(instance.(types.Policy).PolicyId), nil
}

func getPolicyARN(instance any) (string, error) {
	return aws.ToString(instance.(types.Policy).Arn), nil
}

func getPolicyPath(instance any) (string, error) {
	return aws.ToString(instance.(types.Policy).Path), nil
}

func getPolicyAttachmentCount(instance any) (string, error) {
	p := instance.(types.Policy)
	if p.AttachmentCount != nil {
		return fmt.Sprintf("%d", *p.AttachmentCount), nil
	}
	return "", nil
}

func getPolicyDefaultVersionID(instance any) (string, error) {
	return aws.ToString(instance.(types.Policy).DefaultVersionId), nil
}

func getPolicyCreatedDate(instance any) (string, error) {
	p := instance.(types.Policy)
	if p.CreateDate != nil {
		return p.CreateDate.Format(time.RFC3339), nil
	}
	return "", nil
}

func getPolicyUpdatedDate(instance any) (string, error) {
	p := instance.(types.Policy)
	if p.UpdateDate != nil {
		return p.UpdateDate.Format(time.RFC3339), nil
	}
	return "", nil
}

func getPolicyDescription(instance any) (string, error) {
	return aws.ToString(instance.(types.Policy).Description), nil
}

func getPolicyPermissionsBoundaryCount(instance any) (string, error) {
	p := instance.(types.Policy)
	if p.PermissionsBoundaryUsageCount != nil {
		return fmt.Sprintf("%d", *p.PermissionsBoundaryUsageCount), nil
	}
	return "", nil
}
