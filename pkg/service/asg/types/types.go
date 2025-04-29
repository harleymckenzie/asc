package types

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
)

type AddAutoScalingGroupScheduleInput struct {

	// The name of the Auto Scaling Group to add the schedule to
	AutoScalingGroupName string

	// The name of the scheduled action to add
	ScheduledActionName string

	// The minimum size of the Auto Scaling Group
	MinSize *int32

	// The maximum size of the Auto Scaling Group
	MaxSize *int32

	// The desired capacity of the Auto Scaling Group
	DesiredCapacity *int32

	// The recurrence of the scheduled action
	Recurrence *string

	// The start time of the scheduled action
	StartTime *time.Time

	// The end time of the scheduled action
	EndTime *time.Time
}

type GetAutoScalingGroupsInput struct {

	// The names of the Auto Scaling Groups to get
	AutoScalingGroupNames []string
}

type GetAutoScalingGroupInstancesInput struct {

	// The names of the Auto Scaling Groups to get instances from
	AutoScalingGroupNames []string
}

type GetAutoScalingGroupSchedulesInput struct {

	// The name of the Auto Scaling Group to get schedules from
	AutoScalingGroupName string

	// The names of the scheduled actions to get
	ScheduledActionNames []string
}

type ModifyAutoScalingGroupInput struct {

	// The name of the Auto Scaling Group to modify
	AutoScalingGroupName string
	
	// The minimum size of the Auto Scaling Group
	MinSize *int32

	// The maximum size of the Auto Scaling Group
	MaxSize *int32

	// The desired capacity of the Auto Scaling Group
	DesiredCapacity *int32
}

type RemoveAutoScalingGroupScheduleInput struct {

	// The name of the Auto Scaling Group to remove the schedule from
	AutoScalingGroupName string

	// The name of the scheduled action to remove
	ScheduledActionName string
}

type ColumnDef struct {

	// The function to get the value of the column
	GetValue func(*types.AutoScalingGroup) string
}

type InstanceColumnDef struct {

	// The function to get the value of the column
	GetValue func(*types.Instance) string
}

type ScheduleColumnDef struct {

	// The function to get the value of the column
	GetValue func(*types.ScheduledUpdateGroupAction) string
}
