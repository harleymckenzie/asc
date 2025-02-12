package rds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	"github.com/harleymckenzie/asc/pkg/service/base"
)

type RDSClientAPI interface {
	DescribeDBInstances(context.Context, *rds.DescribeDBInstancesInput, ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBClusters(context.Context, *rds.DescribeDBClustersInput, ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
	StartDBInstance(context.Context, *rds.StartDBInstanceInput, ...func(*rds.Options)) (*rds.StartDBInstanceOutput, error)
	StopDBInstance(context.Context, *rds.StopDBInstanceInput, ...func(*rds.Options)) (*rds.StopDBInstanceOutput, error)
}

// RDSService is a struct that holds the RDS client
type RDSService struct {
	*base.AWSService
	Client RDSClientAPI
	ctx    context.Context
}

type columnDef struct {
	id       string
	title    string
	getValue func(*types.DBInstance, []types.DBCluster) string
}

var availableColumns = []columnDef{
	{
		id:    "cluster_identifier",
		title: "Cluster Identifier",
		getValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
			if i.DBClusterIdentifier != nil {
				return aws.ToString(i.DBClusterIdentifier)
			}
			return "None"
		},
	},
	{
		id:    "identifier",
		title: "Identifier",
		getValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
			return aws.ToString(i.DBInstanceIdentifier)
		},
	},
	{
		id:    "status",
		title: "Status",
		getValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
			return tableformat.ResourceState(aws.ToString(i.DBInstanceStatus))
		},
	},
	{
		id:    "engine",
		title: "Engine",
		getValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
			return string(*i.Engine)
		},
	},
	{
		id:    "size",
		title: "Size",
		getValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
			return string(*i.DBInstanceClass)
		},
	},
	{
		id:    "role",
		title: "Role",
		getValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
			return getDBInstanceRole(*i, clusters)
		},
	},
	{
		id:    "endpoint",
		title: "Endpoint",
		getValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
			return aws.ToString(i.Endpoint.Address)
		},
	},
}

func NewRDSService(ctx context.Context, profile string, region string) (*RDSService, error) {
	var cfg aws.Config
	var err error

	opts := []func(*config.LoadOptions) error{}

	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}
	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	cfg, err = config.LoadDefaultConfig(ctx, opts...)

	if err != nil {
		return nil, err
	}

	client := rds.NewFromConfig(cfg)
	return &RDSService{
		AWSService: base.NewAWSService(DefaultTableConfig),
		Client:     client,
	}, nil
}

// ListInstances lists all RDS instances
func (svc *RDSService) ListInstances(ctx context.Context, options base.ListOptions) error {
	cmd := NewListInstancesCommand(svc, options)
	return cmd.Execute(ctx)
}

// StartInstances starts the specified RDS instances
func (svc *RDSService) StartInstances(ctx context.Context, options base.StateChangeOptions) error {
	cmd := NewStartInstanceCommand(svc, options)
	return cmd.Execute(ctx)
}

// StopInstances stops the specified RDS instances
func (svc *RDSService) StopInstances(ctx context.Context, options base.StateChangeOptions) error {
	cmd := NewStopInstanceCommand(svc, options)
	return cmd.Execute(ctx)
}

func getDBInstanceRole(instance types.DBInstance, clusters []types.DBCluster) string {
	if instance.ReadReplicaSourceDBInstanceIdentifier != nil {
		return "Replica"
	}

	if len(instance.ReadReplicaDBInstanceIdentifiers) > 0 {
		return "Primary"
	}

	if instance.DBClusterIdentifier == nil {
		return "None"
	}

	for _, cluster := range clusters {
		for _, member := range cluster.DBClusterMembers {
			if aws.ToString(member.DBInstanceIdentifier) == aws.ToString(instance.DBInstanceIdentifier) {
				if member.IsClusterWriter != nil && *member.IsClusterWriter {
					return "Writer"
				}
				return "Reader"
			}
		}
	}

	return "Unknown"
}
