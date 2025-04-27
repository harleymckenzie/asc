package elasticache

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"

	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/jedib0t/go-pretty/v6/table"
)

type ElasticacheTable struct {
	Instances       []types.CacheCluster
	SelectedColumns []string
}

type ElasticacheClientAPI interface {
	DescribeCacheClusters(context.Context, *elasticache.DescribeCacheClustersInput, ...func(*elasticache.Options)) (*elasticache.DescribeCacheClustersOutput, error)
}

// ElasticacheService is a struct that holds the Elasticache client.
type ElasticacheService struct {
	Client ElasticacheClientAPI
}

type columnDef struct {
	GetValue func(*types.CacheCluster) string
}

func availableColumns() map[string]columnDef {
	return map[string]columnDef{
		"Cache Name": {
			GetValue: func(i *types.CacheCluster) string {
				return aws.ToString(i.CacheClusterId)
			},
		},
		"Status": {
			GetValue: func(i *types.CacheCluster) string {
				return tableformat.ResourceState(string(*i.CacheClusterStatus))
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

func (et *ElasticacheTable) Headers() table.Row {
	return tableformat.BuildHeaders(et.SelectedColumns)
}
func (et *ElasticacheTable) Rows() []table.Row {
	rows := []table.Row{}
	for _, instance := range et.Instances {
		row := table.Row{}
		for _, colID := range et.SelectedColumns {
			row = append(row, availableColumns()[colID].GetValue(&instance))
		}
		rows = append(rows, row)
	}
	return rows
}

func (et *ElasticacheTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{}
}

func (et *ElasticacheTable) TableStyle() table.Style {
	return table.StyleRounded
}

func NewElasticacheService(ctx context.Context, profile string, region string) (*ElasticacheService, error) {
	var cfg aws.Config
	var err error

	opts := []func(*config.LoadOptions) error{}

	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	cfg, err = config.LoadDefaultConfig(ctx, opts...)

	if err != nil {
		return nil, err
	}

	client := elasticache.NewFromConfig(cfg)
	return &ElasticacheService{Client: client}, nil
}

func (svc *ElasticacheService) GetInstances(ctx context.Context) ([]types.CacheCluster, error) {
	output, err := svc.Client.DescribeCacheClusters(ctx, &elasticache.DescribeCacheClustersInput{
		ShowCacheNodeInfo: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	var instances []types.CacheCluster
	instances = append(instances, output.CacheClusters...)
	return instances, nil
}
