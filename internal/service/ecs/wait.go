package ecs

import (
	"context"
	"fmt"
	"strings"

	ecssdk "github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

var serviceTerminalStates = map[string]bool{
	"active":   true,
	"inactive": true,
}

var taskTerminalStates = map[string]bool{
	"running":       true,
	"stopped":       true,
	"deactivating":  true,
}

// IsTerminalServiceState returns true if the ECS service status is stable.
func IsTerminalServiceState(status string) bool {
	return serviceTerminalStates[strings.ToLower(status)]
}

// IsTerminalTaskState returns true if the ECS task status is stable.
func IsTerminalTaskState(status string) bool {
	return taskTerminalStates[strings.ToLower(status)]
}

// GetServiceStatus returns the current status of an ECS service.
func (svc *ECSService) GetServiceStatus(ctx context.Context, cluster, serviceName string) (string, error) {
	output, err := svc.Client.DescribeServices(ctx, &ecssdk.DescribeServicesInput{
		Cluster:  &cluster,
		Services: []string{serviceName},
		Include:  []types.ServiceField{},
	})
	if err != nil {
		return "", fmt.Errorf("describe service: %w", err)
	}
	if len(output.Services) == 0 {
		return "", fmt.Errorf("service %s not found in cluster %s", serviceName, cluster)
	}
	if output.Services[0].Status == nil {
		return "", fmt.Errorf("service %s has no status", serviceName)
	}
	return *output.Services[0].Status, nil
}

// GetTaskStatus returns the current status of an ECS task.
func (svc *ECSService) GetTaskStatus(ctx context.Context, cluster, taskID string) (string, error) {
	output, err := svc.Client.DescribeTasks(ctx, &ecssdk.DescribeTasksInput{
		Cluster: &cluster,
		Tasks:   []string{taskID},
		Include: []types.TaskField{},
	})
	if err != nil {
		return "", fmt.Errorf("describe task: %w", err)
	}
	if len(output.Tasks) == 0 {
		return "", fmt.Errorf("task %s not found in cluster %s", taskID, cluster)
	}
	if output.Tasks[0].LastStatus == nil {
		return "", fmt.Errorf("task %s has no status", taskID)
	}
	return *output.Tasks[0].LastStatus, nil
}
