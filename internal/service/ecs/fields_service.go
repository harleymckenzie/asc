package ecs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

var serviceFieldValueGetters = map[string]FieldValueGetter{
	"Name":              getServiceName,
	"ARN":               getServiceARN,
	"Status":            getServiceStatus,
	"Launch Type":       getServiceLaunchType,
	"Task Definition":   getServiceTaskDefinition,
	"Desired Count":     getServiceDesiredCount,
	"Running Count":     getServiceRunningCount,
	"Pending Count":     getServicePendingCount,
	"Cluster":           getServiceCluster,
	"Created Date":      getServiceCreatedDate,
	"Platform Version":  getServicePlatformVersion,
	"Scheduling":        getServiceSchedulingStrategy,
	"Deployment Config": getServiceDeploymentConfig,
	"Network Mode":      getServiceNetworkMode,
	"Load Balancers":    getServiceLoadBalancers,
	"Role ARN":          getServiceRoleARN,
	"Subnets":           getServiceSubnets,
	"Security Groups":   getServiceSecurityGroups,
	"Public IP":         getServiceAssignPublicIP,
}

func getServiceFieldValue(fieldName string, service types.Service) (string, error) {
	if getter, exists := serviceFieldValueGetters[fieldName]; exists {
		return getter(service)
	}
	return "", fmt.Errorf("field %s not found in serviceFieldValueGetters", fieldName)
}

func getServiceName(instance any) (string, error) {
	return aws.ToString(instance.(types.Service).ServiceName), nil
}

func getServiceARN(instance any) (string, error) {
	return aws.ToString(instance.(types.Service).ServiceArn), nil
}

func getServiceStatus(instance any) (string, error) {
	return format.Status(aws.ToString(instance.(types.Service).Status)), nil
}

func getServiceLaunchType(instance any) (string, error) {
	return string(instance.(types.Service).LaunchType), nil
}

func getServiceTaskDefinition(instance any) (string, error) {
	arn := aws.ToString(instance.(types.Service).TaskDefinition)
	return ShortARN(arn), nil
}

func getServiceDesiredCount(instance any) (string, error) {
	return strconv.Itoa(int(instance.(types.Service).DesiredCount)), nil
}

func getServiceRunningCount(instance any) (string, error) {
	return strconv.Itoa(int(instance.(types.Service).RunningCount)), nil
}

func getServicePendingCount(instance any) (string, error) {
	return strconv.Itoa(int(instance.(types.Service).PendingCount)), nil
}

func getServiceCluster(instance any) (string, error) {
	arn := aws.ToString(instance.(types.Service).ClusterArn)
	return ShortARN(arn), nil
}

func getServiceCreatedDate(instance any) (string, error) {
	return format.TimeToStringOrEmpty(instance.(types.Service).CreatedAt), nil
}

func getServicePlatformVersion(instance any) (string, error) {
	return aws.ToString(instance.(types.Service).PlatformVersion), nil
}

func getServiceSchedulingStrategy(instance any) (string, error) {
	return string(instance.(types.Service).SchedulingStrategy), nil
}

func getServiceDeploymentConfig(instance any) (string, error) {
	svc := instance.(types.Service)
	if svc.DeploymentConfiguration == nil {
		return "-", nil
	}
	min := format.Int32ToStringOrDefault(svc.DeploymentConfiguration.MinimumHealthyPercent, "-")
	max := format.Int32ToStringOrDefault(svc.DeploymentConfiguration.MaximumPercent, "-")
	return fmt.Sprintf("min %s%% / max %s%%", min, max), nil
}

func getServiceNetworkMode(instance any) (string, error) {
	svc := instance.(types.Service)
	if svc.NetworkConfiguration == nil || svc.NetworkConfiguration.AwsvpcConfiguration == nil {
		return "-", nil
	}
	return "awsvpc", nil
}

func getServiceLoadBalancers(instance any) (string, error) {
	svc := instance.(types.Service)
	if len(svc.LoadBalancers) == 0 {
		return "-", nil
	}
	var lbs []string
	for _, lb := range svc.LoadBalancers {
		tg := aws.ToString(lb.TargetGroupArn)
		lbs = append(lbs, ShortARN(tg))
	}
	return strings.Join(lbs, ", "), nil
}

func getServiceRoleARN(instance any) (string, error) {
	return format.StringOrDefault(instance.(types.Service).RoleArn, "-"), nil
}

func getServiceSubnets(instance any) (string, error) {
	svc := instance.(types.Service)
	if svc.NetworkConfiguration == nil || svc.NetworkConfiguration.AwsvpcConfiguration == nil {
		return "-", nil
	}
	return strings.Join(svc.NetworkConfiguration.AwsvpcConfiguration.Subnets, ", "), nil
}

func getServiceSecurityGroups(instance any) (string, error) {
	svc := instance.(types.Service)
	if svc.NetworkConfiguration == nil || svc.NetworkConfiguration.AwsvpcConfiguration == nil {
		return "-", nil
	}
	return strings.Join(svc.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups, ", "), nil
}

func getServiceAssignPublicIP(instance any) (string, error) {
	svc := instance.(types.Service)
	if svc.NetworkConfiguration == nil || svc.NetworkConfiguration.AwsvpcConfiguration == nil {
		return "-", nil
	}
	return string(svc.NetworkConfiguration.AwsvpcConfiguration.AssignPublicIp), nil
}
