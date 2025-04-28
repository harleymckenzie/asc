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
}

type AutoScalingInstanceTable struct {
	Instances       []types.Instance
	SelectedColumns []string
}

type AutoScalingSchedulesTable struct {
	Schedules       []types.ScheduledUpdateGroupAction
	SelectedColumns []string
}

type GetInstancesInput struct {
	AutoScalingGroupNames []string
}

type GetSchedulesInput struct {
	AutoScalingGroupName string
}

type AutoScalingClientAPI interface {
	DescribeAutoScalingGroups(ctx context.Context, params *autoscaling.DescribeAutoScalingGroupsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
	DescribeScheduledActions(ctx context.Context, params *autoscaling.DescribeScheduledActionsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeScheduledActionsOutput, error)
}

// AutoScalingService is a struct that holds the AutoScaling client.
type AutoScalingService struct {
	Client AutoScalingClientAPI
}

// ColumnDef is a definition of a column to display in the table
type columnDef struct {
	GetValue func(*types.AutoScalingGroup) string
}

type instanceColumnDef struct {
	GetValue func(*types.Instance) string
}

type scheduleColumnDef struct {
	GetValue func(*types.ScheduledUpdateGroupAction) string
}

// availableColumns returns a map of column definitions for Auto Scaling Groups
func availableColumns() map[string]columnDef {
	return map[string]columnDef{
		"Name": {
			GetValue: func(i *types.AutoScalingGroup) string {
				return aws.ToString(i.AutoScalingGroupName)
			},
		},
		"Instances": {
			GetValue: func(i *types.AutoScalingGroup) string {
				// TODO: Return count of Instances (Instance[])
				return strconv.Itoa(len(i.Instances))
			},
		},
		"Desired": {
			GetValue: func(i *types.AutoScalingGroup) string {
				return strconv.Itoa(int(*i.DesiredCapacity))
			},
		},
		"Min": {
			GetValue: func(i *types.AutoScalingGroup) string {
				return strconv.Itoa(int(*i.MinSize))
			},
		},
		"Max": {
			GetValue: func(i *types.AutoScalingGroup) string {
				return strconv.Itoa(int(*i.MaxSize))
			},
		},
		"ARN": {
			GetValue: func(i *types.AutoScalingGroup) string {
				return aws.ToString(i.AutoScalingGroupARN)
			},
		},
	}
}

// availableInstanceColumns returns a map of column definitions for instances
func availableInstanceColumns() map[string]instanceColumnDef {
	return map[string]instanceColumnDef{
		"Name": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.InstanceId)
			},
		},
		"State": {
			GetValue: func(i *types.Instance) string {
				return string(i.LifecycleState)
			},
		},
		"Instance Type": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.InstanceType)
			},
		},
		"Launch Template/Configuration": {
			GetValue: func(i *types.Instance) string {
				if i.LaunchTemplate != nil {
					return aws.ToString(i.LaunchTemplate.LaunchTemplateName)
				}
				return aws.ToString(i.LaunchConfigurationName)
			},
		},
		"Availability Zone": {
			GetValue: func(i *types.Instance) string {
				return aws.ToString(i.AvailabilityZone)
			},
		},
		"Health": {
			GetValue: func(i *types.Instance) string {
				return tableformat.ResourceState(aws.ToString(i.HealthStatus))
			},
		},
	}
}

func availableSchedulesColumns() map[string]scheduleColumnDef {
	return map[string]scheduleColumnDef{
		"Auto Scaling Group": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				return aws.ToString(i.AutoScalingGroupName)
			},
		},
		"Name": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				return aws.ToString(i.ScheduledActionName)
			},
		},
		"Recurrence": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				if i.Recurrence == nil {
					return ""
				}
				return aws.ToString(i.Recurrence)
			},
		},
		"Start Time": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				if i.StartTime == nil {
					return ""
				}
				return i.StartTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"End Time": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				if i.EndTime == nil {
					return ""
				}
				return i.EndTime.Local().Format("2006-01-02 15:04:05 MST")
			},
		},
		"Desired Capacity": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				if i.DesiredCapacity == nil {
					return ""
				}
				return strconv.Itoa(int(*i.DesiredCapacity))
			},
		},
		"Min": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				if i.MinSize == nil {
					return ""
				}
				return strconv.Itoa(int(*i.MinSize))
			},
		},
		"Max": {
			GetValue: func(i *types.ScheduledUpdateGroupAction) string {
				if i.MaxSize == nil {
					return ""
				}
				return strconv.Itoa(int(*i.MaxSize))
			},
		},
	}
}

// Header and Row functions for Auto Scaling Groups
func (et *AutoScalingTable) Headers() table.Row {
	return tableformat.BuildHeaders(et.SelectedColumns)
}
func (et *AutoScalingTable) Rows() []table.Row {
	rows := []table.Row{}
	for _, asg := range et.AutoScalingGroups {
		row := table.Row{}
		for _, colID := range et.SelectedColumns {
			row = append(row, availableColumns()[colID].GetValue(&asg))
		}
		rows = append(rows, row)
	}
	return rows
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
	headers := table.Row{}
	for _, colID := range et.SelectedColumns {
		headers = append(headers, colID)
	}
	return headers
}

func (et *AutoScalingInstanceTable) Rows() []table.Row {
	rows := []table.Row{}
	for _, instance := range et.Instances {
		row := table.Row{}
		for _, colID := range et.SelectedColumns {
			row = append(row, availableInstanceColumns()[colID].GetValue(&instance))
		}
		rows = append(rows, row)
	}
	return rows
}

func (et *AutoScalingInstanceTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{}
}

func (et *AutoScalingInstanceTable) TableStyle() table.Style {
	return table.StyleRounded
}

// Header and Row functions for Schedules
func (et *AutoScalingSchedulesTable) Headers() table.Row {
	return tableformat.BuildHeaders(et.SelectedColumns)
}

func (et *AutoScalingSchedulesTable) Rows() []table.Row {
	rows := []table.Row{}
	for _, schedule := range et.Schedules {
		row := table.Row{}
		for _, colID := range et.SelectedColumns {
			row = append(row, availableSchedulesColumns()[colID].GetValue(&schedule))
		}
		rows = append(rows, row)
	}
	return rows
}

func (et *AutoScalingSchedulesTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{
		{Name: "Auto Scaling Group", AutoMerge: true},
	}
}

func (et *AutoScalingSchedulesTable) TableStyle() table.Style {
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

func (svc *AutoScalingService) GetSchedules(ctx context.Context, input *GetSchedulesInput) ([]types.ScheduledUpdateGroupAction, error) {

	describeScheduledActionsInput := &autoscaling.DescribeScheduledActionsInput{}
	if input.AutoScalingGroupName != "" {
		describeScheduledActionsInput.AutoScalingGroupName = &input.AutoScalingGroupName
	}

	output, err := svc.Client.DescribeScheduledActions(ctx, describeScheduledActionsInput)
	if err != nil {
		return nil, err
	}

	var schedules []types.ScheduledUpdateGroupAction
	schedules = append(schedules, output.ScheduledUpdateGroupActions...)
	return schedules, nil
}
