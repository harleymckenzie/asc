package elb

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type FieldValueGetter func(instance any) (string, error)

// Load Balancer field getters
var loadBalancerFieldValueGetters = map[string]FieldValueGetter{
	"Name":               getLoadBalancerName,
	"DNS Name":           getLoadBalancerDNSName,
	"Scheme":             getLoadBalancerScheme,
	"Type":               getLoadBalancerType,
	"State":              getLoadBalancerState,
	"VPC ID":             getLoadBalancerVPCID,
	"IP Type":            getLoadBalancerIPAddressType,
	"Security Groups":    getLoadBalancerSecurityGroups,
	"Subnets":            getLoadBalancerSubnets,
	"Availability Zones": getLoadBalancerAvailabilityZones,
	"Created Time":       getLoadBalancerCreated,
	"ARN":                getLoadBalancerARN,
}

// Target Group field getters
var targetGroupFieldValueGetters = map[string]FieldValueGetter{
	"Name":                      getTargetGroupName,
	"ARN":                       getTargetGroupARN,
	"Protocol":                  getTargetGroupProtocol,
	"Port":                      getTargetGroupPort,
	"VPC ID":                    getTargetGroupVPCID,
	"Target Type":               getTargetGroupTargetType,
	"Load Balancer":             getTargetGroupLoadBalancerField,
	"Health Check Protocol":     getTargetGroupHealthCheckProtocol,
	"Health Check Port":         getTargetGroupHealthCheckPort,
	"Health Check Path":         getTargetGroupHealthCheckPath,
	"Health Check Interval":     getTargetGroupHealthCheckInterval,
	"Health Check Timeout":      getTargetGroupHealthCheckTimeout,
	"Healthy Threshold":         getTargetGroupHealthyThreshold,
	"Unhealthy Threshold":       getTargetGroupUnhealthyThreshold,
	"HTTP Code":                 getTargetGroupHTTPCode,
	"Health Check Enabled":      getTargetGroupHealthCheckEnabled,
	"Health Check Grace Period": getTargetGroupHealthCheckGracePeriod,
}

// GetFieldValue returns the value of a field for the given instance.
func GetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.LoadBalancer:
		return getLoadBalancerFieldValue(fieldName, v)
	case types.TargetGroup:
		return getTargetGroupFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// getLoadBalancerFieldValue returns the value of a field for a Load Balancer
func getLoadBalancerFieldValue(fieldName string, lb types.LoadBalancer) (string, error) {
	if getter, exists := loadBalancerFieldValueGetters[fieldName]; exists {
		return getter(lb)
	}
	return "", fmt.Errorf("field %s not found in loadBalancerFieldValueGetters", fieldName)
}

// getTargetGroupFieldValue returns the value of a field for a Target Group
func getTargetGroupFieldValue(fieldName string, tg types.TargetGroup) (string, error) {
	if getter, exists := targetGroupFieldValueGetters[fieldName]; exists {
		return getter(tg)
	}
	return "", fmt.Errorf("field %s not found in targetGroupFieldValueGetters", fieldName)
}

// GetTagValue returns the value of a tag for the given instance.
func GetTagValue(tagKey string, instance any) (string, error) {
	// ELB tags are handled differently - they need to be fetched separately
	// This function signature matches the existing pattern but ELB tags
	// are not directly accessible from the LoadBalancer/TargetGroup struct
	return "", nil
}

// -----------------------------------------------------------------------------
// Load Balancer field getters
// -----------------------------------------------------------------------------

func getLoadBalancerName(instance any) (string, error) {
	return aws.ToString(instance.(types.LoadBalancer).LoadBalancerName), nil
}

func getLoadBalancerDNSName(instance any) (string, error) {
	return aws.ToString(instance.(types.LoadBalancer).DNSName), nil
}

func getLoadBalancerScheme(instance any) (string, error) {
	return string(instance.(types.LoadBalancer).Scheme), nil
}

func getLoadBalancerType(instance any) (string, error) {
	return string(instance.(types.LoadBalancer).Type), nil
}

func getLoadBalancerState(instance any) (string, error) {
	lb := instance.(types.LoadBalancer)
	if lb.State == nil {
		return "", nil
	}
	return format.Status(string(lb.State.Code)), nil
}

func getLoadBalancerVPCID(instance any) (string, error) {
	return aws.ToString(instance.(types.LoadBalancer).VpcId), nil
}

