package rds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	ascTypes "github.com/harleymckenzie/asc/internal/service/rds/types"

	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)

// RDSClientAPI is the interface for the RDS client.
type RDSClientAPI interface {
	DescribeDBInstances(context.Context, *rds.DescribeDBInstancesInput, ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBClusters(context.Context, *rds.DescribeDBClustersInput, ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
	ModifyDBInstance(context.Context, *rds.ModifyDBInstanceInput, ...func(*rds.Options)) (*rds.ModifyDBInstanceOutput, error)
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
func (svc *RDSService) GetInstances(ctx context.Context, input *ascTypes.GetInstancesInput) ([]types.DBInstance, error) {
	output, err := svc.Client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: &input.InstanceIdentifier,
	})
	if err != nil {
		return nil, err
	}

	var instances []types.DBInstance
	instances = append(instances, output.DBInstances...)
	return instances, nil
}

// GetClusters gets all the RDS clusters.
func (svc *RDSService) GetClusters(ctx context.Context, input *ascTypes.GetClustersInput) ([]types.DBCluster, error) {
	clusterOutput, err := svc.Client.DescribeDBClusters(ctx, &rds.DescribeDBClustersInput{
		DBClusterIdentifier: &input.ClusterIdentifier,
	})
	if err != nil {
		return nil, err
	}

	var clusters []types.DBCluster
	clusters = append(clusters, clusterOutput.DBClusters...)
	return clusters, nil
}

// ModifyInstance modifies an RDS instance.
func (svc *RDSService) ModifyInstance(ctx context.Context, input *ascTypes.ModifyInstanceInput) error {
	_, err := svc.Client.ModifyDBInstance(ctx, &rds.ModifyDBInstanceInput{
		DBInstanceIdentifier: input.DBInstanceIdentifier,
		ApplyImmediately: input.ApplyImmediately,
		DBInstanceClass: input.DBInstanceClass,
		PreferredMaintenanceWindow: input.PreferredMaintenanceWindow,
	})
	return err
}
