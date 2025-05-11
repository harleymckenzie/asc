package types

import (
)

type GetInstancesInput struct {
	
	// The IDs of the instances to get
	InstanceIDs []string
}

type RestartInstanceInput struct {
	
	// The ID of the instance to restart
	InstanceID string
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
