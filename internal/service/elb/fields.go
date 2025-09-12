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
// This function routes field requests to the appropriate type-specific handler.
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
// ELB tags are handled differently - they need to be fetched separately
// This function signature matches the existing pattern but ELB tags
// are not directly accessible from the LoadBalancer/TargetGroup struct
func GetTagValue(tagKey string, instance any) (string, error) {
	// ELB tags are handled differently - they need to be fetched separately
	// This function signature matches the existing pattern but ELB tags
	// are not directly accessible from the LoadBalancer/TargetGroup struct
	return "", nil
}

// -----------------------------------------------------------------------------
// Load Balancer field getters
// -----------------------------------------------------------------------------

// getLoadBalancerName returns the name of the load balancer
func getLoadBalancerName(instance any) (string, error) {
	return aws.ToString(instance.(types.LoadBalancer).LoadBalancerName), nil
}

// getLoadBalancerDNSName returns the DNS name of the load balancer
func getLoadBalancerDNSName(instance any) (string, error) {
	return aws.ToString(instance.(types.LoadBalancer).DNSName), nil
}

// getLoadBalancerScheme returns the scheme of the load balancer (internet-facing or internal)
func getLoadBalancerScheme(instance any) (string, error) {
	return string(instance.(types.LoadBalancer).Scheme), nil
}

// getLoadBalancerType returns the type of the load balancer (application, network, gateway)
func getLoadBalancerType(instance any) (string, error) {
	return string(instance.(types.LoadBalancer).Type), nil
}

// getLoadBalancerState returns the current state of the load balancer
func getLoadBalancerState(instance any) (string, error) {
	lb := instance.(types.LoadBalancer)
	if lb.State == nil {
		return "", nil
	}
	return format.Status(string(lb.State.Code)), nil
}

// getLoadBalancerVPCID returns the VPC ID where the load balancer is deployed
func getLoadBalancerVPCID(instance any) (string, error) {
	return aws.ToString(instance.(types.LoadBalancer).VpcId), nil
}

// getLoadBalancerIPAddressType returns the IP address type of the load balancer
func getLoadBalancerIPAddressType(instance any) (string, error) {
	return string(instance.(types.LoadBalancer).IpAddressType), nil
}

// getLoadBalancerSecurityGroups returns the security groups associated with the load balancer
func getLoadBalancerSecurityGroups(instance any) (string, error) {
	lb := instance.(types.LoadBalancer)
	return strings.Join(lb.SecurityGroups, ", "), nil
}

// getLoadBalancerSubnets returns the subnets where the load balancer is deployed
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

// getLoadBalancerAvailabilityZones returns the availability zones where the load balancer is deployed
func getLoadBalancerAvailabilityZones(instance any) (string, error) {
	lb := instance.(types.LoadBalancer)
	var zones []string
	for _, zone := range lb.AvailabilityZones {
		zones = append(zones, aws.ToString(zone.ZoneName))
	}
	return strings.Join(zones, ", "), nil
}

// getLoadBalancerCreated returns the creation time of the load balancer
func getLoadBalancerCreated(instance any) (string, error) {
	return format.TimeToStringOrEmpty(instance.(types.LoadBalancer).CreatedTime), nil
}

// getLoadBalancerARN returns the ARN of the load balancer
func getLoadBalancerARN(instance any) (string, error) {
	return aws.ToString(instance.(types.LoadBalancer).LoadBalancerArn), nil
}

// -----------------------------------------------------------------------------
// Target Group field getters
// -----------------------------------------------------------------------------

// getTargetGroupName returns the name of the target group
func getTargetGroupName(instance any) (string, error) {
	return aws.ToString(instance.(types.TargetGroup).TargetGroupName), nil
}

// getTargetGroupARN returns the ARN of the target group
func getTargetGroupARN(instance any) (string, error) {
	return aws.ToString(instance.(types.TargetGroup).TargetGroupArn), nil
}

