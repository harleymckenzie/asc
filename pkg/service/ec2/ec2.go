package ec2

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"strconv"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
	"github.com/harleymckenzie/asc/pkg/shared/awsutil"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// EC2Table implements TableData interface
type EC2Table struct {
	Instances       []types.Instance
	SelectedColumns []string
	SortBy          string
}

type EC2DetailTable struct {
	Instance types.Instance
	Fields   []tableformat.Field
}

type EC2ClientAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	RebootInstances(ctx context.Context, params *ec2.RebootInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RebootInstancesOutput, error)
	StartInstances(ctx context.Context, params *ec2.StartInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error)
	StopInstances(ctx context.Context, params *ec2.StopInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StopInstancesOutput, error)
	TerminateInstances(ctx context.Context, params *ec2.TerminateInstancesInput, optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}

// EC2Service is a struct that holds the EC2 client.
type EC2Service struct {
	Client EC2ClientAPI
}

func availableColumns() map[string]ascTypes.ColumnDef {
	return map[string]ascTypes.ColumnDef{
		"Name": {
			GetValue: func(i *types.Instance) string {
				return getInstanceName(*i)
			},
		},
		"Instance ID": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.InstanceId)
			},
		},
		"State": {
			GetValue: func(i *types.Instance) string {
				return tableformat.FormatState(string(i.State.Name))
			},
		},
		"Instance Type": {
			GetValue: func(i *types.Instance) string {
				return string(i.InstanceType)
			},
		},
		"AMI ID": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.ImageId)
			},
		},
		"AMI Name": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.ImageId)
			},
		},
		"Public IP": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.PublicIpAddress)
			},
		},
		"Private IP": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.PrivateIpAddress)
			},
		},
		"Launch Time": {
			GetValue: func(i *types.Instance) string {
				return i.LaunchTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"Subnet ID": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.SubnetId)
			},
		},
		"Security Group(s)": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.SecurityGroups[0].GroupId)
			},
		},
		"Key Name": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.KeyName)
			},
		},
		"VPC ID": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.VpcId)
			},
		},
		"Placement Group": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.Placement.GroupName)
			},
		},
		"Availability Zone": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.Placement.AvailabilityZone)
			},
		},
		"Root Device Type": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString((*string)(&i.RootDeviceType))
			},
		},
		"Root Device Name": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.RootDeviceName)
			},
		},
		"Virtualization Type": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString((*string)(&i.VirtualizationType))
			},
		},
		"vCPUs": {
			GetValue: func(i *types.Instance) string {
				return strconv.Itoa(int(*i.CpuOptions.CoreCount))
			},
		},
	}
}

//
// Table functions
//

// list headers
func (et *EC2Table) Headers() table.Row {
	return tableformat.BuildHeaders(et.SelectedColumns)
}

// list rows
func (et *EC2Table) Rows() []table.Row {
	rows := []table.Row{}
	for _, instance := range et.Instances {
		row := table.Row{}
		for _, colID := range et.SelectedColumns {
			row = append(row, availableColumns()[colID].GetValue(&instance))
		}
		rows = append(rows, row)
	}
	return rows
}

// list column configs
func (et *EC2Table) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{
		{Name: "Name", WidthMax: 40},
		{Name: "Instance ID", WidthMax: 20},
		{Name: "State", WidthMax: 15},
		{Name: "Type", WidthMax: 12},
		{Name: "Public IP", WidthMax: 15},
	}
}

// list table style
func (et *EC2Table) TableStyle() table.Style {
	style := table.StyleRounded
	style.Options.SeparateRows = false
	style.Options.SeparateColumns = true
	style.Options.SeparateHeader = true
	return style
}

// detail headers
func (et *EC2DetailTable) Headers() table.Row {
	headers := []string{}
	for _, field := range et.Fields {
		headers = append(headers, field.ID)
	}
	return tableformat.BuildHeaders(headers)
}

// detail rows
func (et *EC2DetailTable) AppendRows(t table.Writer) {

	for _, field := range et.Fields {
		if field.Header {
			row := table.Row{field.ID, field.ID}
			t.AppendSeparator()
			t.AppendRow(row, table.RowConfig{AutoMerge: true})
			t.AppendSeparator()
		} else {
			row := table.Row{field.ID, availableColumns()[field.ID].GetValue(&et.Instance)}
			t.AppendRow(row)
		}
	}
}

// detail column configs
func (et *EC2DetailTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{
		{Number: 1, Colors: text.Colors{text.Bold}},
	}
}

// detail table style
func (et *EC2DetailTable) TableStyle() table.Style {
	style := table.StyleRounded
	style.Options.SeparateRows = false
	style.Options.SeparateColumns = true
	style.Options.SeparateHeader = true
	return style
}

//
// Service functions
//

func NewEC2Service(ctx context.Context, profile string, region string) (*EC2Service, error) {
	cfg, err := awsutil.LoadDefaultConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}
	client := ec2.NewFromConfig(cfg.Config)

	return &EC2Service{Client: client}, nil
}

// GetInstances fetches EC2 instances and returns them directly.
func (svc *EC2Service) GetInstances(ctx context.Context, input *ascTypes.GetInstancesInput) ([]types.Instance, error) {
	output, err := svc.Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: input.InstanceIDs,
	})
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

func (svc *EC2Service) RestartInstance(ctx context.Context, input *ascTypes.RestartInstanceInput) error {
	_, err := svc.Client.RebootInstances(ctx, &ec2.RebootInstancesInput{
		InstanceIds: []string{input.InstanceID},
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *EC2Service) StartInstance(ctx context.Context, input *ascTypes.StartInstanceInput) error {
	_, err := svc.Client.StartInstances(ctx, &ec2.StartInstancesInput{
		InstanceIds: []string{input.InstanceID},
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *EC2Service) StopInstance(ctx context.Context, input *ascTypes.StopInstanceInput) error {
	_, err := svc.Client.StopInstances(ctx, &ec2.StopInstancesInput{
		InstanceIds: []string{input.InstanceID},
		Force:       &input.Force,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *EC2Service) TerminateInstance(ctx context.Context, input *ascTypes.TerminateInstanceInput) error {
	_, err := svc.Client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{input.InstanceID},
	})
	if err != nil {
		return err
	}
	return nil
}
