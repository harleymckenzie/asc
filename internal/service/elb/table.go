package elb

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// Attribute is a struct that defines a field in a detailed table.
type Attribute struct {
	GetValue func(*types.LoadBalancer) string
}

type LoadBalancerAttribute struct {
	GetValue func(*types.LoadBalancerAttribute) string
}

type TargetGroupAttribute struct {
	GetValue func(*types.TargetGroup) string
}

func GetAttributeValue(fieldID string, instance any) (string, error) {
	lb, ok := instance.(types.LoadBalancer)
	if !ok {
		return "", fmt.Errorf("instance is not a types.LoadBalancer")
	}
	attr, ok := availableAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q", fieldID)
	}
	return attr.GetValue(&lb), nil
}

func availableAttributes() map[string]Attribute {
	return map[string]Attribute{
		"Name": {
			GetValue: func(i *types.LoadBalancer) string {
				return aws.ToString(i.LoadBalancerName)
			},
		},
		"DNS Name": {
			GetValue: func(i *types.LoadBalancer) string {
				return aws.ToString(i.DNSName)
			},
		},
		"Scheme": {
			GetValue: func(i *types.LoadBalancer) string {
				return string(i.Scheme)
			},
		},
		"State": {
			GetValue: func(i *types.LoadBalancer) string {
				return format.Status(string(i.State.Code))
			},
		},
		"Type": {
			GetValue: func(i *types.LoadBalancer) string {
				return string(i.Type)
			},
		},
		"IP Type": {
			GetValue: func(i *types.LoadBalancer) string {
				return string(i.IpAddressType)
			},
		},
		"Hosted Zone": {
			GetValue: func(i *types.LoadBalancer) string {
				return aws.ToString(i.CanonicalHostedZoneId)
			},
		},
		"VPC ID": {
			GetValue: func(i *types.LoadBalancer) string {
				return aws.ToString(i.VpcId)
			},
		},
		"Subnets": {
			GetValue: func(i *types.LoadBalancer) string {
				subnets := []string{}
				for _, az := range i.AvailabilityZones {
					subnets = append(subnets, aws.ToString(az.SubnetId))
				}
				return strings.Join(subnets, ", ")
			},
		},
		"Created Time": {
			GetValue: func(i *types.LoadBalancer) string {
				return i.CreatedTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"ARN": {
			GetValue: func(i *types.LoadBalancer) string {
				return aws.ToString(i.LoadBalancerArn)
			},
		},
		"Availability Zones": {
			GetValue: func(i *types.LoadBalancer) string {
				azs := []string{}
				for _, az := range i.AvailabilityZones {
					azs = append(azs, aws.ToString(az.ZoneName))
				}
				return strings.Join(azs, ", ")
			},
		},
	}
}

// GetLoadBalancerAttributeAttributeValue gets the value of a load balancer "attribute".
func GetLoadBalancerAttributeAttributeValue(fieldID string, instance any) (string, error) {
	attributes, ok := instance.([]types.LoadBalancerAttribute)
	if !ok {
		return "", fmt.Errorf("instance is not a types.LoadBalancerAttribute")
	}
	attr, ok := availableLoadBalancerAttributeAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q", fieldID)
	}
	return attr.GetValue(&attributes[0]), nil
}

func availableLoadBalancerAttributeAttributes() map[string]LoadBalancerAttribute {
	return map[string]LoadBalancerAttribute{
		"Name": {
			GetValue: func(i *types.LoadBalancerAttribute) string {
				return aws.ToString(i.Key)
			},
		},
		"Value": {
			GetValue: func(i *types.LoadBalancerAttribute) string {
				return aws.ToString(i.Value)
			},
		},
	}
}

func GetTargetGroupAttributeValue(fieldID string, instance any) (string, error) {
	tg, ok := instance.(types.TargetGroup)
	if !ok {
		return "", fmt.Errorf("instance is not a types.TargetGroup")
	}
	attr, ok := availableTargetGroupAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q", fieldID)
	}
	return attr.GetValue(&tg), nil
}

func availableTargetGroupAttributes() map[string]TargetGroupAttribute {
	return map[string]TargetGroupAttribute{
		"Name": {
			GetValue: func(i *types.TargetGroup) string {
				return aws.ToString(i.TargetGroupName)
			},
		},
		"Target Type": {
			GetValue: func(i *types.TargetGroup) string {
				return string(i.TargetType)
			},
		},
		"Port": {
			GetValue: func(i *types.TargetGroup) string {
				return strconv.Itoa(int(*i.Port))
			},
		},
		"Protocol": {
			GetValue: func(i *types.TargetGroup) string {
				return string(i.Protocol)
			},
		},
		"ARN": {
			GetValue: func(i *types.TargetGroup) string {
				return aws.ToString(i.TargetGroupArn)
			},
		},
		"VPC ID": {
			GetValue: func(i *types.TargetGroup) string {
				return aws.ToString(i.VpcId)
			},
		},
		"Load Balancer": {
			GetValue: func(i *types.TargetGroup) string {
				return getTargetGroupLoadBalancer(*i)
			},
		},
		"Health Check Enabled": {
			GetValue: func(i *types.TargetGroup) string {
				if i.HealthCheckEnabled != nil {
					return strconv.FormatBool(*i.HealthCheckEnabled)
				}
				return "N/A"
			},
		},
		"Health Check Path": {
			GetValue: func(i *types.TargetGroup) string {
				return aws.ToString(i.HealthCheckPath)
			},
		},
		"Health Check Port": {
			GetValue: func(i *types.TargetGroup) string {
				if i.HealthCheckPort != nil {
					return aws.ToString(i.HealthCheckPort)
				}
				return "N/A"
			},
		},
		"Health Check Protocol": {
			GetValue: func(i *types.TargetGroup) string {
				return string(i.HealthCheckProtocol)
			},
		},
		"Health Check Timeout": {
			GetValue: func(i *types.TargetGroup) string {
				return strconv.Itoa(int(*i.HealthCheckTimeoutSeconds))
			},
		},
		"Healthy Threshold": {
			GetValue: func(i *types.TargetGroup) string {
				return strconv.Itoa(int(*i.HealthyThresholdCount))
			},
		},
		"Unhealthy Threshold": {
			GetValue: func(i *types.TargetGroup) string {
				return strconv.Itoa(int(*i.UnhealthyThresholdCount))
			},
		},
	}
}
