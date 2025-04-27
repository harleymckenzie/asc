package asg

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"

	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/jedib0t/go-pretty/v6/table"
)

type AutoScalingTable struct {
	AutoScalingGroups []types.AutoScalingGroup
	SelectedColumns   []string
	SortOrder         []string
}

type AutoScalingInstanceTable struct {
	Instances       []types.Instance
	SelectedColumns []string
	SortOrder       []string
}

type GetAutoScalingGroupsInput struct {
	List            bool
	SelectedColumns []string
	SortOrder       []string
}

type GetInstancesInput struct {
	AutoScalingGroupNames []string
}

type AutoScalingClientAPI interface {
	DescribeAutoScalingGroups(ctx context.Context, params *autoscaling.DescribeAutoScalingGroupsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
}

// AutoScalingService is a struct that holds the AutoScaling client.
type AutoScalingService struct {
	Client AutoScalingClientAPI
}

// ColumnDef is a definition of a column to display in the table
type columnDef struct {
	Title    string
	GetValue func(*types.AutoScalingGroup) string
}

type instanceColumnDef struct {
	Title    string
	GetValue func(*types.Instance) string
}

// availableColumns returns a map of column definitions for Auto Scaling Groups
func availableColumns() map[string]columnDef {
	return map[string]columnDef{
		"name": {
			Title: "Name",
			GetValue: func(i *types.AutoScalingGroup) string {
				return aws.ToString(i.AutoScalingGroupName)
			},
		},
		"instances": {
			Title: "Instances",
			GetValue: func(i *types.AutoScalingGroup) string {
				// TODO: Return count of Instances (Instance[])
				return strconv.Itoa(len(i.Instances))
			},
		},
		"desired_capacity": {
			Title: "Desired",
			GetValue: func(i *types.AutoScalingGroup) string {
				return strconv.Itoa(int(*i.DesiredCapacity))
			},
		},
		"min_capacity": {
			Title: "Min",
			GetValue: func(i *types.AutoScalingGroup) string {
				return strconv.Itoa(int(*i.MinSize))
			},
		},
		"max_capacity": {
			Title: "Max",
			GetValue: func(i *types.AutoScalingGroup) string {
				return strconv.Itoa(int(*i.MaxSize))
			},
		},
		"arn": {
			Title: "ARN",
			GetValue: func(i *types.AutoScalingGroup) string {
				return aws.ToString(i.AutoScalingGroupARN)
			},
		},
	}
}

// instanceColumns returns a map of column definitions for instances
func instanceColumns() map[string]instanceColumnDef {
	return map[string]instanceColumnDef{
		"instance_name": {
			Title: "Name",
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.InstanceId)
			},
		},
		"state": {
			Title: "State",
			GetValue: func(i *types.Instance) string {
				return string(i.LifecycleState)
			},
		},
		"instance_type": {
			Title: "Instance Type",
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.InstanceType)
			},
		},
		"launch_config": {
			Title: "Launch Template/Configuration",
			GetValue: func(i *types.Instance) string {
				if i.LaunchTemplate != nil {
					return aws.ToString(i.LaunchTemplate.LaunchTemplateName)
				}
				return aws.ToString(i.LaunchConfigurationName)
			},
		},
		"availability_zone": {
			Title: "Availability Zone",
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.AvailabilityZone)
			},
		},
		"health": {
			Title: "Health",
			GetValue: func(i *types.Instance) string {
				return tableformat.ResourceState(aws.ToString(i.HealthStatus))
			},
		},
	}
}

// Header and Row functions for Auto Scaling Groups
func (et *AutoScalingTable) Headers() table.Row {
	columns := availableColumns()
	headers := table.Row{}
	for _, colID := range et.SelectedColumns {
		headers = append(headers, columns[colID].Title)
	}
	return headers
}

func (et *AutoScalingTable) Rows() []table.Row {
	columns := availableColumns()
	rows := []table.Row{}
	for _, asg := range et.AutoScalingGroups {
		row := table.Row{}
		for _, colID := range et.SelectedColumns {
			row = append(row, columns[colID].GetValue(&asg))
		}
		rows = append(rows, row)
	}
	return rows
}

func (et *AutoScalingTable) SortColumns() []string {
	return et.SortOrder
}

func (et *AutoScalingTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{
		{Name: "Instances", WidthMin: 9, WidthMax: 9},
		{Name: "Desired", WidthMin: 7, WidthMax: 7},
		{Name: "Min", WidthMin: 5, WidthMax: 5},
		{Name: "Max", WidthMin: 5, WidthMax: 5},
	}
}

func (et *AutoScalingTable) TableStyle() table.Style {
	return table.StyleRounded
}

// Header and Row functions for Instances
func (et *AutoScalingInstanceTable) Headers() table.Row {
	columns := instanceColumns()
	headers := table.Row{}
	for _, colID := range et.SelectedColumns {
		headers = append(headers, columns[colID].Title)
	}
	return headers
}

