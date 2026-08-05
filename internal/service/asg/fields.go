package asg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type FieldValueGetter func(instance any) (string, error)

// AutoScaling Group field getters
var asgFieldValueGetters = map[string]FieldValueGetter{
	"Name":                      getASGName,
	"Instances":                 getASGInstances,
	"Desired":                   getASGDesired,
	"Min":                       getASGMin,
	"Max":                       getASGMax,
	"ARN":                       getASGARN,
	"Status":                    getASGStatus,
	"Created Time":              getASGCreatedTime,
	"Default Cooldown":          getASGDefaultCooldown,
	"Health Check Type":         getASGHealthCheckType,
	"Health Check Grace Period": getASGHealthCheckGracePeriod,
	"Availability Zones":        getASGAvailabilityZones,
	"Subnets":                   getASGSubnets,
	"Launch Template":           getASGLaunchTemplate,
	"Launch Configuration":      getASGLaunchConfiguration,
}

// Instance field getters
var instanceFieldValueGetters = map[string]FieldValueGetter{
	"Name":                          getInstanceName,
	"State":                         getInstanceState,
	"Instance Type":                 getInstanceType,
	"Launch Template/Configuration": getInstanceLaunchConfig,
	"Availability Zone":             getInstanceAZ,
	"Health":                        getInstanceHealth,
}

// Schedule field getters
var scheduleFieldValueGetters = map[string]FieldValueGetter{
	"Auto Scaling Group": getScheduleASGName,
	"Name":               getScheduleName,
	"Recurrence":         getScheduleRecurrence,
	"Start Time":         getScheduleStartTime,
	"End Time":           getScheduleEndTime,
	"Desired Capacity":   getScheduleDesiredCapacity,
	"Min":                getScheduleMin,
	"Max":                getScheduleMax,
}

// GetFieldValue returns the value of a field for the given instance.
// This function routes field requests to the appropriate type-specific handler.
func GetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.AutoScalingGroup:
		return getASGFieldValue(fieldName, v)
	case types.Instance:
		return getInstanceFieldValue(fieldName, v)
	case types.ScheduledUpdateGroupAction:
		return getScheduleFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// getASGFieldValue returns the value of a field for an Auto Scaling Group
func getASGFieldValue(fieldName string, asg types.AutoScalingGroup) (string, error) {
	if getter, exists := asgFieldValueGetters[fieldName]; exists {
		return getter(asg)
	}
	return "", fmt.Errorf("field %s not found in asgFieldValueGetters", fieldName)
}

// getInstanceFieldValue returns the value of a field for an ASG instance
func getInstanceFieldValue(fieldName string, instance types.Instance) (string, error) {
	if getter, exists := instanceFieldValueGetters[fieldName]; exists {
		return getter(instance)
	}
	return "", fmt.Errorf("field %s not found in instanceFieldValueGetters", fieldName)
}

// getScheduleFieldValue returns the value of a field for a scheduled action
func getScheduleFieldValue(fieldName string, schedule types.ScheduledUpdateGroupAction) (string, error) {
	if getter, exists := scheduleFieldValueGetters[fieldName]; exists {
		return getter(schedule)
	}
	return "", fmt.Errorf("field %s not found in scheduleFieldValueGetters", fieldName)
}

// GetTagValue returns the value of a tag for the given instance.
// This function handles tag retrieval for Auto Scaling Groups and related resources.
func GetTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.AutoScalingGroup:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.ScheduledUpdateGroupAction:
		// Scheduled actions don't have tags in AWS ASG
		return "", nil
	default:
		return "", fmt.Errorf("unsupported instance type for tags: %T", instance)
	}
	return "", nil
}

// -----------------------------------------------------------------------------
// Auto Scaling Group field getters
// -----------------------------------------------------------------------------

// getASGName returns the name of the Auto Scaling Group
func getASGName(instance any) (string, error) {
	return aws.ToString(instance.(types.AutoScalingGroup).AutoScalingGroupName), nil
}

