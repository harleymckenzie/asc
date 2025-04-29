package types

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type ColumnDef struct {
	
    // The function to get the value of the column
    GetValue func(*types.Instance) string
}

type GetInstancesInput struct {
	
	// The IDs of the instances to get
	InstanceIDs []string
}

type StartInstanceInput struct {
	
	// The ID of the instance to start
	InstanceID string
}

type StopInstanceInput struct {
	
	// The ID of the instance to stop
	InstanceID string

	// Whether to force stop the instance
	Force bool
}

type TerminateInstanceInput struct {
	
	// The ID of the instance to terminate
	InstanceID string
}
