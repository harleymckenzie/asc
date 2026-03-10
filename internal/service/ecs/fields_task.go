package ecs

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

var taskFieldValueGetters = map[string]FieldValueGetter{
	"Service":         getTaskService,
	"Task ID":         getTaskID,
	"ARN":             getTaskARN,
	"Status":          getTaskStatus,
	"Desired Status":  getTaskDesiredStatus,
	"Task Definition": getTaskTaskDefinition,
	"Created At":      getTaskCreatedAt,
	"Started At":      getTaskStartedAt,
	"Stopped At":      getTaskStoppedAt,
	"Stop Code":       getTaskStopCode,
	"Stopped Reason":  getTaskStoppedReason,
	"vCPU":            getTaskCPU,
	"Memory":          getTaskMemory,
	"Launch Type":     getTaskLaunchType,
	"Cluster":         getTaskCluster,
	"Group":           getTaskGroup,
	"Connectivity":    getTaskConnectivity,
	"Health Status":   getTaskHealthStatus,
	"Platform Version": getTaskPlatformVersion,
	"Containers":      getTaskContainers,
}

func getTaskFieldValue(fieldName string, task types.Task) (string, error) {
	if getter, exists := taskFieldValueGetters[fieldName]; exists {
		return getter(task)
	}
	return "", fmt.Errorf("field %s not found in taskFieldValueGetters", fieldName)
}

func getTaskID(instance any) (string, error) {
	arn := aws.ToString(instance.(types.Task).TaskArn)
	return ShortARN(arn), nil
}

func getTaskARN(instance any) (string, error) {
	return aws.ToString(instance.(types.Task).TaskArn), nil
}

func getTaskStatus(instance any) (string, error) {
	return format.Status(aws.ToString(instance.(types.Task).LastStatus)), nil
}

func getTaskDesiredStatus(instance any) (string, error) {
	return format.Status(aws.ToString(instance.(types.Task).DesiredStatus)), nil
}

func getTaskTaskDefinition(instance any) (string, error) {
	arn := aws.ToString(instance.(types.Task).TaskDefinitionArn)
	return ShortARN(arn), nil
}

func getTaskCreatedAt(instance any) (string, error) {
	return format.TimeToStringOrEmpty(instance.(types.Task).CreatedAt), nil
}

func getTaskStartedAt(instance any) (string, error) {
	return format.TimeToStringOrEmpty(instance.(types.Task).StartedAt), nil
}

func getTaskStoppedAt(instance any) (string, error) {
	return format.TimeToStringOrEmpty(instance.(types.Task).StoppedAt), nil
}

func getTaskStopCode(instance any) (string, error) {
	return string(instance.(types.Task).StopCode), nil
}

func getTaskStoppedReason(instance any) (string, error) {
	return aws.ToString(instance.(types.Task).StoppedReason), nil
}

func getTaskCPU(instance any) (string, error) {
	return aws.ToString(instance.(types.Task).Cpu), nil
}

func getTaskMemory(instance any) (string, error) {
	return aws.ToString(instance.(types.Task).Memory), nil
}

func getTaskLaunchType(instance any) (string, error) {
	return string(instance.(types.Task).LaunchType), nil
}

func getTaskCluster(instance any) (string, error) {
	arn := aws.ToString(instance.(types.Task).ClusterArn)
	return ShortARN(arn), nil
}

func getTaskGroup(instance any) (string, error) {
	return aws.ToString(instance.(types.Task).Group), nil
}

func getTaskService(instance any) (string, error) {
	group := aws.ToString(instance.(types.Task).Group)
	if name, ok := strings.CutPrefix(group, "service:"); ok {
		return name, nil
	}
	return "-", nil
}

func getTaskConnectivity(instance any) (string, error) {
	return string(instance.(types.Task).Connectivity), nil
}

func getTaskHealthStatus(instance any) (string, error) {
	return string(instance.(types.Task).HealthStatus), nil
}

func getTaskPlatformVersion(instance any) (string, error) {
	return aws.ToString(instance.(types.Task).PlatformVersion), nil
}

func getTaskContainers(instance any) (string, error) {
	task := instance.(types.Task)
	if len(task.Containers) == 0 {
		return "-", nil
	}
	var names []string
	for _, c := range task.Containers {
		names = append(names, aws.ToString(c.Name))
	}
	return strings.Join(names, ", "), nil
}
