package rds

import (
	"context"
	"fmt"
	"strings"

	rdssdk "github.com/aws/aws-sdk-go-v2/service/rds"
)

var rdsInstanceTerminalStates = map[string]bool{
	"available":                           true,
	"stopped":                             true,
	"failed":                              true,
	"deleting":                            true,
	"deleted":                             true,
	"storage-full":                        true,
	"incompatible-credentials":            true,
	"incompatible-parameters":             true,
	"incompatible-restore":                true,
	"inaccessible-encryption-credentials": true,
}

var rdsClusterTerminalStates = map[string]bool{
	"available":                           true,
	"stopped":                             true,
	"failing-over":                        true,
	"inaccessible-encryption-credentials": true,
	"migration-failed":                    true,
}

// IsTerminalInstanceState returns true if the RDS instance status is a stable,
// non-processing state.
func IsTerminalInstanceState(status string) bool {
	return rdsInstanceTerminalStates[strings.ToLower(status)]
}

// IsTerminalClusterState returns true if the RDS cluster status is a stable,
// non-processing state.
func IsTerminalClusterState(status string) bool {
	return rdsClusterTerminalStates[strings.ToLower(status)]
}

// GetInstanceStatus returns the current status of an RDS instance.
func (svc *RDSService) GetInstanceStatus(ctx context.Context, identifier string) (string, error) {
	output, err := svc.Client.DescribeDBInstances(ctx, &rdssdk.DescribeDBInstancesInput{
		DBInstanceIdentifier: &identifier,
	})
	if err != nil {
		return "", fmt.Errorf("describe DB instance: %w", err)
	}
	if len(output.DBInstances) == 0 {
		return "", fmt.Errorf("DB instance %s not found", identifier)
	}
	return *output.DBInstances[0].DBInstanceStatus, nil
}

// GetClusterStatus returns the current status of an RDS cluster.
func (svc *RDSService) GetClusterStatus(ctx context.Context, identifier string) (string, error) {
	output, err := svc.Client.DescribeDBClusters(ctx, &rdssdk.DescribeDBClustersInput{
		DBClusterIdentifier: &identifier,
	})
	if err != nil {
		return "", fmt.Errorf("describe DB cluster: %w", err)
	}
	if len(output.DBClusters) == 0 {
		return "", fmt.Errorf("DB cluster %s not found", identifier)
	}
	return *output.DBClusters[0].Status, nil
}
