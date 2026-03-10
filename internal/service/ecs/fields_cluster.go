package ecs

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

var clusterFieldValueGetters = map[string]FieldValueGetter{
	"Name":                getClusterName,
	"ARN":                 getClusterARN,
	"Status":              getClusterStatus,
	"Active Services":     getClusterActiveServices,
	"Running Tasks":       getClusterRunningTasks,
	"Pending Tasks":       getClusterPendingTasks,
	"Registered Instances": getClusterRegisteredInstances,
	"Capacity Providers":  getClusterCapacityProviders,
	"Default Strategy":    getClusterDefaultCapacityProviderStrategy,
}

func getClusterFieldValue(fieldName string, cluster types.Cluster) (string, error) {
	if getter, exists := clusterFieldValueGetters[fieldName]; exists {
		return getter(cluster)
	}
	return "", fmt.Errorf("field %s not found in clusterFieldValueGetters", fieldName)
}

func getClusterName(instance any) (string, error) {
	return aws.ToString(instance.(types.Cluster).ClusterName), nil
}

func getClusterARN(instance any) (string, error) {
	return aws.ToString(instance.(types.Cluster).ClusterArn), nil
}

func getClusterStatus(instance any) (string, error) {
	return format.Status(aws.ToString(instance.(types.Cluster).Status)), nil
}

func getClusterActiveServices(instance any) (string, error) {
	cluster := instance.(types.Cluster)
	for _, stat := range cluster.Statistics {
		if aws.ToString(stat.Name) == "activeServiceCount" {
			return aws.ToString(stat.Value), nil
		}
	}
	return "0", nil
}

func getClusterRunningTasks(instance any) (string, error) {
	return strconv.Itoa(int(instance.(types.Cluster).RunningTasksCount)), nil
}

func getClusterPendingTasks(instance any) (string, error) {
	return strconv.Itoa(int(instance.(types.Cluster).PendingTasksCount)), nil
}

func getClusterRegisteredInstances(instance any) (string, error) {
	return strconv.Itoa(int(instance.(types.Cluster).RegisteredContainerInstancesCount)), nil
}

func getClusterCapacityProviders(instance any) (string, error) {
	cluster := instance.(types.Cluster)
	if len(cluster.CapacityProviders) == 0 {
		return "-", nil
	}
	result := ""
	for i, cp := range cluster.CapacityProviders {
		if i > 0 {
			result += ", "
		}
		result += cp
	}
	return result, nil
}

func getClusterDefaultCapacityProviderStrategy(instance any) (string, error) {
	cluster := instance.(types.Cluster)
	if len(cluster.DefaultCapacityProviderStrategy) == 0 {
		return "-", nil
	}
	result := ""
	for i, s := range cluster.DefaultCapacityProviderStrategy {
		if i > 0 {
			result += ", "
		}
		result += aws.ToString(s.CapacityProvider)
	}
	return result, nil
}
