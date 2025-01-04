package elasticache

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
)

type mockElasticacheClient struct {
	describeCacheClustersOutput *elasticache.DescribeCacheClustersOutput
	err                         error
}

func (m *mockElasticacheClient) DescribeCacheClusters(
	_ context.Context,
	params *elasticache.DescribeCacheClustersInput,
	_ ...func(*elasticache.Options),
) (*elasticache.DescribeCacheClustersOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.describeCacheClustersOutput, nil
}

func TestListInstances(t *testing.T) {
	testCases := []struct {
		name         string
		clusters     []types.CacheCluster
		showEndpoint bool
		err          error
		wantErr      bool
	}{
		{
			name: "mixed instance types without endpoints",
			clusters: []types.CacheCluster{
				{
					CacheClusterId:     aws.String("redis-cluster"),
					CacheClusterStatus: aws.String("available"),
					Engine:             aws.String("redis"),
					EngineVersion:      aws.String("6.2.6"),
					CacheNodeType:      aws.String("cache.m5.large"),
					CacheNodes: []types.CacheNode{
						{
							CacheNodeId: aws.String("node-1"),
							Endpoint:    &types.Endpoint{Address: aws.String("redis-cluster.1234567890.eu-west-1.cache.amazonaws.com")},
						},
					},
				},
				{
					CacheClusterId:     aws.String("memcached-cluster"),
					CacheClusterStatus: aws.String("available"),
					Engine:             aws.String("memcached"),
					EngineVersion:      aws.String("1.6.6"),
					CacheNodeType:      aws.String("cache.t3.medium"),
					CacheNodes: []types.CacheNode{
						{
							CacheNodeId: aws.String("node-1"),
							Endpoint:    &types.Endpoint{Address: aws.String("memcached-cluster.1234567890.eu-west-1.cache.amazonaws.com")},
						},
					},
				},
			},
			showEndpoint: false,
			err:          nil,
			wantErr:      false,
		},
		{
			name: "mixed instance types with endpoints",
			clusters: []types.CacheCluster{
				{
					CacheClusterId:     aws.String("redis-cluster"),
					CacheClusterStatus: aws.String("available"),
					Engine:             aws.String("redis"),
					EngineVersion:      aws.String("6.2.6"),
					CacheNodeType:      aws.String("cache.m5.large"),
					CacheNodes: []types.CacheNode{
						{
							CacheNodeId: aws.String("node-1"),
							Endpoint:    &types.Endpoint{Address: aws.String("redis-cluster.1234567890.eu-west-1.cache.amazonaws.com")},
						},
					},
				},
				{
					CacheClusterId:     aws.String("memcached-cluster"),
					CacheClusterStatus: aws.String("available"),
					Engine:             aws.String("memcached"),
					EngineVersion:      aws.String("1.6.6"),
					CacheNodeType:      aws.String("cache.t3.medium"),
					CacheNodes: []types.CacheNode{
						{
							CacheNodeId: aws.String("node-1"),
							Endpoint:    &types.Endpoint{Address: aws.String("memcached-cluster.1234567890.eu-west-1.cache.amazonaws.com")},
						},
					},
				},
			},
			showEndpoint: true,
			err:          nil,
			wantErr:      false,
		},
		{
			name:     "empty response",
			clusters: []types.CacheCluster{},
			err:      nil,
			wantErr:  false,
		},
		{
			name:     "api error",
			clusters: nil,
			err:      &types.CacheClusterNotFoundFault{Message: aws.String("Cache cluster not found")},
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockElasticacheClient{
				describeCacheClustersOutput: &elasticache.DescribeCacheClustersOutput{
					CacheClusters: tc.clusters,
				},
				err: tc.err,
			}

			svc := &ElasticacheService{
				Client: mockClient,
				ctx:    context.Background(),
			}

			err := svc.ListInstances(context.Background(), tc.showEndpoint)
			if (err != nil) != tc.wantErr {
				t.Errorf("ListInstances() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
