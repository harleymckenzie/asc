package rds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	"github.com/harleymckenzie/asc/pkg/shared/awsutil"
)

// RDSClientAPI is the interface for the RDS client.
type RDSClientAPI interface {
	DescribeDBInstances(context.Context, *rds.DescribeDBInstancesInput, ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBClusters(context.Context, *rds.DescribeDBClustersInput, ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
}

// RDSService is the service for the RDS client.
type RDSService struct {
	Client RDSClientAPI
}

// NewRDSService creates a new RDS service.
func NewRDSService(ctx context.Context, profile string, region string) (*RDSService, error) {
	cfg, err := awsutil.LoadDefaultConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	client := rds.NewFromConfig(cfg.Config)
	return &RDSService{Client: client}, nil
}

// GetInstances gets all the RDS instances.
func (svc *RDSService) GetInstances(ctx context.Context) ([]types.DBInstance, error) {
	output, err := svc.Client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, err
	}

	var instances []types.DBInstance
	instances = append(instances, output.DBInstances...)
	return instances, nil
}

// GetClusters gets all the RDS clusters.
func (svc *RDSService) GetClusters(ctx context.Context) ([]types.DBCluster, error) {
	clusterOutput, err := svc.Client.DescribeDBClusters(ctx, &rds.DescribeDBClustersInput{})
	if err != nil {
		return nil, err
	}

	var clusters []types.DBCluster
	clusters = append(clusters, clusterOutput.DBClusters...)
	return clusters, nil
}

// getDBInstanceRole gets the role of the RDS instance.
func getDBInstanceRole(instance types.DBInstance, clusters []types.DBCluster) string {
	// If ReadReplicaSourceDBInstanceIdentifier is set, then this is a replica. If
	// if ReadReplicaDBInstanceIdentifiers is set, then this is a primary.
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