// getASGInstances returns the number of instances in the Auto Scaling Group
func getASGInstances(instance any) (string, error) {
	return strconv.Itoa(len(instance.(types.AutoScalingGroup).Instances)), nil
}

// getASGDesired returns the desired capacity of the Auto Scaling Group
func getASGDesired(instance any) (string, error) {
	return strconv.Itoa(int(*instance.(types.AutoScalingGroup).DesiredCapacity)), nil
}

// getASGMin returns the minimum size of the Auto Scaling Group
func getASGMin(instance any) (string, error) {
	return strconv.Itoa(int(*instance.(types.AutoScalingGroup).MinSize)), nil
}

// getASGMax returns the maximum size of the Auto Scaling Group
func getASGMax(instance any) (string, error) {
	return strconv.Itoa(int(*instance.(types.AutoScalingGroup).MaxSize)), nil
}

// getASGARN returns the ARN of the Auto Scaling Group
func getASGARN(instance any) (string, error) {
	return aws.ToString(instance.(types.AutoScalingGroup).AutoScalingGroupARN), nil
}

// getASGStatus returns the status of the Auto Scaling Group
func getASGStatus(instance any) (string, error) {
	return aws.ToString(instance.(types.AutoScalingGroup).Status), nil
}

// getASGCreatedTime returns the creation time of the Auto Scaling Group
func getASGCreatedTime(instance any) (string, error) {
	return format.TimeToStringOrEmpty(instance.(types.AutoScalingGroup).CreatedTime), nil
}

// getASGDefaultCooldown returns the default cooldown period (in seconds) of the Auto Scaling Group
func getASGDefaultCooldown(instance any) (string, error) {
	return format.Int32ToStringOrEmpty(instance.(types.AutoScalingGroup).DefaultCooldown), nil
}

// getASGHealthCheckType returns the health check type of the Auto Scaling Group
func getASGHealthCheckType(instance any) (string, error) {
	return aws.ToString(instance.(types.AutoScalingGroup).HealthCheckType), nil
}

// getASGHealthCheckGracePeriod returns the health check grace period (in seconds) of the Auto Scaling Group
func getASGHealthCheckGracePeriod(instance any) (string, error) {
	return format.Int32ToStringOrEmpty(instance.(types.AutoScalingGroup).HealthCheckGracePeriod), nil
}

// getASGAvailabilityZones returns the availability zones of the Auto Scaling Group
func getASGAvailabilityZones(instance any) (string, error) {
	return strings.Join(instance.(types.AutoScalingGroup).AvailabilityZones, "\n"), nil
}

// getASGSubnets returns the subnets (VPC zone identifier) of the Auto Scaling Group
func getASGSubnets(instance any) (string, error) {
	vpcZone := aws.ToString(instance.(types.AutoScalingGroup).VPCZoneIdentifier)
	if vpcZone == "" {
		return "", nil
	}
	return strings.Join(strings.Split(vpcZone, ","), "\n"), nil
}

// getASGLaunchTemplate returns the launch template name of the Auto Scaling Group
func getASGLaunchTemplate(instance any) (string, error) {
	spec := GetLaunchTemplateSpecification(instance.(types.AutoScalingGroup))
	if spec == nil {
		return "", nil
	}
	return aws.ToString(spec.LaunchTemplateName), nil
}

// getASGLaunchConfiguration returns the launch configuration name of the Auto Scaling Group
func getASGLaunchConfiguration(instance any) (string, error) {
	return aws.ToString(instance.(types.AutoScalingGroup).LaunchConfigurationName), nil
}

// GetLaunchTemplateSpecification returns the launch template specification for the Auto
// Scaling Group, resolving it from either a directly attached launch template or a
// mixed instances policy. It returns nil if the group uses a launch configuration.
func GetLaunchTemplateSpecification(group types.AutoScalingGroup) *types.LaunchTemplateSpecification {
	if group.LaunchTemplate != nil {
		return group.LaunchTemplate
	}
	if group.MixedInstancesPolicy != nil && group.MixedInstancesPolicy.LaunchTemplate != nil {
		return group.MixedInstancesPolicy.LaunchTemplate.LaunchTemplateSpecification
	}
	return nil
}

