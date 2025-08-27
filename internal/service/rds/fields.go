package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type FieldValueGetter func(instance any) (string, error)

// RDS DB Instance field getters
var dbInstanceFieldValueGetters = map[string]FieldValueGetter{
	"Cluster Identifier": getDBInstanceClusterID,
	"Identifier":         getDBInstanceID,
	"Status":             getDBInstanceStatus,
	"Engine":             getDBInstanceEngine,
	"Engine Version":     getDBInstanceEngineVersion,
	"Size":               getDBInstanceSize,
	"Role":               getDBInstanceRoleField,
	"Endpoint":           getDBInstanceEndpoint,
}

// Clusters context - stored globally to support role calculation
var clustersContext []types.DBCluster

// GetFieldValue returns the value of a field for the given instance.
func GetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.DBInstance:
		return getDBInstanceFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// SetClustersContext sets the clusters context for role calculation
func SetClustersContext(clusters []types.DBCluster) {
	clustersContext = clusters
}

// getDBInstanceFieldValue returns the value of a field for an RDS DB instance
func getDBInstanceFieldValue(fieldName string, instance types.DBInstance) (string, error) {
	if getter, exists := dbInstanceFieldValueGetters[fieldName]; exists {
		return getter(instance)
	}
	return "", fmt.Errorf("field %s not found in dbInstanceFieldValueGetters", fieldName)
}

// GetTagValue returns the value of a tag for the given instance.
// Note: This duplicates the function in table.go but will replace it when table.go is removed
func GetTagValueNew(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.DBInstance:
		for _, tag := range v.TagList {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	default:
		return "", fmt.Errorf("unsupported instance type for tags: %T", instance)
	}
	return "", nil
}

// -----------------------------------------------------------------------------
// DB Instance field getters
// -----------------------------------------------------------------------------

func getDBInstanceClusterID(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.DBClusterIdentifier != nil {
		return aws.ToString(dbInstance.DBClusterIdentifier), nil
	}
	return "-", nil
}

func getDBInstanceID(instance any) (string, error) {
	return aws.ToString(instance.(types.DBInstance).DBInstanceIdentifier), nil
}

func getDBInstanceStatus(instance any) (string, error) {
	return format.Status(aws.ToString(instance.(types.DBInstance).DBInstanceStatus)), nil
}

func getDBInstanceEngine(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.Engine == nil {
		return "", nil
	}
	return string(*dbInstance.Engine), nil
}

func getDBInstanceEngineVersion(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.EngineVersion == nil {
		return "", nil
	}
	return string(*dbInstance.EngineVersion), nil
}

func getDBInstanceSize(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.DBInstanceClass == nil {
		return "", nil
	}
	return string(*dbInstance.DBInstanceClass), nil
}

func getDBInstanceRoleField(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	return calculateDBInstanceRole(dbInstance, clustersContext), nil
}

func getDBInstanceEndpoint(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.Endpoint == nil || dbInstance.Endpoint.Address == nil {
		return "", nil
	}
	return aws.ToString(dbInstance.Endpoint.Address), nil
}

// -----------------------------------------------------------------------------
// Helper functions
// -----------------------------------------------------------------------------

// calculateDBInstanceRole calculates the role of the RDS instance (same logic as original)
func calculateDBInstanceRole(instance types.DBInstance, clusters []types.DBCluster) string {
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
