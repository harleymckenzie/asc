package ec2

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/jedib0t/go-pretty/v6/table"
)

// EC2Table implements TableData interface
type EC2Table struct {
	Instances       []types.Instance
	SelectedColumns []string
	SortOrder       []string
}

type ListInstancesInput struct {
	// List of instance IDs to describe
	List bool
	// Columns to display in the table
	SelectedColumns []string
	// Sort order for the instances
	SortOrder []string
}

type EC2ClientAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// EC2Service is a struct that holds the EC2 client.
type EC2Service struct {
	Client EC2ClientAPI
}

// ColumnDef is a definition of a column to display in the table
type columnDef struct {
	Title    string
	GetValue func(*types.Instance) string
}

func availableColumns() map[string]columnDef {
	return map[string]columnDef{
		"name": {
			Title: "Name",
			GetValue: func(i *types.Instance) string {
				return getInstanceName(*i)
			},
		},
		"instance_id": {
			Title: "Instance ID",
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.InstanceId)
			},
		},
		"state": {
			Title: "State",
			GetValue: func(i *types.Instance) string {
				return tableformat.ResourceState(string(i.State.Name))
			},
		},
		"instance_type": {
			Title: "Type",
			GetValue: func(i *types.Instance) string {
				return string(i.InstanceType)
			},
		},
		"ami_id": {
			Title: "AMI ID",
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.ImageId)
			},
		},
		"public_ip": {
			Title: "Public IP",
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.PublicIpAddress)
			},
		},
		"private_ip": {
			Title: "Private IP",
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.PrivateIpAddress)
			},
		},
		"launch_time": {
			Title: "Launch Time",
			GetValue: func(i *types.Instance) string {
				return i.LaunchTime.Format(time.RFC3339)
			},
		},
	}
}

func (et *EC2Table) Headers() table.Row {
	columns := availableColumns()
	headers := table.Row{}
	for _, colID := range et.SelectedColumns {
		headers = append(headers, columns[colID].Title)
	}
	return headers
}

func (et *EC2Table) Rows() []table.Row {
	columns := availableColumns()
	rows := []table.Row{}
	for _, instance := range et.Instances {
		row := table.Row{}
		for _, colID := range et.SelectedColumns {
			row = append(row, columns[colID].GetValue(&instance))
		}
		rows = append(rows, row)
	}
	return rows
}

func (et *EC2Table) SortColumns() []string {
	return et.SortOrder
}

func (et *EC2Table) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{
		{Name: "Name", WidthMax: 40},
		{Name: "Instance ID", WidthMax: 20},
		{Name: "State", WidthMax: 15},
		{Name: "Type", WidthMax: 12},
		{Name: "Public IP", WidthMax: 15},
	}
}

func (et *EC2Table) TableStyle() table.Style {
	style := table.StyleRounded
	style.Options.SeparateRows = false
	style.Options.SeparateColumns = true
	style.Options.SeparateHeader = true
	return style
}

func NewEC2Service(ctx context.Context, profile string, region string) (*EC2Service, error) {
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

	client := ec2.NewFromConfig(cfg)
	return &EC2Service{Client: client}, nil
}

// GetInstances fetches EC2 instances and returns them directly.
func (svc *EC2Service) GetInstances(ctx context.Context) ([]types.Instance, error) {
	output, err := svc.Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	var instances []types.Instance
	for _, reservation := range output.Reservations {
		instances = append(instances, reservation.Instances...)
	}

	return instances, nil
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