// -----------------------------------------------------------------------------
// Instance field getters
// -----------------------------------------------------------------------------

// getInstanceName returns the instance ID
func getInstanceName(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).InstanceId), nil
}

// getInstanceState returns the lifecycle state of the instance
func getInstanceState(instance any) (string, error) {
	return string(instance.(types.Instance).LifecycleState), nil
}

// getInstanceType returns the instance type
func getInstanceType(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).InstanceType), nil
}

// getInstanceLaunchConfig returns the launch template or configuration name
func getInstanceLaunchConfig(instance any) (string, error) {
	inst := instance.(types.Instance)
	if inst.LaunchTemplate != nil {
		return aws.ToString(inst.LaunchTemplate.LaunchTemplateName), nil
	}
	return aws.ToString(inst.LaunchConfigurationName), nil
}

// getInstanceAZ returns the availability zone of the instance
func getInstanceAZ(instance any) (string, error) {
	return aws.ToString(instance.(types.Instance).AvailabilityZone), nil
}

// getInstanceHealth returns the health status of the instance
func getInstanceHealth(instance any) (string, error) {
	return format.Status(aws.ToString(instance.(types.Instance).HealthStatus)), nil
}

// -----------------------------------------------------------------------------
// Schedule field getters
// -----------------------------------------------------------------------------

// getScheduleASGName returns the Auto Scaling Group name for the scheduled action
func getScheduleASGName(instance any) (string, error) {
	return aws.ToString(instance.(types.ScheduledUpdateGroupAction).AutoScalingGroupName), nil
}

// getScheduleName returns the name of the scheduled action
func getScheduleName(instance any) (string, error) {
	return aws.ToString(instance.(types.ScheduledUpdateGroupAction).ScheduledActionName), nil
}

// getScheduleRecurrence returns the recurrence pattern for the scheduled action
func getScheduleRecurrence(instance any) (string, error) {
	sched := instance.(types.ScheduledUpdateGroupAction)
	if sched.Recurrence == nil {
		return "", nil
	}
	return aws.ToString(sched.Recurrence), nil
}

// getScheduleStartTime returns the start time of the scheduled action
func getScheduleStartTime(instance any) (string, error) {
	sched := instance.(types.ScheduledUpdateGroupAction)
	if sched.StartTime == nil {
		return "", nil
	}
	return sched.StartTime.Local().Format("2006-01-02 15:04:05 MST"), nil
}

// getScheduleEndTime returns the end time of the scheduled action
func getScheduleEndTime(instance any) (string, error) {
	sched := instance.(types.ScheduledUpdateGroupAction)
	if sched.EndTime == nil {
		return "", nil
	}
	return sched.EndTime.Local().Format("2006-01-02 15:04:05 MST"), nil
}

// getScheduleDesiredCapacity returns the desired capacity for the scheduled action
func getScheduleDesiredCapacity(instance any) (string, error) {
	sched := instance.(types.ScheduledUpdateGroupAction)
	if sched.DesiredCapacity == nil {
		return "", nil
	}
	return strconv.Itoa(int(*sched.DesiredCapacity)), nil
}

// getScheduleMin returns the minimum size for the scheduled action
func getScheduleMin(instance any) (string, error) {
	sched := instance.(types.ScheduledUpdateGroupAction)
	if sched.MinSize == nil {
		return "", nil
	}
	return strconv.Itoa(int(*sched.MinSize)), nil
}

// getScheduleMax returns the maximum size for the scheduled action
func getScheduleMax(instance any) (string, error) {
	sched := instance.(types.ScheduledUpdateGroupAction)
	if sched.MaxSize == nil {
		return "", nil
	}
	return strconv.Itoa(int(*sched.MaxSize)), nil
}
