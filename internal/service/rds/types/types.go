package types

type GetInstancesInput struct {

	// The identifier of the instance to get
	InstanceIdentifier string
}

type GetClustersInput struct {

	// The IDs of the clusters to get
	ClusterIdentifier string
}

type ModifyInstanceInput struct {

	// The identifier of the instance to modify
	DBInstanceIdentifier *string

	// Whether to apply the changes immediately
	ApplyImmediately *bool

	// The new instance class
	DBInstanceClass *string

	// The preferred maintenance window
	// Must be in the format ddd:hh24:mi-ddd:hh24:mi .
	// The day values must be mon | tue | wed | thu | fri | sat | sun .
	// Must be in Universal Coordinated Time (UTC).
	// Must not conflict with the preferred backup window.
	// Must be at least 30 minutes.
	PreferredMaintenanceWindow *string
}
