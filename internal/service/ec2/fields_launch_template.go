package ec2

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

var launchTemplateFieldValueGetters = map[string]FieldValueGetter{
	"Launch Template Name": getLaunchTemplateName,
	"Launch Template ID":   getLaunchTemplateID,
	"Version":              getLaunchTemplateVersion,
	"Default Version":      getLaunchTemplateDefaultVersion,
	"Version Description":  getLaunchTemplateVersionDescription,
	"Created Time":         getLaunchTemplateCreatedTime,
	"Created By":           getLaunchTemplateCreatedBy,
	"Instance Type":        getLaunchTemplateInstanceType,
	"AMI ID":               getLaunchTemplateAMIID,
	"Key Name":             getLaunchTemplateKeyName,
	"EBS Optimized":        getLaunchTemplateEBSOptimized,
	"IAM Instance Profile": getLaunchTemplateIAMInstanceProfile,
	"Monitoring":           getLaunchTemplateMonitoring,
	"Security Group IDs":   getLaunchTemplateSecurityGroupIDs,
	"Security Groups":      getLaunchTemplateSecurityGroups,
}

// getLaunchTemplateFieldValue returns the value of a field for a launch template version
func getLaunchTemplateFieldValue(fieldName string, version types.LaunchTemplateVersion) (string, error) {
	if getter, exists := launchTemplateFieldValueGetters[fieldName]; exists {
		value, err := getter(version)
		if err != nil {
			return "", fmt.Errorf("failed to get field value for %s: %w", fieldName, err)
		}
		return value, nil
	}
	return "", fmt.Errorf("field %s not found in launch template fieldValueGetters", fieldName)
}

// Version metadata getters

func getLaunchTemplateName(instance any) (string, error) {
	return aws.ToString(instance.(types.LaunchTemplateVersion).LaunchTemplateName), nil
}

func getLaunchTemplateID(instance any) (string, error) {
	return aws.ToString(instance.(types.LaunchTemplateVersion).LaunchTemplateId), nil
}

func getLaunchTemplateVersion(instance any) (string, error) {
	return format.Int64ToStringOrEmpty(instance.(types.LaunchTemplateVersion).VersionNumber), nil
}

func getLaunchTemplateDefaultVersion(instance any) (string, error) {
	return format.BoolToLabel(instance.(types.LaunchTemplateVersion).DefaultVersion, "Yes", "No"), nil
}

func getLaunchTemplateVersionDescription(instance any) (string, error) {
	return aws.ToString(instance.(types.LaunchTemplateVersion).VersionDescription), nil
}

func getLaunchTemplateCreatedTime(instance any) (string, error) {
	return format.TimeToStringOrEmpty(instance.(types.LaunchTemplateVersion).CreateTime), nil
}

func getLaunchTemplateCreatedBy(instance any) (string, error) {
	return aws.ToString(instance.(types.LaunchTemplateVersion).CreatedBy), nil
}

// Launch template data getters

func getLaunchTemplateInstanceType(instance any) (string, error) {
	data := instance.(types.LaunchTemplateVersion).LaunchTemplateData
	if data == nil {
		return "", nil
	}
	return string(data.InstanceType), nil
}

func getLaunchTemplateAMIID(instance any) (string, error) {
	data := instance.(types.LaunchTemplateVersion).LaunchTemplateData
	if data == nil {
		return "", nil
	}
	return aws.ToString(data.ImageId), nil
}

func getLaunchTemplateKeyName(instance any) (string, error) {
	data := instance.(types.LaunchTemplateVersion).LaunchTemplateData
	if data == nil {
		return "", nil
	}
	return aws.ToString(data.KeyName), nil
}

func getLaunchTemplateEBSOptimized(instance any) (string, error) {
	data := instance.(types.LaunchTemplateVersion).LaunchTemplateData
	if data == nil || data.EbsOptimized == nil {
		return "", nil
	}
	return format.BoolToLabel(data.EbsOptimized, "Yes", "No"), nil
}

func getLaunchTemplateIAMInstanceProfile(instance any) (string, error) {
	data := instance.(types.LaunchTemplateVersion).LaunchTemplateData
	if data == nil || data.IamInstanceProfile == nil {
		return "", nil
	}
	if arn := aws.ToString(data.IamInstanceProfile.Arn); arn != "" {
		return arn, nil
	}
	return aws.ToString(data.IamInstanceProfile.Name), nil
}

func getLaunchTemplateMonitoring(instance any) (string, error) {
	data := instance.(types.LaunchTemplateVersion).LaunchTemplateData
	if data == nil || data.Monitoring == nil {
		return "", nil
	}
	return format.BoolToLabel(data.Monitoring.Enabled, "Enabled", "Disabled"), nil
}

func getLaunchTemplateSecurityGroupIDs(instance any) (string, error) {
	data := instance.(types.LaunchTemplateVersion).LaunchTemplateData
	if data == nil {
		return "", nil
	}
	return strings.Join(data.SecurityGroupIds, "\n"), nil
}

func getLaunchTemplateSecurityGroups(instance any) (string, error) {
	data := instance.(types.LaunchTemplateVersion).LaunchTemplateData
	if data == nil {
		return "", nil
	}
	return strings.Join(data.SecurityGroups, "\n"), nil
}
