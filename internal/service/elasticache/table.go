package elasticache

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// Attribute is a struct that defines a field in a detailed table.
type Attribute struct {
    GetValue func(*types.CacheCluster) string
}

func GetAttributeValue(fieldID string, instance any) (string, error) {
	inst, ok := instance.(types.CacheCluster)
	if !ok {
		return "", fmt.Errorf("instance is not a types.CacheCluster")
	}
	attr, ok := availableAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
        return "", fmt.Errorf("error getting attribute %q", fieldID)
	}
	return attr.GetValue(&inst), nil
}

func availableAttributes() map[string]Attribute {
	return map[string]Attribute{
        "Cache Name": {
            GetValue: func(i *types.CacheCluster) string {
                return aws.ToString(i.CacheClusterId)
            },
        },
        "Status": {
            GetValue: func(i *types.CacheCluster) string {
                return format.Status(string(*i.CacheClusterStatus))
            },
        },
        "Engine Version": {
            GetValue: func(i *types.CacheCluster) string {
                return fmt.Sprintf("%s (%s)", *i.EngineVersion, *i.Engine)
            },
        },
        "Configuration": {
            GetValue: func(i *types.CacheCluster) string {
                return string(*i.CacheNodeType)
            },
        },
        "Endpoint": {
            GetValue: func(i *types.CacheCluster) string {
                return string(*i.CacheNodes[0].Endpoint.Address)
            },
        },
    }
}