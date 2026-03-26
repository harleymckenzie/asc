package elasticache

import (
	"context"
	"fmt"
	"strings"

	ecsdk "github.com/aws/aws-sdk-go-v2/service/elasticache"
)

var elasticacheTerminalStates = map[string]bool{
	"available":             true,
	"deleted":               true,
	"create-failed":         true,
	"incompatible-network":  true,
	"restore-failed":        true,
	"snapshotting":          true,
}

// IsTerminalClusterState returns true if the ElastiCache cluster status is stable.
func IsTerminalClusterState(status string) bool {
	return elasticacheTerminalStates[strings.ToLower(status)]
}

// GetClusterStatus returns the current status of an ElastiCache cluster.
func (svc *ElasticacheService) GetClusterStatus(ctx context.Context, clusterID string) (string, error) {
	output, err := svc.Client.DescribeCacheClusters(ctx, &ecsdk.DescribeCacheClustersInput{
		CacheClusterId: &clusterID,
	})
	if err != nil {
		return "", fmt.Errorf("describe cache cluster: %w", err)
	}
	if len(output.CacheClusters) == 0 {
		return "", fmt.Errorf("cache cluster %s not found", clusterID)
	}
	if output.CacheClusters[0].CacheClusterStatus == nil {
		return "", fmt.Errorf("cache cluster %s has no status", clusterID)
	}
	return *output.CacheClusters[0].CacheClusterStatus, nil
}
