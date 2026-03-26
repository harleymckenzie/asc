package wait

import (
	"context"
	"fmt"
	"time"

	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// NewWaitCmd creates the top-level wait command.
func NewWaitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wait <protocol://resource>",
		Short: "Wait for an AWS resource to reach a stable state",
		Long: `Wait for an AWS resource to reach a stable state (e.g. available, running, stopped).

Supports protocol-style URIs to identify the service and resource type:
  ec2://i-xxx                            EC2 instance
  ec2://volume/vol-xxx                   EC2 volume
  ec2://snapshot/snap-xxx                EC2 snapshot
  ec2://image/ami-xxx                    EC2 AMI
  rds://my-database                      RDS instance
  rds://cluster/my-cluster               RDS cluster
  cf://my-stack                          CloudFormation stack
  elasticache://my-cluster               ElastiCache cluster
  elb://my-lb                            ELB load balancer
  vpc://nat-gateway/nat-xxx              VPC NAT gateway
  ecs://service/my-cluster/my-service    ECS service
  ecs://task/my-cluster/task-id          ECS task

Resources with known ID prefixes can omit the protocol:
  i-xxx, vol-xxx, snap-xxx, ami-xxx, nat-xxx`,
		Example: `  asc wait ec2://i-1234567890abcdef0
  asc wait rds://my-database
  asc wait cf://my-stack
  asc wait ecs://service/my-cluster/my-service
  asc wait i-1234567890abcdef0
  asc wait nat-1234567890abcdef0`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdutil.DefaultErrorHandler(runWait(cmd, args))
		},
	}
	return cmd
}

func runWait(cmd *cobra.Command, args []string) error {
	uri, err := awsutil.ParseResourceURI(args[0])
	if err != nil {
		return err
	}

	profile, region := cmdutil.GetPersistentFlags(cmd)
	return ExecuteWait(cmd.Context(), profile, region, uri)
}

// ExecuteWait is the shared wait implementation used by both the top-level
// `asc wait` command and per-service `asc <service> wait` commands.
func ExecuteWait(ctx context.Context, profile, region string, uri *awsutil.ResourceURI) error {
	handler, err := getHandler(uri)
	if err != nil {
		return err
	}

	statusFunc, isTerminal, err := handler(ctx, profile, region, uri)
	if err != nil {
		return fmt.Errorf("create %s service: %w", uri.Service, err)
	}

	fmt.Printf("Waiting for %s %s %s to reach a stable state...\n", uri.Service, uri.ResourceType, uri.Resource)

	finalStatus, err := awsutil.WaitForStatus(ctx, awsutil.WaitConfig{
		ResourceName: uri.Resource,
		PollInterval: 10 * time.Second,
		MaxWait:      30 * time.Minute,
		StatusFunc:   statusFunc,
		IsTerminal:   isTerminal,
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s %s %s reached state: %s\n", uri.Service, uri.ResourceType, uri.Resource, finalStatus)
	return nil
}
