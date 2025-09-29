package elasticache

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// Cache Cluster field getters
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

// -----------------------------------------------------------------------------
// Cache Cluster field getters
// -----------------------------------------------------------------------------

// getCacheName returns the cache cluster identifier
func getCacheName(cluster any) (string, error) {
	return aws.ToString(cluster.(types.CacheCluster).CacheClusterId), nil
}

// getCacheStatus returns the current status of the cache cluster
func getCacheStatus(cluster any) (string, error) {
	status := cluster.(types.CacheCluster).CacheClusterStatus
	if status == nil {
		return "", nil
	}
	return format.Status(string(*status)), nil
}

// getCacheEngineVersion returns the engine version and type in a formatted string
func getCacheEngineVersion(cluster any) (string, error) {
	c := cluster.(types.CacheCluster)
	engine := aws.ToString(c.Engine)
	version := aws.ToString(c.EngineVersion)
	if engine == "" && version == "" {
		return "", nil
	}
	return fmt.Sprintf("%s (%s)", version, engine), nil
}

// getCacheConfiguration returns the node type configuration of the cache cluster
func getCacheConfiguration(cluster any) (string, error) {
	nodeType := cluster.(types.CacheCluster).CacheNodeType
	if nodeType == nil {
		return "", nil
	}
	return string(*nodeType), nil
}

// getCacheEndpoint returns the endpoint address of the cache cluster
func getCacheEndpoint(cluster any) (string, error) {
	c := cluster.(types.CacheCluster)
	if len(c.CacheNodes) == 0 || c.CacheNodes[0].Endpoint == nil {
		return "", nil
	}
	return aws.ToString(c.CacheNodes[0].Endpoint.Address), nil
}
