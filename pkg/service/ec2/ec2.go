package ec2

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/harleymckenzie/asc-go/pkg/shared/tableformat"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type EC2ClientAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// EC2Service is a struct that holds the EC2 client.
type EC2Service struct {
	Client EC2ClientAPI
}

// ColumnDef is a definition of a column to display in the table
type columnDef struct {
	id       string
	title    string
	getValue func(*types.Instance) string
}

var availableColumns = []columnDef{
	{
		id:    "name",
		title: "Name",
		getValue: func(i *types.Instance) string {
			return getInstanceName(*i)
		},
	},
	{
		id:    "instance_id",
		title: "Instance ID",
		getValue: func(i *types.Instance) string {
			return aws.ToString(i.InstanceId)
		},
	},
	{
		id:    "state",
		title: "State",
		getValue: func(i *types.Instance) string {
			return tableformat.ResourceState(string(i.State.Name))
		},
	},
	{
		id:    "instance_type",
		title: "Type",
		getValue: func(i *types.Instance) string {
			return string(i.InstanceType)
		},
	},
	{
		id:    "public_ip",
		title: "Public IP",
		getValue: func(i *types.Instance) string {
			return aws.ToString(i.PublicIpAddress)
		},
	},
	{
		id:    "created_time",
		title: "Created Time",
		getValue: func(i *types.Instance) string {
			// Get created time from attachment time for primary network interface
			return i.NetworkInterfaces[0].Attachment.AttachTime.Format(time.RFC3339)
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

func (svc *EC2Service) ListInstances(ctx context.Context, sortOrder []string, list bool) error {
	// Define which columns to display
	selectedColumns := []string{"name", "instance_id", "state", "instance_type", "public_ip"}

	output, err := svc.Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		log.Printf("Failed to describe instances: %v", err)
		return err
	}

	// At this point we have our instances
	var instances []types.Instance
	for _, reservation := range output.Reservations {
		instances = append(instances, reservation.Instances...)
	}

	// Create the table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	headerRow := make(table.Row, 0)
	for _, colID := range selectedColumns {
		for _, col := range availableColumns {
			if col.id == colID {
				headerRow = append(headerRow, col.title)
				break
			}
		}
	}
	t.AppendHeader(headerRow)

	// The following loop is the same across different services, and will eventually
	// be replaced with a shared function.
	for _, instance := range instances {
		// Create empty row for selected instance. Iterate through selected columns
		row := make(table.Row, len(selectedColumns))
		for i, colID := range selectedColumns {
			// Iterate through available columns
			for _, col := range availableColumns {
				// If selected column = selected available column
				if col.id == colID {
					// Add value of getValue to index value (i) in row slice
					row[i] = col.getValue(&instance)
					break
				}
			}
		}
		t.AppendRow(row)
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:     "Name",
			WidthMax: 40,
		},
		{
			Name:     "Instance ID",
			WidthMax: 20,
		},
		{
			Name:     "State",
			WidthMax: 15,
		},
		{
			Name:     "Type",
			WidthMax: 12,
		},
		{
			Name:     "Public IP",
			WidthMax: 15,
		},
	})

	t.SortBy(sortBy(sortOrder))
	setStyle(t, list)
	t.Render()
	return nil
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

func sortBy(sortOrder []string) []table.SortBy {
	sortBy := []table.SortBy{}

	if len(sortOrder) == 0 {
		sortOrder = []string{"Name"}
	}

	for _, sortField := range sortOrder {
		sortBy = append(sortBy, table.SortBy{Name: sortField, Mode: table.Asc})
	}
	return sortBy
}

func setStyle(t table.Writer, list bool) {

	t.SetStyle(table.StyleRounded)
	if list {
		t.Style().Options.DrawBorder = false
		t.Style().Options.SeparateColumns = false
		t.Style().Options.SeparateHeader = false
	} else {
		t.Style().Format.Header = text.FormatTitle
	}
}
