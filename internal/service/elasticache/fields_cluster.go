package elasticache

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// cacheClusterFieldValueGetters provides quick lookup for common cache cluster fields.
var cacheClusterFieldValueGetters = map[string]FieldValueGetter{
	"Cache Name":     getCacheName,
	"Status":         getCacheStatus,
	"Engine Version": getCacheEngineVersion,
	"Configuration":  getCacheConfiguration,
	"Endpoint":       getCacheEndpoint,
}

// getCacheClusterFieldValue returns the value of a field for an ElastiCache cluster
func getCacheClusterFieldValue(fieldName string, cluster types.CacheCluster) (string, error) {
	if getter, exists := cacheClusterFieldValueGetters[fieldName]; exists {
		return getter(cluster)
	}
	return "", fmt.Errorf("field %s not found in cacheClusterFieldValueGetters", fieldName)
}

// Individual field value getters

func getCacheName(cluster any) (string, error) {
	return aws.ToString(cluster.(types.CacheCluster).CacheClusterId), nil
}

func getCacheStatus(cluster any) (string, error) {
	status := cluster.(types.CacheCluster).CacheClusterStatus
	if status == nil {
		return "", nil
	}
	return format.Status(string(*status)), nil
}

func getCacheEngineVersion(cluster any) (string, error) {
	c := cluster.(types.CacheCluster)
	engine := aws.ToString(c.Engine)
	version := aws.ToString(c.EngineVersion)
	if engine == "" && version == "" {
		return "", nil
	}
	return fmt.Sprintf("%s (%s)", version, engine), nil
}

func getCacheConfiguration(cluster any) (string, error) {
	nodeType := cluster.(types.CacheCluster).CacheNodeType
	if nodeType == nil {
		return "", nil
	}
	return string(*nodeType), nil
}

func getCacheEndpoint(cluster any) (string, error) {
	c := cluster.(types.CacheCluster)
	if len(c.CacheNodes) == 0 || c.CacheNodes[0].Endpoint == nil {
		return "", nil
	}
	return aws.ToString(c.CacheNodes[0].Endpoint.Address), nil
}