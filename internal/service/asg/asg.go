package asg

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"

	ascTypes "github.com/harleymckenzie/asc/internal/service/asg/types"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)

type AutoScalingClientAPI interface {
	DescribeAutoScalingGroups(ctx context.Context, params *autoscaling.DescribeAutoScalingGroupsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
	DescribeScheduledActions(ctx context.Context, params *autoscaling.DescribeScheduledActionsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeScheduledActionsOutput, error)
	PutScheduledUpdateGroupAction(ctx context.Context, params *autoscaling.PutScheduledUpdateGroupActionInput, optFns ...func(*autoscaling.Options)) (*autoscaling.PutScheduledUpdateGroupActionOutput, error)
	DeleteScheduledAction(ctx context.Context, params *autoscaling.DeleteScheduledActionInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DeleteScheduledActionOutput, error)
	UpdateAutoScalingGroup(ctx context.Context, params *autoscaling.UpdateAutoScalingGroupInput, optFns ...func(*autoscaling.Options)) (*autoscaling.UpdateAutoScalingGroupOutput, error)
}

// AutoScalingService is a struct that holds the AutoScaling client.
type AutoScalingService struct {
	Client AutoScalingClientAPI
}

//
// Service functions
//

func NewAutoScalingService(ctx context.Context, profile string, region string) (*AutoScalingService, error) {
	cfg, err := awsutil.LoadDefaultConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}


	client := autoscaling.NewFromConfig(cfg.Config)
	return &AutoScalingService{Client: client}, nil
}

func (svc *AutoScalingService) AddAutoScalingGroupSchedule(ctx context.Context, input *ascTypes.AddAutoScalingGroupScheduleInput) error {
	putScheduledUpdateGroupActionInput := &autoscaling.PutScheduledUpdateGroupActionInput{
		AutoScalingGroupName: &input.AutoScalingGroupName,
		ScheduledActionName:  &input.ScheduledActionName,
		MinSize:              input.MinSize,
		MaxSize:              input.MaxSize,
		DesiredCapacity:      input.DesiredCapacity,
		Recurrence:           input.Recurrence,
		StartTime:            input.StartTime,
		EndTime:              input.EndTime,
	}

	_, err := svc.Client.PutScheduledUpdateGroupAction(ctx, putScheduledUpdateGroupActionInput)
	if err != nil {
		return err
	}
	return nil
}

func (svc *AutoScalingService) GetAutoScalingGroups(ctx context.Context, input *ascTypes.GetAutoScalingGroupsInput) ([]types.AutoScalingGroup, error) {
	output, err := svc.Client.DescribeAutoScalingGroups(ctx, &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: input.AutoScalingGroupNames,
	})
	if err != nil {
		return nil, err
	}

	var autoScalingGroups []types.AutoScalingGroup
	autoScalingGroups = append(autoScalingGroups, output.AutoScalingGroups...)
	return autoScalingGroups, nil
}

func (svc *AutoScalingService) GetAutoScalingGroupInstances(ctx context.Context, input *ascTypes.GetAutoScalingGroupInstancesInput) ([]types.Instance, error) {
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

func (svc *AutoScalingService) GetAutoScalingGroupSchedules(ctx context.Context, input *ascTypes.GetAutoScalingGroupSchedulesInput) ([]types.ScheduledUpdateGroupAction, error) {
	describeScheduledActionsInput := &autoscaling.DescribeScheduledActionsInput{}
	if input.AutoScalingGroupName != "" {
		describeScheduledActionsInput.AutoScalingGroupName = &input.AutoScalingGroupName
	}
	if len(input.ScheduledActionNames) > 0 {
		describeScheduledActionsInput.ScheduledActionNames = input.ScheduledActionNames
	}

	output, err := svc.Client.DescribeScheduledActions(ctx, describeScheduledActionsInput)
	if err != nil {
		return nil, err
	}

	var schedules []types.ScheduledUpdateGroupAction
	schedules = append(schedules, output.ScheduledUpdateGroupActions...)
	return schedules, nil
}

func (svc *AutoScalingService) ModifyAutoScalingGroup(ctx context.Context, input *ascTypes.ModifyAutoScalingGroupInput) error {
	modifyAutoScalingGroupInput := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: &input.AutoScalingGroupName,
		MinSize:              input.MinSize,
		MaxSize:              input.MaxSize,
		DesiredCapacity:      input.DesiredCapacity,
	}
	
	_, err := svc.Client.UpdateAutoScalingGroup(ctx, modifyAutoScalingGroupInput)
	if err != nil {
		return err
	}
	return nil
}

func (svc *AutoScalingService) RemoveAutoScalingGroupSchedule(ctx context.Context, input *ascTypes.RemoveAutoScalingGroupScheduleInput) error {
	_, err := svc.Client.DeleteScheduledAction(ctx, &autoscaling.DeleteScheduledActionInput{
		AutoScalingGroupName: &input.AutoScalingGroupName,
		ScheduledActionName:  &input.ScheduledActionName,
	})
	if err != nil {
		return err
	}
	return nil
}
