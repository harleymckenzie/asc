package iam

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

var userFieldValueGetters = map[string]FieldValueGetter{
	"User Name":          getUserName,
	"User ID":            getUserID,
	"ARN":                getUserARN,
	"Path":               getUserPath,
	"Created Date":       getUserCreatedDate,
	"Password Last Used": getUserPasswordLastUsed,
}

func getUserName(instance any) (string, error) {
	return aws.ToString(instance.(types.User).UserName), nil
}

func getUserID(instance any) (string, error) {
	return aws.ToString(instance.(types.User).UserId), nil
}

func getUserARN(instance any) (string, error) {
	return aws.ToString(instance.(types.User).Arn), nil
}

func getUserPath(instance any) (string, error) {
	return aws.ToString(instance.(types.User).Path), nil
}

func getUserCreatedDate(instance any) (string, error) {
	u := instance.(types.User)
	if u.CreateDate != nil {
		return u.CreateDate.Format(time.RFC3339), nil
	}
	return "", nil
}

func getUserPasswordLastUsed(instance any) (string, error) {
	u := instance.(types.User)
	if u.PasswordLastUsed != nil {
		return u.PasswordLastUsed.Format(time.RFC3339), nil
	}
	return "", nil
}
