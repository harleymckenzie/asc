package rds

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/harleymckenzie/asc/pkg/service/base"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
)

// ListInstancesCommand handles listing RDS instances
type ListInstancesCommand struct {
	service *RDSService
	options base.ListOptions
}

func NewListInstancesCommand(svc *RDSService, options base.ListOptions) *ListInstancesCommand {
	return &ListInstancesCommand{
		service: svc,
		options: options,
	}
}

func (cmd *ListInstancesCommand) Execute(ctx context.Context) error {
	output, err := cmd.service.Client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		return fmt.Errorf("failed to describe DB instances: %w", err)
	}

	clusterOutput, err := cmd.service.Client.DescribeDBClusters(ctx, &rds.DescribeDBClustersInput{})
	if err != nil {
		return fmt.Errorf("failed to describe DB clusters: %w", err)
	}

	tableData := cmd.convertToTableData(output.DBInstances, clusterOutput.DBClusters)
	return cmd.service.RenderResourceTable(ctx, tableData, tableformat.TableOptions{
		List:            cmd.options.List,
		SortOrder:       cmd.options.SortOrder,
		SelectedColumns: cmd.options.SelectedColumns,
	})
}

func (cmd *ListInstancesCommand) convertToTableData(instances []types.DBInstance,
	clusters []types.DBCluster) []map[string]string {
	tableData := make([]map[string]string, 0, len(instances))

	for _, instance := range instances {
		row := map[string]string{
			"identifier": aws.ToString(instance.DBInstanceIdentifier),
			"status":     tableformat.ResourceState(aws.ToString(instance.DBInstanceStatus)),
			"engine":     aws.ToString(instance.Engine),
			"size":       aws.ToString(instance.DBInstanceClass),
			"role":       getDBInstanceRole(instance, clusters),
		}

		if instance.DBClusterIdentifier != nil {
			row["cluster_identifier"] = aws.ToString(instance.DBClusterIdentifier)
		} else {
			row["cluster_identifier"] = "None"
		}

		if instance.Endpoint != nil {
			row["endpoint"] = aws.ToString(instance.Endpoint.Address)
		}

		tableData = append(tableData, row)
	}

	return tableData
}

// StopInstanceCommand handles stopping RDS instances
type StopInstanceCommand struct {
	service *RDSService
	options base.StateChangeOptions
}

func NewStopInstanceCommand(svc *RDSService, options base.StateChangeOptions) *StopInstanceCommand {
	return &StopInstanceCommand{
		service: svc,
		options: options,
	}
}

func (cmd *StopInstanceCommand) Execute(ctx context.Context) error {
	for _, resource := range cmd.options.ResourceIDs {
		input := &rds.StopDBInstanceInput{
			DBInstanceIdentifier: aws.String(resource.Name),
		}

		if _, err := cmd.service.Client.StopDBInstance(ctx, input); err != nil {
			return fmt.Errorf("failed to stop DB instance %s: %w", resource.Name, err)
		}
	}
	return nil
}

// StartInstanceCommand handles starting RDS instances
type StartInstanceCommand struct {
	service *RDSService
	options base.StateChangeOptions
}

func NewStartInstanceCommand(svc *RDSService, options base.StateChangeOptions) *StartInstanceCommand {
	return &StartInstanceCommand{
		service: svc,
		options: options,
	}
}

func (cmd *StartInstanceCommand) Execute(ctx context.Context) error {
	for _, resource := range cmd.options.ResourceIDs {
		input := &rds.StartDBInstanceInput{
			DBInstanceIdentifier: aws.String(resource.Name),
		}

		if _, err := cmd.service.Client.StartDBInstance(ctx, input); err != nil {
			return fmt.Errorf("failed to start DB instance %s: %w", resource.Name, err)
		}
	}
	return nil
}
