package ecs

import (
	"context"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ecs/types"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)

// ECSClientAPI is the interface for the ECS client.
type ECSClientAPI interface {
	ListClusters(context.Context, *ecs.ListClustersInput, ...func(*ecs.Options)) (*ecs.ListClustersOutput, error)
	DescribeClusters(context.Context, *ecs.DescribeClustersInput, ...func(*ecs.Options)) (*ecs.DescribeClustersOutput, error)
	ListServices(context.Context, *ecs.ListServicesInput, ...func(*ecs.Options)) (*ecs.ListServicesOutput, error)
	DescribeServices(context.Context, *ecs.DescribeServicesInput, ...func(*ecs.Options)) (*ecs.DescribeServicesOutput, error)
	ListTasks(context.Context, *ecs.ListTasksInput, ...func(*ecs.Options)) (*ecs.ListTasksOutput, error)
	DescribeTasks(context.Context, *ecs.DescribeTasksInput, ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error)
	ListTaskDefinitionFamilies(context.Context, *ecs.ListTaskDefinitionFamiliesInput, ...func(*ecs.Options)) (*ecs.ListTaskDefinitionFamiliesOutput, error)
	ListTaskDefinitions(context.Context, *ecs.ListTaskDefinitionsInput, ...func(*ecs.Options)) (*ecs.ListTaskDefinitionsOutput, error)
	DescribeTaskDefinition(context.Context, *ecs.DescribeTaskDefinitionInput, ...func(*ecs.Options)) (*ecs.DescribeTaskDefinitionOutput, error)
}

// ECSService is the service for the ECS client.
type ECSService struct {
	Client ECSClientAPI
}

// NewECSService creates a new ECS service.
func NewECSService(ctx context.Context, profile string, region string) (*ECSService, error) {
	cfg, err := awsutil.LoadDefaultConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	client := ecs.NewFromConfig(cfg.Config)
	return &ECSService{Client: client}, nil
}

// ListClusters lists all ECS cluster ARNs.
func (svc *ECSService) ListClusters(ctx context.Context, input *ascTypes.ListClustersInput) ([]string, error) {
	var allARNs []string
	var nextToken *string

	for {
		output, err := svc.Client.ListClusters(ctx, &ecs.ListClustersInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		allARNs = append(allARNs, output.ClusterArns...)
		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
	}

	return allARNs, nil
}

// DescribeClusters describes the specified ECS clusters.
func (svc *ECSService) DescribeClusters(ctx context.Context, input *ascTypes.DescribeClustersInput) ([]types.Cluster, error) {
	output, err := svc.Client.DescribeClusters(ctx, &ecs.DescribeClustersInput{
		Clusters: input.ClusterARNs,
		Include:  []types.ClusterField{types.ClusterFieldTags, types.ClusterFieldStatistics},
	})
	if err != nil {
		return nil, err
	}

	return output.Clusters, nil
}

// ListServices lists all ECS services in the specified cluster.
func (svc *ECSService) ListServices(ctx context.Context, input *ascTypes.ListServicesInput) ([]string, error) {
	var allARNs []string
	var nextToken *string

	params := &ecs.ListServicesInput{
		NextToken: nextToken,
	}
	if input.Cluster != "" {
		params.Cluster = &input.Cluster
	}

	for {
		output, err := svc.Client.ListServices(ctx, params)
		if err != nil {
			return nil, err
		}
		allARNs = append(allARNs, output.ServiceArns...)
		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return allARNs, nil
}

// DescribeServices describes the specified ECS services.
func (svc *ECSService) DescribeServices(ctx context.Context, input *ascTypes.DescribeServicesInput) ([]types.Service, error) {
	output, err := svc.Client.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  &input.Cluster,
		Services: input.Services,
		Include:  []types.ServiceField{types.ServiceFieldTags},
	})
	if err != nil {
		return nil, err
	}

	return output.Services, nil
}

// ListTasks lists all ECS tasks in the specified cluster.
func (svc *ECSService) ListTasks(ctx context.Context, input *ascTypes.ListTasksInput) ([]string, error) {
	var allARNs []string
	var nextToken *string

	params := &ecs.ListTasksInput{
		NextToken: nextToken,
	}
	if input.Cluster != "" {
		params.Cluster = &input.Cluster
	}
	if input.ServiceName != "" {
		params.ServiceName = &input.ServiceName
	}

	for {
		output, err := svc.Client.ListTasks(ctx, params)
		if err != nil {
			return nil, err
		}
		allARNs = append(allARNs, output.TaskArns...)
		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return allARNs, nil
}

// DescribeTasks describes the specified ECS tasks.
func (svc *ECSService) DescribeTasks(ctx context.Context, input *ascTypes.DescribeTasksInput) ([]types.Task, error) {
	output, err := svc.Client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &input.Cluster,
		Tasks:   input.Tasks,
		Include: []types.TaskField{types.TaskFieldTags},
	})
	if err != nil {
		return nil, err
	}

	return output.Tasks, nil
}

