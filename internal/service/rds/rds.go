package rds

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	ascTypes "github.com/harleymckenzie/asc/internal/service/rds/types"

	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)

// RDSClientAPI is the interface for the RDS client.
type RDSClientAPI interface {
	DescribeDBInstances(context.Context, *rds.DescribeDBInstancesInput, ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBClusters(context.Context, *rds.DescribeDBClustersInput, ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
	DescribeDBSnapshots(context.Context, *rds.DescribeDBSnapshotsInput, ...func(*rds.Options)) (*rds.DescribeDBSnapshotsOutput, error)
	DescribeDBClusterSnapshots(context.Context, *rds.DescribeDBClusterSnapshotsInput, ...func(*rds.Options)) (*rds.DescribeDBClusterSnapshotsOutput, error)
	ModifyDBInstance(context.Context, *rds.ModifyDBInstanceInput, ...func(*rds.Options)) (*rds.ModifyDBInstanceOutput, error)
	CreateDBSnapshot(context.Context, *rds.CreateDBSnapshotInput, ...func(*rds.Options)) (*rds.CreateDBSnapshotOutput, error)
	CreateDBClusterSnapshot(context.Context, *rds.CreateDBClusterSnapshotInput, ...func(*rds.Options)) (*rds.CreateDBClusterSnapshotOutput, error)
	DeleteDBSnapshot(context.Context, *rds.DeleteDBSnapshotInput, ...func(*rds.Options)) (*rds.DeleteDBSnapshotOutput, error)
	DeleteDBClusterSnapshot(context.Context, *rds.DeleteDBClusterSnapshotInput, ...func(*rds.Options)) (*rds.DeleteDBClusterSnapshotOutput, error)
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

// CreateSnapshot creates a snapshot of an RDS instance or cluster.
func (svc *RDSService) CreateSnapshot(ctx context.Context, input *ascTypes.CreateSnapshotInput) error {
	if input.IsCluster {
		_, err := svc.Client.CreateDBClusterSnapshot(ctx, &rds.CreateDBClusterSnapshotInput{
			DBClusterIdentifier:         &input.Identifier,
			DBClusterSnapshotIdentifier: &input.SnapshotIdentifier,
		})
		return err
	}

	_, err := svc.Client.CreateDBSnapshot(ctx, &rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: &input.Identifier,
		DBSnapshotIdentifier: &input.SnapshotIdentifier,
	})
	return err
}

// GetSnapshots gets RDS DB instance snapshots, optionally filtered by source instance
// or a specific snapshot identifier.
func (svc *RDSService) GetSnapshots(ctx context.Context, input *ascTypes.GetSnapshotsInput) ([]types.DBSnapshot, error) {
	params := &rds.DescribeDBSnapshotsInput{}
	if input.Identifier != "" {
		params.DBInstanceIdentifier = &input.Identifier
	}
	if input.SnapshotIdentifier != "" {
		params.DBSnapshotIdentifier = &input.SnapshotIdentifier
	}

	output, err := svc.Client.DescribeDBSnapshots(ctx, params)
	if err != nil {
		return nil, err
	}

	var snapshots []types.DBSnapshot
	snapshots = append(snapshots, output.DBSnapshots...)
	return snapshots, nil
}

// GetClusterSnapshots gets RDS DB cluster snapshots, optionally filtered by source
// cluster or a specific snapshot identifier.
func (svc *RDSService) GetClusterSnapshots(ctx context.Context, input *ascTypes.GetSnapshotsInput) ([]types.DBClusterSnapshot, error) {
	params := &rds.DescribeDBClusterSnapshotsInput{}
	if input.Identifier != "" {
		params.DBClusterIdentifier = &input.Identifier
	}
	if input.SnapshotIdentifier != "" {
		params.DBClusterSnapshotIdentifier = &input.SnapshotIdentifier
	}

	output, err := svc.Client.DescribeDBClusterSnapshots(ctx, params)
	if err != nil {
		return nil, err
	}

	var snapshots []types.DBClusterSnapshot
	snapshots = append(snapshots, output.DBClusterSnapshots...)
	return snapshots, nil
}

// DeleteSnapshot deletes an RDS DB instance or cluster snapshot.
func (svc *RDSService) DeleteSnapshot(ctx context.Context, input *ascTypes.DeleteSnapshotInput) error {
	if input.IsCluster {
		_, err := svc.Client.DeleteDBClusterSnapshot(ctx, &rds.DeleteDBClusterSnapshotInput{
			DBClusterSnapshotIdentifier: &input.SnapshotIdentifier,
		})
		return err
	}

	_, err := svc.Client.DeleteDBSnapshot(ctx, &rds.DeleteDBSnapshotInput{
		DBSnapshotIdentifier: &input.SnapshotIdentifier,
	})
	return err
}

// WaitForSnapshot waits for an RDS snapshot to become available.
func (svc *RDSService) WaitForSnapshot(ctx context.Context, input *ascTypes.CreateSnapshotInput, maxWait time.Duration) error {
	if input.IsCluster {
		waiter := rds.NewDBClusterSnapshotAvailableWaiter(svc.Client)
		err := waiter.Wait(ctx, &rds.DescribeDBClusterSnapshotsInput{
			DBClusterSnapshotIdentifier: &input.SnapshotIdentifier,
		}, maxWait)
		if err != nil {
			return fmt.Errorf("wait for cluster snapshot: %w", err)
		}
		return nil
	}

	waiter := rds.NewDBSnapshotAvailableWaiter(svc.Client)
	err := waiter.Wait(ctx, &rds.DescribeDBSnapshotsInput{
		DBSnapshotIdentifier: &input.SnapshotIdentifier,
	}, maxWait)
	if err != nil {
		return fmt.Errorf("wait for snapshot: %w", err)
	}
	return nil
}
