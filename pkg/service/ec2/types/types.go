package types

type GetInstancesInput struct {

	// The IDs of the instances to get
	InstanceIDs []string
}

type GetVolumesInput struct {

	// The IDs of the volumes to get
	VolumeIDs []string
}

type GetSnapshotsInput struct {

	// The IDs of the snapshots to get
	SnapshotIDs []string
}

type GetImagesInput struct {

	// The IDs of the images to get
	ImageIDs []string
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

type GetSecurityGroupsInput struct {
	// The IDs of the security groups to get
	GroupIDs []string
}
