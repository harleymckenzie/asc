package elasticache

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"

	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)


type ElasticacheClientAPI interface {
	DescribeCacheClusters(context.Context, *elasticache.DescribeCacheClustersInput, ...func(*elasticache.Options)) (*elasticache.DescribeCacheClustersOutput, error)
	ListTagsForResource(context.Context, *elasticache.ListTagsForResourceInput, ...func(*elasticache.Options)) (*elasticache.ListTagsForResourceOutput, error)
}

// CacheClusterWithTags wraps a CacheCluster with its tags.
type CacheClusterWithTags struct {
	types.CacheCluster
	Tags []types.Tag
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

func (svc *ElasticacheService) GetInstances(ctx context.Context) ([]CacheClusterWithTags, error) {
	output, err := svc.Client.DescribeCacheClusters(ctx, &elasticache.DescribeCacheClustersInput{
		ShowCacheNodeInfo: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	var instances []CacheClusterWithTags
	for _, cluster := range output.CacheClusters {
		tagged := CacheClusterWithTags{CacheCluster: cluster}
		if cluster.ARN != nil {
			tagOutput, err := svc.Client.ListTagsForResource(ctx, &elasticache.ListTagsForResourceInput{
				ResourceName: cluster.ARN,
			})
			if err == nil {
				tagged.Tags = tagOutput.TagList
			}
		}
		instances = append(instances, tagged)
	}
	return instances, nil
}
