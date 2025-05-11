package elasticache

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"

	"github.com/harleymckenzie/asc/pkg/shared/awsutil"
)


type ElasticacheClientAPI interface {
	DescribeCacheClusters(context.Context, *elasticache.DescribeCacheClustersInput, ...func(*elasticache.Options)) (*elasticache.DescribeCacheClustersOutput, error)
}

// ElasticacheService is a struct that holds the Elasticache client.
type ElasticacheService struct {
	Client ElasticacheClientAPI
}

//
// Service functions
//

func NewElasticacheService(ctx context.Context, profile string, region string) (*ElasticacheService, error) {
	cfg, err := awsutil.LoadDefaultConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	client := elasticache.NewFromConfig(cfg.Config)
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
