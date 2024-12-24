package elasticache

import (
	"context"
    "fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/olekukonko/tablewriter"
)

// ElasticacheService is a struct that holds the Elasticache client.
type ElasticacheService struct {
	Client *elasticache.Client
}

type Column struct {
	Header    string
	GetValue  func(types.CacheCluster) string
	GetColour func(types.CacheCluster) tablewriter.Colors
}

var availableColumns = map[string]Column{
	"name": {
		Header: "Cache name",
		GetValue: func(i types.CacheCluster) string {
			return aws.ToString(i.CacheClusterId)
		},
		GetColour: func(i types.CacheCluster) tablewriter.Colors {
			return tablewriter.Colors{}
		},
	},
	"status": {
		Header: "Status",
		GetValue: func(i types.CacheCluster) string {
			return string(*i.CacheClusterStatus)
		},
		GetColour: func(i types.CacheCluster) tablewriter.Colors {
			stateColors := map[string]tablewriter.Colors{
				"available": {tablewriter.FgGreenColor},
				"deleting":  {tablewriter.FgRedColor},
				"deleted":   {tablewriter.FgRedColor},
				"rebooting": {tablewriter.FgYellowColor},
			}
			if color, exists := stateColors[aws.ToString(i.CacheClusterStatus)]; exists {
				return color
			}
			return tablewriter.Colors{}
		},
	},
	"engine_version": {
        Header: "Engine version",
		GetValue: func(i types.CacheCluster) string {
            return fmt.Sprintf("%s (%s)", *i.EngineVersion, *i.Engine)
		},
		GetColour: func(i types.CacheCluster) tablewriter.Colors {
			return tablewriter.Colors{}
		},
    },
	"instance_type": {
		Header: "Configuration",
		GetValue: func(i types.CacheCluster) string {
			return string(*i.CacheNodeType)
		},
		GetColour: func(i types.CacheCluster) tablewriter.Colors {
			return tablewriter.Colors{}
		},
	},
}

func NewElasticacheService(ctx context.Context, profile string) (*ElasticacheService, error) {
	var cfg aws.Config
	var err error

	if profile != "" {
		// Load the configuration for the specified profile
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	} else {
		// Use the default configuration
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return nil, err
	}

	client := elasticache.NewFromConfig(cfg)
	return &ElasticacheService{Client: client}, nil
}

func (svc *ElasticacheService) ListInstances(ctx context.Context) error {
	output, err := svc.Client.DescribeCacheClusters(ctx, &elasticache.DescribeCacheClustersInput{})
	if err != nil {
		log.Printf("Failed to describe instances: %v", err)
		return err
	}

	var instances []types.CacheCluster
	for _, instance := range output.CacheClusters {
		instances = append(instances, instance)
	}

	return PrintInstances(instances)
}

func PrintInstances(instances []types.CacheCluster) error {
	// Define which columns to display
	selectedColumns := []string{"name", "status", "engine_version", "instance_type"}

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
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator(" ")
	table.SetBorder(false)

	// Build table data
	data, colours := buildTableData(instances, selectedColumns)

	// Add rows to table
	for i := range data {
		table.Rich(data[i], colours[i])
	}

	table.Render()
	return nil
}

func buildTableData(instances []types.CacheCluster,
	selectedColumns []string) ([][]string, [][]tablewriter.Colors) {

	var data [][]string
	var colours [][]tablewriter.Colors

	for _, instance := range instances {
		var row []string
		var rowColors []tablewriter.Colors

		for _, colKey := range selectedColumns {
			if col, exists := availableColumns[colKey]; exists {
				row = append(row, col.GetValue(instance))
				rowColors = append(rowColors, col.GetColour(instance))
			}
		}

		data = append(data, row)
		colours = append(colours, rowColors)
	}

	return data, colours
}
