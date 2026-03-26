package awsutil

import (
	"fmt"
	"strings"
)

// ResourceURI represents a parsed protocol-style resource identifier.
// e.g. "rds://my-database" → {Service: "rds", ResourceType: "instance", Resource: "my-database"}
// e.g. "ecs://service/my-cluster/my-svc" → {Service: "ecs", ResourceType: "service", Resource: "my-svc", Params: {"cluster": "my-cluster"}}
type ResourceURI struct {
	Service      string
	ResourceType string
	Resource     string
	Params       map[string]string
}

// resourceTypeConfig defines the path parameters for a resource type.
// PathParams lists named parameters extracted from path segments before the
// final resource identifier.
//
// For example, ECS services require a cluster context:
//
//	PathParams: []string{"cluster"}
//	ecs://service/my-cluster/my-svc → Params["cluster"] = "my-cluster", Resource = "my-svc"
type resourceTypeConfig struct {
	PathParams []string
}

// serviceConfig defines the default resource type and known sub-types for a service.
type serviceConfig struct {
	DefaultType   string
	ResourceTypes map[string]resourceTypeConfig
}

var services = map[string]serviceConfig{
	"ec2": {
		DefaultType: "instance",
		ResourceTypes: map[string]resourceTypeConfig{
			"instance": {},
			"volume":   {},
			"snapshot": {},
			"image":    {},
		},
	},
	"rds": {
		DefaultType: "instance",
		ResourceTypes: map[string]resourceTypeConfig{
			"instance": {},
			"cluster":  {},
		},
	},
	"cf": {
		DefaultType: "stack",
		ResourceTypes: map[string]resourceTypeConfig{
			"stack": {},
		},
	},
	"elasticache": {
		DefaultType: "cluster",
		ResourceTypes: map[string]resourceTypeConfig{
			"cluster": {},
		},
	},
	"elb": {
		DefaultType: "load-balancer",
		ResourceTypes: map[string]resourceTypeConfig{
			"load-balancer": {},
		},
	},
	"vpc": {
		DefaultType: "nat-gateway",
		ResourceTypes: map[string]resourceTypeConfig{
			"nat-gateway": {},
		},
	},
	"ecs": {
		DefaultType: "service",
		ResourceTypes: map[string]resourceTypeConfig{
			"service": {PathParams: []string{"cluster"}},
			"task":    {PathParams: []string{"cluster"}},
		},
	},
}

// ParseResourceURI parses a protocol-style resource identifier.
//
// Supported formats:
//   - "ec2://i-xxx"                        → EC2 instance
//   - "ec2://volume/vol-xxx"               → EC2 volume
//   - "rds://my-database"                  → RDS instance (default type)
//   - "rds://cluster/my-cluster"           → RDS cluster
//   - "cf://my-stack"                      → CloudFormation stack
//   - "elasticache://my-cluster"           → ElastiCache cluster
//   - "elb://my-lb"                        → ELB load balancer
//   - "vpc://nat-gateway/nat-xxx"          → VPC NAT gateway
//   - "ecs://service/my-cluster/my-svc"    → ECS service (cluster extracted as param)
//   - "ecs://task/my-cluster/task-id"      → ECS task (cluster extracted as param)
//   - "i-xxx"                              → EC2 instance (detected by prefix)
//   - "nat-xxx"                            → VPC NAT gateway (detected by prefix)
//
// Returns an error if the input cannot be resolved to a service.
func ParseResourceURI(input string) (*ResourceURI, error) {
	if scheme, rest, ok := strings.Cut(input, "://"); ok {
		if rest == "" {
			return nil, fmt.Errorf("empty resource in URI: %s", input)
		}
		svcConfig, exists := services[scheme]
		if !exists {
			return nil, fmt.Errorf("unknown service: %s", scheme)
		}

		// Check if the first path segment is a known resource type
		resourceType := svcConfig.DefaultType
		resource := rest
		if segment, remainder, ok := strings.Cut(rest, "/"); ok {
			if _, isKnownType := svcConfig.ResourceTypes[segment]; isKnownType {
				resourceType = segment
				resource = remainder
			}
		}

		if resource == "" {
			return nil, fmt.Errorf("empty resource in URI: %s", input)
		}

		// Extract path parameters defined for this resource type
		params := make(map[string]string)
		typeConfig := svcConfig.ResourceTypes[resourceType]
		for _, paramName := range typeConfig.PathParams {
			segment, remainder, ok := strings.Cut(resource, "/")
			if !ok {
				return nil, fmt.Errorf("missing %s in URI: %s (expected %s://%s/%s/<resource>)",
					paramName, input, scheme, resourceType, paramName)
			}
			params[paramName] = segment
			resource = remainder
		}

		if resource == "" {
			return nil, fmt.Errorf("empty resource in URI: %s", input)
		}

		return &ResourceURI{
			Service:      scheme,
			ResourceType: resourceType,
			Resource:     resource,
			Params:       params,
		}, nil
	}

	// Fall back to prefix-based detection
	if info := IdentifyResource(input); info != nil {
		return &ResourceURI{
			Service:      info.Service,
			ResourceType: info.ResourceType,
			Resource:     input,
			Params:       map[string]string{},
		}, nil
	}

	return nil, fmt.Errorf("cannot determine service for %q: use protocol syntax (e.g. rds://my-database)", input)
}