func getLoadBalancerIPAddressType(instance any) (string, error) {
	return string(instance.(types.LoadBalancer).IpAddressType), nil
}

func getLoadBalancerSecurityGroups(instance any) (string, error) {
	lb := instance.(types.LoadBalancer)
	return strings.Join(lb.SecurityGroups, ", "), nil
}

func getLoadBalancerSubnets(instance any) (string, error) {
	lb := instance.(types.LoadBalancer)
	var subnets []string
	for _, az := range lb.AvailabilityZones {
		if az.SubnetId != nil {
			subnets = append(subnets, aws.ToString(az.SubnetId))
		}
	}
	return strings.Join(subnets, ", "), nil
}

func getLoadBalancerAvailabilityZones(instance any) (string, error) {
	lb := instance.(types.LoadBalancer)
	var zones []string
	for _, zone := range lb.AvailabilityZones {
		zones = append(zones, aws.ToString(zone.ZoneName))
	}
	return strings.Join(zones, ", "), nil
}

func getLoadBalancerCreated(instance any) (string, error) {
	return format.TimeToStringOrEmpty(instance.(types.LoadBalancer).CreatedTime), nil
}

func getLoadBalancerARN(instance any) (string, error) {
	return aws.ToString(instance.(types.LoadBalancer).LoadBalancerArn), nil
}

// -----------------------------------------------------------------------------
// Target Group field getters
// -----------------------------------------------------------------------------

func getTargetGroupName(instance any) (string, error) {
	return aws.ToString(instance.(types.TargetGroup).TargetGroupName), nil
}

func getTargetGroupARN(instance any) (string, error) {
	return aws.ToString(instance.(types.TargetGroup).TargetGroupArn), nil
}

func getTargetGroupProtocol(instance any) (string, error) {
	return string(instance.(types.TargetGroup).Protocol), nil
}

func getTargetGroupPort(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if tg.Port == nil {
		return "", nil
	}
	return strconv.Itoa(int(*tg.Port)), nil
}

func getTargetGroupVPCID(instance any) (string, error) {
	return aws.ToString(instance.(types.TargetGroup).VpcId), nil
}

func getTargetGroupTargetType(instance any) (string, error) {
	return string(instance.(types.TargetGroup).TargetType), nil
}

func getTargetGroupLoadBalancerField(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if len(tg.LoadBalancerArns) > 0 {
		name := strings.Split(tg.LoadBalancerArns[0], "/")
		return name[len(name)-2], nil
	}
	return "", nil
}

func getTargetGroupHealthCheckProtocol(instance any) (string, error) {
	return string(instance.(types.TargetGroup).HealthCheckProtocol), nil
}

func getTargetGroupHealthCheckPort(instance any) (string, error) {
	return aws.ToString(instance.(types.TargetGroup).HealthCheckPort), nil
}

func getTargetGroupHealthCheckPath(instance any) (string, error) {
	return aws.ToString(instance.(types.TargetGroup).HealthCheckPath), nil
}

func getTargetGroupHealthCheckInterval(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if tg.HealthCheckIntervalSeconds == nil {
		return "", nil
	}
	return strconv.Itoa(int(*tg.HealthCheckIntervalSeconds)), nil
}

func getTargetGroupHealthCheckTimeout(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if tg.HealthCheckTimeoutSeconds == nil {
		return "", nil
	}
	return strconv.Itoa(int(*tg.HealthCheckTimeoutSeconds)), nil
}

func getTargetGroupHealthyThreshold(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if tg.HealthyThresholdCount == nil {
		return "", nil
	}
	return strconv.Itoa(int(*tg.HealthyThresholdCount)), nil
}

func getTargetGroupUnhealthyThreshold(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if tg.UnhealthyThresholdCount == nil {
		return "", nil
	}
	return strconv.Itoa(int(*tg.UnhealthyThresholdCount)), nil
}

func getTargetGroupHTTPCode(instance any) (string, error) {
	return aws.ToString(instance.(types.TargetGroup).Matcher.HttpCode), nil
}

func getTargetGroupHealthCheckEnabled(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if tg.HealthCheckEnabled == nil {
		return "", nil
	}
	return format.BoolToLabel(tg.HealthCheckEnabled, "Yes", "No"), nil
}

func getTargetGroupHealthCheckGracePeriod(instance any) (string, error) {
	// HealthCheckGracePeriodSeconds is not available in the TargetGroup type
	// This field may be specific to certain target group types or require separate API calls
	return "-", nil
}
