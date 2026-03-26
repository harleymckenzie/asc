package elb

import (
	"context"
	"fmt"
	"strings"

	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

var elbTerminalStates = map[string]bool{
	"active":          true,
	"active_impaired": true,
	"failed":          true,
}

// IsTerminalLoadBalancerState returns true if the load balancer state is stable.
func IsTerminalLoadBalancerState(status string) bool {
	return elbTerminalStates[strings.ToLower(status)]
}

// GetLoadBalancerStatus returns the current state of a load balancer.
func (svc *ELBService) GetLoadBalancerStatus(ctx context.Context, lbName string) (string, error) {
	output, err := svc.Client.DescribeLoadBalancers(ctx, &elbv2.DescribeLoadBalancersInput{
		Names: []string{lbName},
	})
	if err != nil {
		return "", fmt.Errorf("describe load balancer: %w", err)
	}
	if len(output.LoadBalancers) == 0 {
		return "", fmt.Errorf("load balancer %s not found", lbName)
	}
	lb := output.LoadBalancers[0]
	if lb.State == nil {
		return "", fmt.Errorf("load balancer %s has no state", lbName)
	}
	return string(lb.State.Code), nil
}
