package types

import (
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

type LoadBalancerColumnDef struct {
	GetValue func(*types.LoadBalancer) string
}

type TargetGroupColumnDef struct {
	GetValue func(*types.TargetGroup) string
}

type GetTargetGroupsInput struct {
	ListTargetGroupsInput ListTargetGroupsInput
}

type GetLoadBalancersInput struct {
	ListLoadBalancersInput ListLoadBalancersInput
}

type ListTargetGroupsInput struct {
	Names []string
}

type ListLoadBalancersInput struct {
	Names []string
}
