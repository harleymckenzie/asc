package elasticache

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/jedib0t/go-pretty/v6/table"
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
			}

			instances, err := svc.GetInstances(context.Background())
			if (err != nil) != tc.wantErr {
				t.Errorf("ListInstances() error = %v, wantErr %v", err, tc.wantErr)
			}

			if len(instances) != len(tc.clusters) {
				t.Errorf("ListInstances() returned %d instances, want %d", len(instances), len(tc.clusters))
			}

			for i, instance := range instances {
				if instance.CacheClusterId != tc.clusters[i].CacheClusterId {
					t.Errorf("ListInstances() returned instance %d with ID %s, want %s", i, *instance.CacheClusterId, *tc.clusters[i].CacheClusterId)
				}
			}
			
		})
	}
}

func TestTableOutput(t *testing.T) {
	clusters := []types.CacheCluster{
		{
			CacheClusterId:     aws.String("redis-cluster"),
			CacheClusterStatus: aws.String("available"),
			Engine:             aws.String("redis"),
			EngineVersion:      aws.String("6.2.6"),
			CacheNodeType:      aws.String("cache.m5.large"),
			CacheNodes: []types.CacheNode{
				{
					CacheNodeId: aws.String("node-1"),
					Endpoint:    &types.Endpoint{Address: aws.String("redis-cluster.example.com")},
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
					Endpoint:    &types.Endpoint{Address: aws.String("memcached-cluster.example.com")},
				},
			},
		},
	}

	testCases := []struct {
		name            string
		selectedColumns []string
		wantHeaders     table.Row
		wantRowCount    int
	}{
		{
			name:            "full cluster details",
			selectedColumns: []string{"Cache Name", "Status", "Engine Version", "Configuration", "Endpoint"},
			wantHeaders:     table.Row{"Cache Name", "Status", "Engine Version", "Configuration", "Endpoint"},
			wantRowCount:    2,
		},
		{
			name:            "minimal columns",
			selectedColumns: []string{"Cache Name", "Status"},
			wantHeaders:     table.Row{"Cache Name", "Status"},
			wantRowCount:    2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			elasticacheTable := &ElasticacheTable{
				Instances:       clusters,
				SelectedColumns: tc.selectedColumns,
			}

			// Test Headers
			headers := elasticacheTable.Headers()
			if len(headers) != len(tc.wantHeaders) {
				t.Errorf("Headers() returned %d columns, want %d", len(headers), len(tc.wantHeaders))
			}

			// Test Rows
			rows := elasticacheTable.Rows()
			if len(rows) != tc.wantRowCount {
				t.Errorf("Rows() returned %d rows, want %d", len(rows), tc.wantRowCount)
			}

			// Print the actual table output
			tw := table.NewWriter()
			tw.AppendHeader(headers)
			tw.AppendRows(rows)
			tw.SetStyle(elasticacheTable.TableStyle())
			t.Logf("\nTable Output:\n%s", tw.Render())
		})
	}
}