func (et *AutoScalingInstanceTable) Rows() []table.Row {
	columns := instanceColumns()
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

func (et *AutoScalingInstanceTable) SortColumns() []string {
	return et.SortOrder
}

func (et *AutoScalingInstanceTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{}
}

func (et *AutoScalingInstanceTable) TableStyle() table.Style {
	return table.StyleRounded
}

func NewAutoScalingService(ctx context.Context, profile string, region string) (*AutoScalingService, error) {
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

	client := autoscaling.NewFromConfig(cfg)
	return &AutoScalingService{Client: client}, nil
}

func (svc *AutoScalingService) GetAutoScalingGroups(ctx context.Context) ([]types.AutoScalingGroup, error) {
	output, err := svc.Client.DescribeAutoScalingGroups(ctx, &autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		return nil, err
	}

	var autoScalingGroups []types.AutoScalingGroup
	autoScalingGroups = append(autoScalingGroups, output.AutoScalingGroups...)
	return autoScalingGroups, nil
}

func (svc *AutoScalingService) GetInstances(ctx context.Context, input *GetInstancesInput) ([]types.Instance, error) {
	output, err := svc.Client.DescribeAutoScalingGroups(ctx, &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: input.AutoScalingGroupNames,
	})
	if err != nil {
		return nil, err
	}

	var instances []types.Instance
	for _, autoScalingGroups := range output.AutoScalingGroups {
		instances = append(instances, autoScalingGroups.Instances...)
	}
	return instances, nil
}

// func (svc *AutoScalingService) ListAutoScalingGroups(ctx context.Context, sortOrder []string, list bool, selectedColumns []string) error {
// 	// Set the default sort order to name if no sort order is provided
// 	if len(sortOrder) == 0 {
// 		sortOrder = []string{"Name"}
// 	}

// 	output, err := svc.Client.DescribeAutoScalingGroups(ctx, &autoscaling.DescribeAutoScalingGroupsInput{})
// 	if err != nil {
// 		return err
// 	}

// 	// At this point we have our Auto Scaling Groups
// 	autoScalingGroups := output.AutoScalingGroups

// 	// Create the table
// 	t := table.NewWriter()
// 	t.SetOutputMirror(os.Stdout)

// 	headerRow := make(table.Row, 0)
// 	for _, colID := range selectedColumns {
// 		for _, col := range availableColumns {
// 			if col.id == colID {
// 				headerRow = append(headerRow, col.title)
// 				break
// 			}
// 		}
// 	}
// 	t.AppendHeader(headerRow)

// 	// The following loop is the same across different services, and will eventually
// 	// be replaced with a shared function.
// 	for _, asg := range autoScalingGroups {
// 		// Create empty row for selected instance. Iterate through selected columns
// 		row := make(table.Row, len(selectedColumns))
// 		for i, colID := range selectedColumns {
// 			// Iterate through available columns
// 			for _, col := range availableColumns {
// 				// If selected column = selected available column
// 				if col.id == colID {
// 					// Add value of getValue to index value (i) in row slice
// 					row[i] = col.getValue(&asg)
// 					break
// 				}
// 			}
// 		}
// 		t.AppendRow(row)
// 	}

// 	t.SortBy(tableformat.SortBy(sortOrder))
// 	tableformat.SetStyle(t, list, false, nil)
// 	t.Render()
// 	return nil
// }

// func (svc *AutoScalingService) ListAutoScalingGroupInstances(ctx context.Context, autoScalingGroupName string, sortOrder []string, list bool, selectedColumns []string) error {
// 	// Input provided Auto Scaling Group Name into DescribeAutoScalingGroupsInput
// 	output, err := svc.Client.DescribeAutoScalingGroups(ctx, &autoscaling.DescribeAutoScalingGroupsInput{
// 		AutoScalingGroupNames: []string{autoScalingGroupName},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	var instances []types.Instance
// 	for _, autoScalingGroups := range output.AutoScalingGroups {
// 		instances = append(instances, autoScalingGroups.Instances...)
// 	}

// 	// Create the table
// 	t := table.NewWriter()
// 	t.SetOutputMirror(os.Stdout)

// 	headerRow := make(table.Row, 0)
// 	for _, colID := range selectedColumns {
// 		for _, col := range instanceColumns {
// 			if col.id == colID {
// 				headerRow = append(headerRow, col.title)
// 				break
// 			}
// 		}
// 	}
// 	t.AppendHeader(headerRow)

// 	for _, instance := range instances {
// 		row := make(table.Row, len(selectedColumns))
// 		for i, colID := range selectedColumns {
// 			for _, col := range instanceColumns {
// 				if col.id == colID {
// 					row[i] = col.getValue(&instance)
// 					break
// 				}
// 			}
// 		}
// 		t.AppendRow(row)
// 	}

// 	t.SortBy(tableformat.SortBy(sortOrder))
// 	tableformat.SetStyle(t, list, false, nil)
// 	t.Render()
// 	return nil
// }