// getTargetGroupProtocol returns the protocol used by the target group
func getTargetGroupProtocol(instance any) (string, error) {
	return string(instance.(types.TargetGroup).Protocol), nil
}

// getTargetGroupPort returns the port used by the target group
func getTargetGroupPort(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if tg.Port == nil {
		return "", nil
	}
	return strconv.Itoa(int(*tg.Port)), nil
}

// getTargetGroupVPCID returns the VPC ID where the target group is deployed
func getTargetGroupVPCID(instance any) (string, error) {
	return aws.ToString(instance.(types.TargetGroup).VpcId), nil
}

// getTargetGroupTargetType returns the target type of the target group
func getTargetGroupTargetType(instance any) (string, error) {
	return string(instance.(types.TargetGroup).TargetType), nil
}

// getTargetGroupLoadBalancerField returns the associated load balancer name
func getTargetGroupLoadBalancerField(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if len(tg.LoadBalancerArns) > 0 {
		name := strings.Split(tg.LoadBalancerArns[0], "/")
		return name[len(name)-2], nil
	}
	return "", nil
}

// getTargetGroupHealthCheckProtocol returns the health check protocol
func getTargetGroupHealthCheckProtocol(instance any) (string, error) {
	return string(instance.(types.TargetGroup).HealthCheckProtocol), nil
}

// getTargetGroupHealthCheckPort returns the health check port
func getTargetGroupHealthCheckPort(instance any) (string, error) {
	return aws.ToString(instance.(types.TargetGroup).HealthCheckPort), nil
}

// getTargetGroupHealthCheckPath returns the health check path
func getTargetGroupHealthCheckPath(instance any) (string, error) {
	return aws.ToString(instance.(types.TargetGroup).HealthCheckPath), nil
}

// getTargetGroupHealthCheckInterval returns the health check interval in seconds
func getTargetGroupHealthCheckInterval(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if tg.HealthCheckIntervalSeconds == nil {
		return "", nil
	}
	return strconv.Itoa(int(*tg.HealthCheckIntervalSeconds)), nil
}

// getTargetGroupHealthCheckTimeout returns the health check timeout in seconds
func getTargetGroupHealthCheckTimeout(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if tg.HealthCheckTimeoutSeconds == nil {
		return "", nil
	}
	return strconv.Itoa(int(*tg.HealthCheckTimeoutSeconds)), nil
}

// getTargetGroupHealthyThreshold returns the healthy threshold count
func getTargetGroupHealthyThreshold(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if tg.HealthyThresholdCount == nil {
		return "", nil
	}
	return strconv.Itoa(int(*tg.HealthyThresholdCount)), nil
}

// getTargetGroupUnhealthyThreshold returns the unhealthy threshold count
func getTargetGroupUnhealthyThreshold(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if tg.UnhealthyThresholdCount == nil {
		return "", nil
	}
	return strconv.Itoa(int(*tg.UnhealthyThresholdCount)), nil
}

// getTargetGroupHTTPCode returns the HTTP code used for health checks
func getTargetGroupHTTPCode(instance any) (string, error) {
	return aws.ToString(instance.(types.TargetGroup).Matcher.HttpCode), nil
}

// getTargetGroupHealthCheckEnabled returns whether health checks are enabled
func getTargetGroupHealthCheckEnabled(instance any) (string, error) {
	tg := instance.(types.TargetGroup)
	if tg.HealthCheckEnabled == nil {
		return "", nil
	}
	return format.BoolToLabel(tg.HealthCheckEnabled, "Yes", "No"), nil
}

// getTargetGroupHealthCheckGracePeriod returns the health check grace period
func getTargetGroupHealthCheckGracePeriod(instance any) (string, error) {
	// HealthCheckGracePeriodSeconds is not available in the TargetGroup type
	// This field may be specific to certain target group types or require separate API calls
	return "-", nil
}
