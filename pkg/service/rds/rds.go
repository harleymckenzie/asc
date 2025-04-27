package rds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/jedib0t/go-pretty/v6/table"
)

type RDSTable struct {
	Clusters        []types.DBCluster
	Instances       []types.DBInstance
	SelectedColumns []string
}

type RDSClientAPI interface {
	DescribeDBInstances(context.Context, *rds.DescribeDBInstancesInput, ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBClusters(context.Context, *rds.DescribeDBClustersInput, ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
}

type RDSService struct {
	Client RDSClientAPI
}

type columnDef struct {
	GetValue func(*types.DBInstance, []types.DBCluster) string
}

func availableColumns() map[string]columnDef {
	return map[string]columnDef{
		"Cluster Identifier": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				if i.DBClusterIdentifier != nil {
					return aws.ToString(i.DBClusterIdentifier)
				}
				return "-"
			},
		},
		"Identifier": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.DBInstanceIdentifier)
			},
		},
		"Status": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return tableformat.ResourceState(aws.ToString(i.DBInstanceStatus))
			},
		},
		"Engine": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return string(*i.Engine)
			},
		},
		"Engine Version": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return string(*i.EngineVersion)
			},
		},
		"Size": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return string(*i.DBInstanceClass)
			},
		},
		"Role": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return getDBInstanceRole(*i, clusters)
			},
		},
		"Endpoint": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.Endpoint.Address)
			},
		},
	}
}

func (et *RDSTable) Headers() table.Row {
	return tableformat.BuildHeaders(et.SelectedColumns)
}
func (et *RDSTable) Rows() []table.Row {
	rows := []table.Row{}
	for _, instance := range et.Instances {
		row := table.Row{}
		for _, colID := range et.SelectedColumns {
			row = append(row, availableColumns()[colID].GetValue(&instance, et.Clusters))
		}
		rows = append(rows, row)
	}
	return rows
}

func (et *RDSTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{
		{Name: "Cluster Identifier", WidthMax: 40},
		// {Name: "Identifier", WidthMax: 20},
		{Name: "Status", WidthMax: 15},
		{Name: "Engine", WidthMax: 12},
		{Name: "Engine Version", WidthMax: 15},
		// {Name: "Size", WidthMax: 12},
		{Name: "Role", WidthMax: 15},
		{Name: "Endpoint", WidthMax: 15},
	}
}

func (et *RDSTable) TableStyle() table.Style {
	style := table.StyleRounded
	style.Options.SeparateRows = true

	style.Options.SeparateColumns = true
	style.Options.SeparateHeader = true
	return style
}

func NewRDSService(ctx context.Context, profile string, region string) (*RDSService, error) {
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

	client := rds.NewFromConfig(cfg)
	return &RDSService{Client: client}, nil
}

func (svc *RDSService) GetInstances(ctx context.Context) ([]types.DBInstance, error) {
	output, err := svc.Client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, err
	}

	var instances []types.DBInstance
	instances = append(instances, output.DBInstances...)
	return instances, nil
}

func (svc *RDSService) GetClusters(ctx context.Context) ([]types.DBCluster, error) {
	clusterOutput, err := svc.Client.DescribeDBClusters(ctx, &rds.DescribeDBClustersInput{})
	if err != nil {
		return nil, err
	}

	var clusters []types.DBCluster
	clusters = append(clusters, clusterOutput.DBClusters...)
	return clusters, nil
}

func getDBInstanceRole(instance types.DBInstance, clusters []types.DBCluster) string {
	// If ReadReplicaSourceDBInstanceIdentifier is set, then this is a replica. If
	// if ReadReplicaDBInstanceIdentifiers is set, then this is a primary.
	if instance.ReadReplicaSourceDBInstanceIdentifier != nil {
		return "Replica"
	}

	if len(instance.ReadReplicaDBInstanceIdentifiers) > 0 {
		return "Primary"
	}

	if instance.DBClusterIdentifier == nil {
		return "None"
	}

	for _, cluster := range clusters {
		for _, member := range cluster.DBClusterMembers {
			if aws.ToString(member.DBInstanceIdentifier) == aws.ToString(instance.DBInstanceIdentifier) {
				if member.IsClusterWriter != nil && *member.IsClusterWriter {
					return "Writer"
				}
				return "Reader"
			}
		}
	}

	return "Unknown"
}