// ListTaskDefinitionFamilies lists all ECS task definition families.
func (svc *ECSService) ListTaskDefinitionFamilies(ctx context.Context, input *ascTypes.ListTaskDefinitionFamiliesInput) ([]string, error) {
	var allFamilies []string
	var nextToken *string

	for {
		output, err := svc.Client.ListTaskDefinitionFamilies(ctx, &ecs.ListTaskDefinitionFamiliesInput{
			NextToken: nextToken,
			Status:    types.TaskDefinitionFamilyStatusActive,
		})
		if err != nil {
			return nil, err
		}
		allFamilies = append(allFamilies, output.Families...)
		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
	}

	return allFamilies, nil
}

// ListTaskDefinitionRevisions lists all revisions for a task definition family.
func (svc *ECSService) ListTaskDefinitionRevisions(ctx context.Context, input *ascTypes.ListTaskDefinitionRevisionsInput) ([]string, error) {
	var allARNs []string
	var nextToken *string

	for {
		output, err := svc.Client.ListTaskDefinitions(ctx, &ecs.ListTaskDefinitionsInput{
			FamilyPrefix: &input.FamilyName,
			NextToken:    nextToken,
		})
		if err != nil {
			return nil, err
		}
		allARNs = append(allARNs, output.TaskDefinitionArns...)
		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
	}

	return allARNs, nil
}

// DescribeTaskDefinition describes a task definition.
func (svc *ECSService) DescribeTaskDefinition(ctx context.Context, input *ascTypes.DescribeTaskDefinitionInput) (*types.TaskDefinition, error) {
	output, err := svc.Client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &input.TaskDefinition,
		Include:        []types.TaskDefinitionField{types.TaskDefinitionFieldTags},
	})
	if err != nil {
		return nil, err
	}

	return output.TaskDefinition, nil
}

// GetAllClusters lists and describes all clusters.
func (svc *ECSService) GetAllClusters(ctx context.Context) ([]types.Cluster, error) {
	arns, err := svc.ListClusters(ctx, &ascTypes.ListClustersInput{})
	if err != nil {
		return nil, err
	}

	if len(arns) == 0 {
		return nil, nil
	}

	return svc.DescribeClusters(ctx, &ascTypes.DescribeClustersInput{ClusterARNs: arns})
}

// GetAllServices lists and describes all services, optionally filtered by cluster.
func (svc *ECSService) GetAllServices(ctx context.Context, cluster string) ([]types.Service, error) {
	if cluster != "" {
		return svc.getServicesForCluster(ctx, cluster)
	}

	// Get all clusters first, then get services for each
	clusters, err := svc.GetAllClusters(ctx)
	if err != nil {
		return nil, err
	}

	var allServices []types.Service
	for _, c := range clusters {
		services, err := svc.getServicesForCluster(ctx, *c.ClusterArn)
		if err != nil {
			return nil, fmt.Errorf("list services for cluster %s: %w", *c.ClusterName, err)
		}
		allServices = append(allServices, services...)
	}

	return allServices, nil
}

// getServicesForCluster lists and describes all services in a specific cluster.
func (svc *ECSService) getServicesForCluster(ctx context.Context, cluster string) ([]types.Service, error) {
	arns, err := svc.ListServices(ctx, &ascTypes.ListServicesInput{Cluster: cluster})
	if err != nil {
		return nil, err
	}

	if len(arns) == 0 {
		return nil, nil
	}

	return svc.DescribeServices(ctx, &ascTypes.DescribeServicesInput{
		Cluster:  cluster,
		Services: arns,
	})
}

// GetAllTasks lists and describes all tasks, optionally filtered by cluster and/or service.
func (svc *ECSService) GetAllTasks(ctx context.Context, cluster string, serviceName string) ([]types.Task, error) {
	if cluster != "" {
		return svc.getTasksForCluster(ctx, cluster, serviceName)
	}

	clusters, err := svc.GetAllClusters(ctx)
	if err != nil {
		return nil, err
	}

	var allTasks []types.Task
	for _, c := range clusters {
		tasks, err := svc.getTasksForCluster(ctx, *c.ClusterArn, serviceName)
		if err != nil {
			return nil, fmt.Errorf("list tasks for cluster %s: %w", *c.ClusterName, err)
		}
		allTasks = append(allTasks, tasks...)
	}

	return allTasks, nil
}

// getTasksForCluster lists and describes all tasks in a specific cluster.
func (svc *ECSService) getTasksForCluster(ctx context.Context, cluster string, serviceName string) ([]types.Task, error) {
	arns, err := svc.ListTasks(ctx, &ascTypes.ListTasksInput{Cluster: cluster, ServiceName: serviceName})
	if err != nil {
		return nil, err
	}

	if len(arns) == 0 {
		return nil, nil
	}

	return svc.DescribeTasks(ctx, &ascTypes.DescribeTasksInput{
		Cluster: cluster,
		Tasks:   arns,
	})
}

// ShortARN extracts the short name from an ECS ARN (last segment after /).
func ShortARN(arn string) string {
	return path.Base(arn)
}
