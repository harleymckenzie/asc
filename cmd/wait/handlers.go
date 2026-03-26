package wait

import (
	"context"

	"github.com/harleymckenzie/asc/internal/service/cloudformation"
	"github.com/harleymckenzie/asc/internal/service/ec2"
	"github.com/harleymckenzie/asc/internal/service/ecs"
	"github.com/harleymckenzie/asc/internal/service/elasticache"
	"github.com/harleymckenzie/asc/internal/service/elb"
	"github.com/harleymckenzie/asc/internal/service/rds"
	"github.com/harleymckenzie/asc/internal/service/vpc"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)

func init() {
	// EC2
	RegisterHandler("ec2/instance", newEC2Handler(
		func(svc *ec2.EC2Service, uri *awsutil.ResourceURI) func(ctx context.Context) (string, error) {
			return func(ctx context.Context) (string, error) { return svc.GetInstanceStatus(ctx, uri.Resource) }
		}, ec2.IsTerminalInstanceState,
	))
	RegisterHandler("ec2/volume", newEC2Handler(
		func(svc *ec2.EC2Service, uri *awsutil.ResourceURI) func(ctx context.Context) (string, error) {
			return func(ctx context.Context) (string, error) { return svc.GetVolumeStatus(ctx, uri.Resource) }
		}, ec2.IsTerminalVolumeState,
	))
	RegisterHandler("ec2/snapshot", newEC2Handler(
		func(svc *ec2.EC2Service, uri *awsutil.ResourceURI) func(ctx context.Context) (string, error) {
			return func(ctx context.Context) (string, error) { return svc.GetSnapshotStatus(ctx, uri.Resource) }
		}, ec2.IsTerminalSnapshotState,
	))
	RegisterHandler("ec2/image", newEC2Handler(
		func(svc *ec2.EC2Service, uri *awsutil.ResourceURI) func(ctx context.Context) (string, error) {
			return func(ctx context.Context) (string, error) { return svc.GetImageStatus(ctx, uri.Resource) }
		}, ec2.IsTerminalImageState,
	))

	// RDS
	RegisterHandler("rds/instance", newRDSHandler(
		func(svc *rds.RDSService, uri *awsutil.ResourceURI) func(ctx context.Context) (string, error) {
			return func(ctx context.Context) (string, error) { return svc.GetInstanceStatus(ctx, uri.Resource) }
		}, rds.IsTerminalInstanceState,
	))
	RegisterHandler("rds/cluster", newRDSHandler(
		func(svc *rds.RDSService, uri *awsutil.ResourceURI) func(ctx context.Context) (string, error) {
			return func(ctx context.Context) (string, error) { return svc.GetClusterStatus(ctx, uri.Resource) }
		}, rds.IsTerminalClusterState,
	))

	// CloudFormation
	RegisterHandler("cf/stack", func(ctx context.Context, profile, region string, uri *awsutil.ResourceURI) (func(ctx context.Context) (string, error), func(string) bool, error) {
		svc, err := cloudformation.NewCloudFormationService(ctx, profile, region)
		if err != nil {
			return nil, nil, err
		}
		return func(ctx context.Context) (string, error) {
			return svc.GetStackStatus(ctx, uri.Resource)
		}, cloudformation.IsTerminalStackStatus, nil
	})

	// ElastiCache
	RegisterHandler("elasticache/cluster", func(ctx context.Context, profile, region string, uri *awsutil.ResourceURI) (func(ctx context.Context) (string, error), func(string) bool, error) {
		svc, err := elasticache.NewElasticacheService(ctx, profile, region)
		if err != nil {
			return nil, nil, err
		}
		return func(ctx context.Context) (string, error) {
			return svc.GetClusterStatus(ctx, uri.Resource)
		}, elasticache.IsTerminalClusterState, nil
	})

	// ELB
	RegisterHandler("elb/load-balancer", func(ctx context.Context, profile, region string, uri *awsutil.ResourceURI) (func(ctx context.Context) (string, error), func(string) bool, error) {
		svc, err := elb.NewELBService(ctx, profile, region)
		if err != nil {
			return nil, nil, err
		}
		return func(ctx context.Context) (string, error) {
			return svc.GetLoadBalancerStatus(ctx, uri.Resource)
		}, elb.IsTerminalLoadBalancerState, nil
	})

	// VPC
	RegisterHandler("vpc/nat-gateway", func(ctx context.Context, profile, region string, uri *awsutil.ResourceURI) (func(ctx context.Context) (string, error), func(string) bool, error) {
		svc, err := vpc.NewVPCService(ctx, profile, region)
		if err != nil {
			return nil, nil, err
		}
		return func(ctx context.Context) (string, error) {
			return svc.GetNatGatewayStatus(ctx, uri.Resource)
		}, vpc.IsTerminalNatGatewayState, nil
	})

	// ECS
	RegisterHandler("ecs/service", newECSHandler(
		func(svc *ecs.ECSService, cluster string, uri *awsutil.ResourceURI) func(ctx context.Context) (string, error) {
			return func(ctx context.Context) (string, error) {
				return svc.GetServiceStatus(ctx, cluster, uri.Resource)
			}
		}, ecs.IsTerminalServiceState,
	))
	RegisterHandler("ecs/task", newECSHandler(
		func(svc *ecs.ECSService, cluster string, uri *awsutil.ResourceURI) func(ctx context.Context) (string, error) {
			return func(ctx context.Context) (string, error) {
				return svc.GetTaskStatus(ctx, cluster, uri.Resource)
			}
		}, ecs.IsTerminalTaskState,
	))
}

// Helper constructors to reduce boilerplate for services with multiple resource types.

func newEC2Handler(
	makeStatusFunc func(*ec2.EC2Service, *awsutil.ResourceURI) func(ctx context.Context) (string, error),
	isTerminal func(string) bool,
) WaitHandler {
	return func(ctx context.Context, profile, region string, uri *awsutil.ResourceURI) (func(ctx context.Context) (string, error), func(string) bool, error) {
		svc, err := ec2.NewEC2Service(ctx, profile, region)
		if err != nil {
			return nil, nil, err
		}
		return makeStatusFunc(svc, uri), isTerminal, nil
	}
}

func newRDSHandler(
	makeStatusFunc func(*rds.RDSService, *awsutil.ResourceURI) func(ctx context.Context) (string, error),
	isTerminal func(string) bool,
) WaitHandler {
	return func(ctx context.Context, profile, region string, uri *awsutil.ResourceURI) (func(ctx context.Context) (string, error), func(string) bool, error) {
		svc, err := rds.NewRDSService(ctx, profile, region)
		if err != nil {
			return nil, nil, err
		}
		return makeStatusFunc(svc, uri), isTerminal, nil
	}
}

func newECSHandler(
	makeStatusFunc func(*ecs.ECSService, string, *awsutil.ResourceURI) func(ctx context.Context) (string, error),
	isTerminal func(string) bool,
) WaitHandler {
	return func(ctx context.Context, profile, region string, uri *awsutil.ResourceURI) (func(ctx context.Context) (string, error), func(string) bool, error) {
		svc, err := ecs.NewECSService(ctx, profile, region)
		if err != nil {
			return nil, nil, err
		}
		cluster := uri.Params["cluster"]
		return makeStatusFunc(svc, cluster, uri), isTerminal, nil
	}
}
