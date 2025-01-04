package rds

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/olekukonko/tablewriter"
)

type RDSClientAPI interface {
	DescribeDBInstances(context.Context, *rds.DescribeDBInstancesInput, ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
	DescribeDBClusters(context.Context, *rds.DescribeDBClustersInput, ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
}

type RDSService struct {
	Client RDSClientAPI
	ctx    context.Context
}

type Column struct {
	Header    string
	GetValue  func(types.DBInstance) string
	GetColour func(types.DBInstance) tablewriter.Colors
}

var availableColumns = map[string]Column{
	"cluster_identifier": {
		Header: "Cluster Identifier",
		GetValue: func(i types.DBInstance) string {
			if i.DBClusterIdentifier != nil {
				return aws.ToString(i.DBClusterIdentifier)
			}
			return "None"
		},
		GetColour: func(i types.DBInstance) tablewriter.Colors {
			return tablewriter.Colors{}
		},
	},
	"identifier": {
		Header: "Identifier",
		GetValue: func(i types.DBInstance) string {
			return aws.ToString(i.DBInstanceIdentifier)
		},
		GetColour: func(i types.DBInstance) tablewriter.Colors {
			return tablewriter.Colors{}
		},
	},
	"status": {
		Header: "Status",
		GetValue: func(i types.DBInstance) string {
			return aws.ToString(i.DBInstanceStatus)
		},
		GetColour: func(i types.DBInstance) tablewriter.Colors {
			stateColors := map[string]tablewriter.Colors{
				"available":  {tablewriter.FgGreenColor},
				"backing-up": {tablewriter.FgYellowColor},
				"creating":   {tablewriter.FgYellowColor},
				"deleting":   {tablewriter.FgRedColor},
				"modifying":  {tablewriter.FgYellowColor},
				"rebooting":  {tablewriter.FgYellowColor},
			}
			return stateColors[aws.ToString(i.DBInstanceStatus)]
		},
	},
	"engine": {
		Header: "Engine",
		GetValue: func(i types.DBInstance) string {
			return string(*i.Engine)
		},
		GetColour: func(i types.DBInstance) tablewriter.Colors {
			return tablewriter.Colors{}
		},
	},
	"size": {
		Header: "Size",
		GetValue: func(i types.DBInstance) string {
			return string(*i.DBInstanceClass)
		},
		GetColour: func(i types.DBInstance) tablewriter.Colors {
			return tablewriter.Colors{}
		},
	},
	"role": {
		Header: "Role",
		GetValue: func(i types.DBInstance) string {
			if i.DBClusterIdentifier == nil {
				return "None"
			}
			return "Pending" // Will be updated in buildTableData
		},
		GetColour: func(i types.DBInstance) tablewriter.Colors {
			return tablewriter.Colors{}
		},
	},
}

func NewRDSService(ctx context.Context, profile string) (*RDSService, error) {
	var cfg aws.Config
	var err error

	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return nil, err
	}

	client := rds.NewFromConfig(cfg)
	return &RDSService{Client: client, ctx: ctx}, nil
}

func (svc *RDSService) getInstanceRole(instance types.DBInstance) string {

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

	output, err := svc.Client.DescribeDBClusters(svc.ctx, &rds.DescribeDBClustersInput{
		DBClusterIdentifier: instance.DBClusterIdentifier,
	})

	if err != nil {
		log.Printf("Failed to describe DB clusters: %v", err)
		return "Unknown"
	}

	for _, cluster := range output.DBClusters {
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

func buildTableData(svc *RDSService, instances []types.DBInstance,
	selectedColumns []string) ([][]string, [][]tablewriter.Colors) {

	var data [][]string
	var colours [][]tablewriter.Colors

	for _, instance := range instances {
		var row []string
		var rowColors []tablewriter.Colors

		for _, colKey := range selectedColumns {
			if col, exists := availableColumns[colKey]; exists {
				value := col.GetValue(instance)
				if colKey == "role" {
					value = svc.getInstanceRole(instance)
				}
				row = append(row, value)
				rowColors = append(rowColors, col.GetColour(instance))
			}
		}

		data = append(data, row)
		colours = append(colours, rowColors)
	}

	return data, colours
}

func (svc *RDSService) ListInstances(ctx context.Context) error {
	output, err := svc.Client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		log.Printf("Failed to describe instances: %v", err)
		return err
	}

	var instances []types.DBInstance
	for _, instance := range output.DBInstances {
		instances = append(instances, instance)
	}

	return svc.PrintInstances(instances)
}

func (svc *RDSService) PrintInstances(instances []types.DBInstance) error {
	selectedColumns := []string{"cluster_identifier", "identifier", "status", "engine", "size", "role"}

	var headers []string
	for _, colKey := range selectedColumns {
		if col, exists := availableColumns[colKey]; exists {
			headers = append(headers, col.Header)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader(headers)
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetCenterSeparator("-")
	table.SetColumnSeparator(" ")
	table.SetBorder(false)
	table.SetRowLine(true)

	data, colours := buildTableData(svc, instances, selectedColumns)

	for i := range data {
		table.Rich(data[i], colours[i])
	}

	table.Render()
	return nil
}
