package cloudformation

import (
	"context"
	"fmt"
	"strings"

	cfsdk "github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

// IsTerminalStackStatus returns true if the CloudFormation stack status is
// stable (not in progress).
func IsTerminalStackStatus(status string) bool {
	return !strings.HasSuffix(strings.ToUpper(status), "_IN_PROGRESS")
}

// GetStackStatus returns the current status of a CloudFormation stack.
func (svc *CloudFormationService) GetStackStatus(ctx context.Context, stackName string) (string, error) {
	output, err := svc.Client.DescribeStacks(ctx, &cfsdk.DescribeStacksInput{
		StackName: &stackName,
	})
	if err != nil {
		return "", fmt.Errorf("describe stack: %w", err)
	}
	if len(output.Stacks) == 0 {
		return "", fmt.Errorf("stack %s not found", stackName)
	}
	return string(output.Stacks[0].StackStatus), nil
}
