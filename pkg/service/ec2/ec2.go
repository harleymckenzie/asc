package ec2

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/olekukonko/tablewriter"
)

// EC2Service is a struct that holds the EC2 client.
type EC2Service struct {
	Client *ec2.Client
}

type Column struct {
	Header    string
	GetValue  func(types.Instance) string
	GetColour func(types.Instance) tablewriter.Colors
}

var availableColumns = map[string]Column{
	"name": {
		Header: "Name",
		GetValue: func(i types.Instance) string {
			return getInstanceName(i)
		},
		GetColour: func(i types.Instance) tablewriter.Colors {
			return tablewriter.Colors{}
		},
	},
	"instance_id": {
		Header: "Instance ID",
		GetValue: func(i types.Instance) string {
			return aws.ToString(i.InstanceId)
		},
		GetColour: func(i types.Instance) tablewriter.Colors {
			return tablewriter.Colors{}
		},
	},
	"state": {
		Header: "State",
		GetValue: func(i types.Instance) string {
			return string(i.State.Name)
		},
		GetColour: func(i types.Instance) tablewriter.Colors {
			stateColors := map[types.InstanceStateName]tablewriter.Colors{
				types.InstanceStateNameRunning:    {tablewriter.FgGreenColor},
				types.InstanceStateNameStopped:    {tablewriter.FgRedColor},
				types.InstanceStateNamePending:    {tablewriter.FgYellowColor},
				types.InstanceStateNameTerminated: {tablewriter.FgRedColor},
			}
			if color, exists := stateColors[i.State.Name]; exists {
				return color
			}
			return tablewriter.Colors{}
		},
	},
	"instance_type": {
		Header: "Type",
		GetValue: func(i types.Instance) string {
			return string(i.InstanceType)
		},
		GetColour: func(i types.Instance) tablewriter.Colors {
			return tablewriter.Colors{}
		},
	},
	"public_ip": {
		Header: "Public IP",
		GetValue: func(i types.Instance) string {
			return aws.ToString(i.PublicIpAddress)
		},
		GetColour: func(i types.Instance) tablewriter.Colors {
			return tablewriter.Colors{}
		},
	},
}

func NewEC2Service(ctx context.Context, profile string) (*EC2Service, error) {
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

	client := ec2.NewFromConfig(cfg)
	return &EC2Service{Client: client}, nil
}

func (svc *EC2Service) ListInstances(ctx context.Context) error {
	output, err := svc.Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		log.Printf("Failed to describe instances: %v", err)
		return err
	}

	var instances []types.Instance
	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			if instance.State.Name == types.InstanceStateNameRunning {
				instances = append(instances, instance)
			}
		}
	}

	return PrintInstances(instances)
}

func PrintInstances(instances []types.Instance) error {
	// Define which columns to display
	selectedColumns := []string{"name", "instance_id", "state", "instance_type", "public_ip"}

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

func buildTableData(instances []types.Instance,
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

func getInstanceName(instance types.Instance) string {
	// Get instance name from tags
	name := "-" // Use as default name if "Name" tag doesn't exist
	for _, tag := range instance.Tags {
		if aws.ToString(tag.Key) == "Name" {
			name = aws.ToString(tag.Value)
			break
		}
	}

	return name
}
