package asg

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/harleymckenzie/asc/pkg/shared/format"
)

// Attribute is a struct that defines a field in a detailed table.
type Attribute struct {
	GetValue func(*types.AutoScalingGroup) string
}

// InstanceAttribute is a struct that defines a field in a detailed table for an instance.
type InstanceAttribute struct {
	GetValue func(*types.Instance) string
}

// ScheduleAttribute is a struct that defines a field in a detailed table for a scheduled update group action.
type ScheduleAttribute struct {
	GetValue func(*types.ScheduledUpdateGroupAction) string
}

// GetAttributeValue returns the value of a field for an Auto Scaling Group
func GetAttributeValue(fieldID string, instance any) string {
	asg, ok := instance.(types.AutoScalingGroup)
	if !ok {
		return ""
	}
	attr := availableAttributes()[fieldID]
	return attr.GetValue(&asg)
}

// availableAttributes returns a map of column definitions for Auto Scaling Groups
func availableAttributes() map[string]Attribute {
	return map[string]Attribute{
		"Name": {
			GetValue: func(i *types.AutoScalingGroup) string {
				return aws.ToString(i.AutoScalingGroupName)
			},
		},
		"Instances": {
			GetValue: func(i *types.AutoScalingGroup) string {
				return strconv.Itoa(len(i.Instances))
			},
		},
		"Desired": {
			GetValue: func(i *types.AutoScalingGroup) string {
				return strconv.Itoa(int(*i.DesiredCapacity))
			},
		},
		"Min": {
			GetValue: func(i *types.AutoScalingGroup) string {
				return strconv.Itoa(int(*i.MinSize))
			},
		},
		"Max": {
			GetValue: func(i *types.AutoScalingGroup) string {
				return strconv.Itoa(int(*i.MaxSize))
			},
		},
		"ARN": {
			GetValue: func(i *types.AutoScalingGroup) string {
				return aws.ToString(i.AutoScalingGroupARN)
			},
		},
	}
}

// GetInstanceAttributeValue returns the value of a field for an instance
func GetInstanceAttributeValue(fieldID string, instance any) string {
	inst, ok := instance.(types.Instance)
	if !ok {
		return ""
	}
	attr := availableInstanceAttributes()[fieldID]
	return attr.GetValue(&inst)
}

// availableInstanceAttributes returns a map of column definitions for instances
func availableInstanceAttributes() map[string]InstanceAttribute {
	return map[string]InstanceAttribute{
		"Name": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.InstanceId)
			},
		},
		"State": {
			GetValue: func(i *types.Instance) string {
				return string(i.LifecycleState)
			},
		},
		"Instance Type": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.InstanceType)
			},
		},
		"Launch Template/Configuration": {
			GetValue: func(i *types.Instance) string {
				if i.LaunchTemplate != nil {
					return aws.ToString(i.LaunchTemplate.LaunchTemplateName)
				}
				return aws.ToString(i.LaunchConfigurationName)
			},
		},
		"Availability Zone": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.AvailabilityZone)
			},
		},
		"Health": {
			GetValue: func(i *types.Instance) string {
				return format.Status(aws.ToString(i.HealthStatus))
			},
		},
	}
}

// GetScheduleAttributeValue returns the value of a field for a scheduled update group action
func GetScheduleAttributeValue(fieldID string, instance any) string {
	sched, ok := instance.(types.ScheduledUpdateGroupAction)
	if !ok {
		return ""
	}
	attr := availableSchedulesAttributes()[fieldID]
	return attr.GetValue(&sched)
}

// availableSchedulesAttributes returns a map of column definitions for scheduled update group actions
func availableSchedulesAttributes() map[string]ScheduleAttribute {
	return map[string]ScheduleAttribute{
		"Auto Scaling Group": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				return aws.ToString(i.AutoScalingGroupName)
			},
		},
		"Name": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				return aws.ToString(i.ScheduledActionName)
			},
		},
		"Recurrence": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				if i.Recurrence == nil {
					return ""
				}
				return aws.ToString(i.Recurrence)
			},
		},
		"Start Time": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				if i.StartTime == nil {
					return ""
				}
				return i.StartTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"End Time": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				if i.EndTime == nil {
					return ""
				}
				return i.EndTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"Desired Capacity": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				if i.DesiredCapacity == nil {
					return ""
				}
				return strconv.Itoa(int(*i.DesiredCapacity))
			},
		},
		"Min": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				if i.MinSize == nil {
					return ""
				}
				return strconv.Itoa(int(*i.MinSize))
			},
		},
		"Max": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				if i.MaxSize == nil {
					return ""
				}
				return strconv.Itoa(int(*i.MaxSize))
			},
		},
	}
}
