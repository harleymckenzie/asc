package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// Attribute is a struct that defines a field in a detailed table.
type Attribute struct {
	GetValue func(*types.DBInstance, []types.DBCluster) string
}

// GetAttributeValue returns the value for a given field and DBInstance.
func GetAttributeValue(fieldID string, instance any, clusters []types.DBCluster) (string, error) {
	inst, ok := instance.(types.DBInstance)
	if !ok {
		return "", fmt.Errorf("instance is not a types.DBInstance")
	}
	attr, ok := availableAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q", fieldID)
	}
	return attr.GetValue(&inst, clusters), nil
}

func availableAttributes() map[string]Attribute {
	return map[string]Attribute{
		"Cluster Identifier": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				if i.DBClusterIdentifier != nil {
					return aws.ToString(i.DBClusterIdentifier)
				}
				return "-"
			},
		},
		"Identifier": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.DBInstanceIdentifier)
			},
		},
		"Status": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.Status(aws.ToString(i.DBInstanceStatus))
			},
		},
		"Engine": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return string(*i.Engine)
			},
		},
		"Engine Version": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return string(*i.EngineVersion)
			},
		},
		"Size": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return string(*i.DBInstanceClass)
			},
		},
		"Role": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return getDBInstanceRole(*i, clusters)
			},
		},
		"Endpoint": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.Endpoint.Address)
			},
		},
	}
}
