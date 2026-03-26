package vpc

import (
	"context"
	"fmt"
	"strings"

	ec2sdk "github.com/aws/aws-sdk-go-v2/service/ec2"
)

var natGatewayTerminalStates = map[string]bool{
	"available": true,
	"failed":    true,
	"deleted":   true,
	"deleting":  true,
}

// IsTerminalNatGatewayState returns true if the NAT gateway state is stable.
func IsTerminalNatGatewayState(status string) bool {
	return natGatewayTerminalStates[strings.ToLower(status)]
}

// GetNatGatewayStatus returns the current state of a NAT gateway.
func (svc *VPCService) GetNatGatewayStatus(ctx context.Context, natGatewayID string) (string, error) {
	output, err := svc.Client.DescribeNatGateways(ctx, &ec2sdk.DescribeNatGatewaysInput{
		NatGatewayIds: []string{natGatewayID},
	})
	if err != nil {
		return "", fmt.Errorf("describe NAT gateway: %w", err)
	}
	if len(output.NatGateways) == 0 {
		return "", fmt.Errorf("NAT gateway %s not found", natGatewayID)
	}
	return string(output.NatGateways[0].State), nil
}
