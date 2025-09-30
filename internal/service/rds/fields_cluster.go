package rds

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// RDS DB Cluster field getters
var dbClusterFieldValueGetters = map[string]FieldValueGetter{
	"Activity Stream Kinesis Stream": getDBClusterActivityStreamKinesisStream,
	"Activity Stream KMS Key":        getDBClusterActivityStreamKMSKey,
	"Activity Stream Mode":           getDBClusterActivityStreamMode,
	"Activity Stream Status":         getDBClusterActivityStreamStatus,
	"Allocated Storage":              getDBClusterAllocatedStorage,
	"Associated Roles":               getDBClusterAssociatedRoles,
	"Auto Minor Version Upgrade":     getDBClusterAutoMinorVersionUpgrade,
	"Automatic Restart Time":         getDBClusterAutomaticRestartTime,
	"Availability Zones":             getDBClusterAvailabilityZones,
	"Backup Retention Period":        getDBClusterBackupRetentionPeriod,
	"Capacity":                       getDBClusterCapacity,
	"Certificate Authority":          getDBClusterCertificateAuthority,
	"Certificate Expiry Date":        getDBClusterCertificateExpiryDate,
	"Character Set Name":             getDBClusterCharacterSetName,
	"Clone Group ID":                 getDBClusterCloneGroupID,
	"Cluster Create Time":            getDBClusterClusterCreateTime,
	"Cluster Identifier":             getDBClusterID,
	"Cluster Scalability Type":       getDBClusterClusterScalabilityType,
	"Copy Tags To Snapshot":          getDBClusterCopyTagsToSnapshot,
	"Cross Account Clone":            getDBClusterCrossAccountClone,
	"Custom Endpoints":               getDBClusterCustomEndpoints,
	"DB Cluster ARN":                 getDBClusterARN,
	"DB Cluster Instance Class":      getDBClusterInstanceClass,
	"DB Name":                        getDBClusterDBName,
	"DB Subnet Group":                getDBClusterSubnetGroup,
	"Database Name":                  getDBClusterDatabaseName,
	"Deletion Protection":            getDBClusterDeletionProtection,
	"Domain Memberships":             getDBClusterDomainMemberships,
	"Earliest Backtrack Time":        getDBClusterEarliestBacktrackTime,
	"Earliest Restorable Time":       getDBClusterEarliestRestorableTime,
	"Enabled CloudWatch Logs":        getDBClusterEnabledCloudWatchLogs,
	"Encryption":                     getDBClusterEncryption,
	"Endpoint":                       getDBClusterEndpoint,
	"Engine":                         getDBClusterEngine,
	"Engine Lifecycle Support":       getDBClusterEngineLifecycleSupport,
	"Engine Mode":                    getDBClusterEngineMode,
	"Engine Version":                 getDBClusterEngineVersion,
	"Global Cluster Identifier":      getDBClusterGlobalClusterIdentifier,
	"Global Write Forwarding":        getDBClusterGlobalWriteForwarding,
	"Hosted Zone ID":                 getDBClusterHostedZoneID,
	"HTTP Endpoint Enabled":          getDBClusterHTTPEndpointEnabled,
	"IAM Database Authentication":    getDBClusterIAMDatabaseAuthentication,
	"Identifier":                     getDBClusterID,
	"IOPS":                           getDBClusterIOPS,
	"KMS Key ID":                     getDBClusterKMSKey,
	"Latest Restorable Time":         getDBClusterLatestRestorableTime,
	"Local Write Forwarding":         getDBClusterLocalWriteForwarding,
	"Master Username":                getDBClusterMasterUsername,
	"Monitoring Interval":            getDBClusterMonitoringInterval,
	"Monitoring Role":                getDBClusterMonitoringRole,
	"Multi AZ":                       getDBClusterMultiAZ,
	"Network Type":                   getDBClusterNetworkType,
	"Option Group":                   getDBClusterOptionGroup,
	"Parameter Group":                getDBClusterParameterGroup,
	"Pending Modifications":          getDBClusterPendingModifications,
	"Performance Insights":           getDBClusterPerformanceInsights,
	"Performance Insights KMS Key":   getDBClusterPerformanceInsightsKMSKey,
	"Performance Insights Retention": getDBClusterPerformanceInsightsRetention,
	"Port":                           getDBClusterPort,
	"Preferred Backup Window":        getDBClusterPreferredBackupWindow,
	"Preferred Maintenance Window":   getDBClusterMaintenanceWindow,
	"Publicly Accessible":            getDBClusterPubliclyAccessible,
	"Read Replica Identifiers":       getDBClusterReadReplicaIdentifiers,
	"Reader Endpoint":                getDBClusterReaderEndpoint,
	"Replication Source Identifier":  getDBClusterReplicationSourceIdentifier,
	"Resource ID":                    getDBClusterResourceID,
	"Scaling Configuration":          getDBClusterScalingConfiguration,
	"Serverless V2 Platform Version": getDBClusterServerlessV2PlatformVersion,
	"Serverless V2 Scaling Config":   getDBClusterServerlessV2ScalingConfig,
	"Status":                         getDBClusterStatus,
	"Storage Encrypted":              getDBClusterStorageEncrypted,
	"Storage Throughput":             getDBClusterStorageThroughput,
	"Storage Type":                   getDBClusterStorageType,
	"VPC Security Groups":            getDBClusterVPCSecurityGroups,
}

// SetClustersContext sets the clusters context for role calculation
func SetClustersContext(clusters []types.DBCluster) {
	clustersContext = clusters
}

// getDBClusterFieldValue returns the value of a field for an RDS DB cluster
func getDBClusterFieldValue(fieldName string, instance types.DBCluster) (string, error) {
	if getter, exists := dbClusterFieldValueGetters[fieldName]; exists {
		return getter(instance)
	}
	return "", fmt.Errorf("field %s not found in dbClusterFieldValueGetters", fieldName)
}

// getDBClusterID returns the cluster identifier
func getDBClusterID(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).DBClusterIdentifier), nil
}

// getDBClusterActivityStreamKinesisStream returns the activity stream Kinesis stream name
func getDBClusterActivityStreamKinesisStream(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).ActivityStreamKinesisStreamName), nil
}

// getDBClusterActivityStreamKMSKey returns the activity stream KMS key ID
func getDBClusterActivityStreamKMSKey(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).ActivityStreamKmsKeyId), nil
}

// getDBClusterActivityStreamMode returns the activity stream mode
func getDBClusterActivityStreamMode(instance any) (string, error) {
	return string(instance.(types.DBCluster).ActivityStreamMode), nil
}

// getDBClusterActivityStreamStatus returns the activity stream status
func getDBClusterActivityStreamStatus(instance any) (string, error) {
	return string(instance.(types.DBCluster).ActivityStreamStatus), nil
}

// getDBClusterAllocatedStorage returns the allocated storage size
func getDBClusterAllocatedStorage(instance any) (string, error) {
	storage := instance.(types.DBCluster).AllocatedStorage
	if storage == nil {
		return "", nil
	}
	return fmt.Sprintf("%d GiB", *storage), nil
}

// getDBClusterAssociatedRoles returns the associated IAM roles
func getDBClusterAssociatedRoles(instance any) (string, error) {
	roles := instance.(types.DBCluster).AssociatedRoles
	if len(roles) == 0 {
		return "", nil
	}

	var roleArns []string
	for _, role := range roles {
		roleArns = append(roleArns, aws.ToString(role.RoleArn))
	}
	return strings.Join(roleArns, ", "), nil
}

// getDBClusterAutoMinorVersionUpgrade returns the auto minor version upgrade setting
func getDBClusterAutoMinorVersionUpgrade(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	if cluster.AutoMinorVersionUpgrade == nil {
		return "", nil
	}
	if *cluster.AutoMinorVersionUpgrade {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBClusterAutomaticRestartTime returns the automatic restart time
func getDBClusterAutomaticRestartTime(instance any) (string, error) {
	restartTime := instance.(types.DBCluster).AutomaticRestartTime
	if restartTime == nil {
		return "", nil
	}
	return format.TimeToStringOrDefault(restartTime, ""), nil
}

// getDBClusterAvailabilityZones returns the availability zones
func getDBClusterAvailabilityZones(instance any) (string, error) {
	azs := instance.(types.DBCluster).AvailabilityZones
	if len(azs) == 0 {
		return "", nil
	}
	return strings.Join(azs, ", "), nil
}

// getDBClusterBackupRetentionPeriod returns the backup retention period
func getDBClusterBackupRetentionPeriod(instance any) (string, error) {
	period := instance.(types.DBCluster).BackupRetentionPeriod
	if period == nil {
		return "", nil
	}
	return fmt.Sprintf("%d days", *period), nil
}

// getDBClusterCapacity returns the cluster capacity
func getDBClusterCapacity(instance any) (string, error) {
	capacity := instance.(types.DBCluster).Capacity
	if capacity == nil {
		return "", nil
	}
	return fmt.Sprintf("%d", *capacity), nil
}

// getDBClusterCertificateAuthority returns the certificate authority
func getDBClusterCertificateAuthority(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	cert := cluster.CertificateDetails
	if cert == nil {
		return "Not configured", nil
	}
	caIdentifier := aws.ToString(cert.CAIdentifier)
	if caIdentifier == "" {
		return "Not configured", nil
	}
	return caIdentifier, nil
}

// getDBClusterCertificateExpiryDate returns the certificate expiry date
func getDBClusterCertificateExpiryDate(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	cert := cluster.CertificateDetails
	if cert == nil || cert.ValidTill == nil {
		return "Not configured", nil
	}
	return format.TimeToStringOrDefault(cert.ValidTill, "Not configured"), nil
}

// getDBClusterCharacterSetName returns the character set name
func getDBClusterCharacterSetName(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).CharacterSetName), nil
}

// getDBClusterCloneGroupID returns the clone group ID
func getDBClusterCloneGroupID(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).CloneGroupId), nil
}

// getDBClusterClusterCreateTime returns the cluster creation time
func getDBClusterClusterCreateTime(instance any) (string, error) {
	createTime := instance.(types.DBCluster).ClusterCreateTime
	if createTime == nil {
		return "", nil
	}
	return format.TimeToStringOrDefault(createTime, ""), nil
}

// getDBClusterClusterScalabilityType returns the cluster scalability type
func getDBClusterClusterScalabilityType(instance any) (string, error) {
	return string(instance.(types.DBCluster).ClusterScalabilityType), nil
}

// getDBClusterCopyTagsToSnapshot returns the copy tags to snapshot setting
func getDBClusterCopyTagsToSnapshot(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	if cluster.CopyTagsToSnapshot == nil {
		return "", nil
	}
	if *cluster.CopyTagsToSnapshot {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBClusterCrossAccountClone returns the cross account clone setting
func getDBClusterCrossAccountClone(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	if cluster.CrossAccountClone == nil {
		return "", nil
	}
	if *cluster.CrossAccountClone {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBClusterCustomEndpoints returns the custom endpoints
func getDBClusterCustomEndpoints(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	endpoints := cluster.CustomEndpoints
	if len(endpoints) == 0 {
		return "None", nil
	}
	return strings.Join(endpoints, ", "), nil
}

// getDBClusterARN returns the cluster ARN
func getDBClusterARN(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).DBClusterArn), nil
}

// getDBClusterInstanceClass returns the cluster instance class
func getDBClusterInstanceClass(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).DBClusterInstanceClass), nil
}

// getDBClusterDBName returns the database name
func getDBClusterDBName(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).DatabaseName), nil
}

// getDBClusterDatabaseName returns the database name (alias)
func getDBClusterDatabaseName(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).DatabaseName), nil
}

// getDBClusterSubnetGroup returns the subnet group
func getDBClusterSubnetGroup(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).DBSubnetGroup), nil
}

// getDBClusterDeletionProtection returns the deletion protection setting
func getDBClusterDeletionProtection(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	if cluster.DeletionProtection == nil {
		return "", nil
	}
	if *cluster.DeletionProtection {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBClusterDomainMemberships returns the domain memberships
func getDBClusterDomainMemberships(instance any) (string, error) {
	memberships := instance.(types.DBCluster).DomainMemberships
	if len(memberships) == 0 {
		return "", nil
	}

	var domains []string
	for _, membership := range memberships {
		domains = append(domains, aws.ToString(membership.Domain))
	}
	return strings.Join(domains, ", "), nil
}

// getDBClusterEarliestBacktrackTime returns the earliest backtrack time
func getDBClusterEarliestBacktrackTime(instance any) (string, error) {
	backtrackTime := instance.(types.DBCluster).EarliestBacktrackTime
	if backtrackTime == nil {
		return "", nil
	}
	return format.TimeToStringOrDefault(backtrackTime, ""), nil
}

// getDBClusterEarliestRestorableTime returns the earliest restorable time
func getDBClusterEarliestRestorableTime(instance any) (string, error) {
	restoreTime := instance.(types.DBCluster).EarliestRestorableTime
	if restoreTime == nil {
		return "", nil
	}
	return format.TimeToStringOrDefault(restoreTime, ""), nil
}

// getDBClusterEnabledCloudWatchLogs returns the enabled CloudWatch logs
func getDBClusterEnabledCloudWatchLogs(instance any) (string, error) {
	logs := instance.(types.DBCluster).EnabledCloudwatchLogsExports
	if len(logs) == 0 {
		return "", nil
	}
	return strings.Join(logs, ", "), nil
}

// getDBClusterEncryption returns the encryption setting
func getDBClusterEncryption(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	if cluster.StorageEncrypted == nil {
		return "", nil
	}
	if *cluster.StorageEncrypted {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBClusterEndpoint returns the cluster endpoint
func getDBClusterEndpoint(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).Endpoint), nil
}

// getDBClusterEngine returns the database engine
func getDBClusterEngine(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).Engine), nil
}

// getDBClusterEngineLifecycleSupport returns the engine lifecycle support
func getDBClusterEngineLifecycleSupport(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).EngineLifecycleSupport), nil
}

// getDBClusterEngineMode returns the engine mode
func getDBClusterEngineMode(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).EngineMode), nil
}

// getDBClusterEngineVersion returns the engine version
func getDBClusterEngineVersion(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).EngineVersion), nil
}

// getDBClusterGlobalClusterIdentifier returns the global cluster identifier
func getDBClusterGlobalClusterIdentifier(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).GlobalClusterIdentifier), nil
}

// getDBClusterGlobalWriteForwarding returns the global write forwarding status
func getDBClusterGlobalWriteForwarding(instance any) (string, error) {
	return string(instance.(types.DBCluster).GlobalWriteForwardingStatus), nil
}

// getDBClusterHostedZoneID returns the hosted zone ID
func getDBClusterHostedZoneID(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).HostedZoneId), nil
}

// getDBClusterHTTPEndpointEnabled returns the HTTP endpoint enabled setting
func getDBClusterHTTPEndpointEnabled(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	if cluster.HttpEndpointEnabled == nil {
		return "", nil
	}
	if *cluster.HttpEndpointEnabled {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBClusterIAMDatabaseAuthentication returns the IAM database authentication setting
func getDBClusterIAMDatabaseAuthentication(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	if cluster.IAMDatabaseAuthenticationEnabled == nil {
		return "", nil
	}
	if *cluster.IAMDatabaseAuthenticationEnabled {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBClusterIOPS returns the IOPS
func getDBClusterIOPS(instance any) (string, error) {
	iops := instance.(types.DBCluster).Iops
	if iops == nil {
		return "", nil
	}
	return fmt.Sprintf("%d", *iops), nil
}

// getDBClusterKMSKey returns the KMS key ID
func getDBClusterKMSKey(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).KmsKeyId), nil
}

// getDBClusterLatestRestorableTime returns the latest restorable time
func getDBClusterLatestRestorableTime(instance any) (string, error) {
	restoreTime := instance.(types.DBCluster).LatestRestorableTime
	if restoreTime == nil {
		return "", nil
	}
	return format.TimeToStringOrDefault(restoreTime, ""), nil
}

// getDBClusterLocalWriteForwarding returns the local write forwarding status
func getDBClusterLocalWriteForwarding(instance any) (string, error) {
	return string(instance.(types.DBCluster).LocalWriteForwardingStatus), nil
}

// getDBClusterMasterUsername returns the master username
func getDBClusterMasterUsername(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).MasterUsername), nil
}

// getDBClusterMonitoringInterval returns the monitoring interval
func getDBClusterMonitoringInterval(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	if cluster.MonitoringInterval == nil {
		return "Disabled", nil
	}
	if *cluster.MonitoringInterval == 0 {
		return "Disabled", nil
	}
	return strconv.Itoa(int(*cluster.MonitoringInterval)), nil
}

// getDBClusterMonitoringRole returns the monitoring role ARN
func getDBClusterMonitoringRole(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	roleArn := aws.ToString(cluster.MonitoringRoleArn)
	if roleArn == "" {
		return "Not configured", nil
	}
	return roleArn, nil
}

// getDBClusterMultiAZ returns the Multi-AZ setting
func getDBClusterMultiAZ(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	if cluster.MultiAZ == nil {
		return "", nil
	}
	if *cluster.MultiAZ {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBClusterNetworkType returns the network type
func getDBClusterNetworkType(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).NetworkType), nil
}

// getDBClusterOptionGroup returns the option group
func getDBClusterOptionGroup(instance any) (string, error) {
	optionGroups := instance.(types.DBCluster).DBClusterOptionGroupMemberships
	if len(optionGroups) == 0 {
		return "", nil
	}

	var groups []string
	for _, group := range optionGroups {
		groups = append(groups, aws.ToString(group.DBClusterOptionGroupName))
	}
	return strings.Join(groups, ", "), nil
}

// getDBClusterParameterGroup returns the parameter group
func getDBClusterParameterGroup(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).DBClusterParameterGroup), nil
}

// getDBClusterPendingModifications returns pending modifications
func getDBClusterPendingModifications(instance any) (string, error) {
	pending := instance.(types.DBCluster).PendingModifiedValues
	if pending == nil {
		return "None", nil
	}

	var modifications []string
	if pending.EngineVersion != nil {
		modifications = append(modifications, fmt.Sprintf("Engine Version: %s", aws.ToString(pending.EngineVersion)))
	}
	if pending.MasterUserPassword != nil {
		modifications = append(modifications, "Master User Password: ***")
	}

	if len(modifications) == 0 {
		return "None", nil
	}
	return strings.Join(modifications, ", "), nil
}

// getDBClusterPerformanceInsights returns the Performance Insights setting
func getDBClusterPerformanceInsights(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	if cluster.PerformanceInsightsEnabled == nil {
		return "Disabled", nil
	}
	if *cluster.PerformanceInsightsEnabled {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBClusterPerformanceInsightsKMSKey returns the Performance Insights KMS key
func getDBClusterPerformanceInsightsKMSKey(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).PerformanceInsightsKMSKeyId), nil
}

// getDBClusterPerformanceInsightsRetention returns the Performance Insights retention period
func getDBClusterPerformanceInsightsRetention(instance any) (string, error) {
	retention := instance.(types.DBCluster).PerformanceInsightsRetentionPeriod
	if retention == nil {
		return "", nil
	}
	return fmt.Sprintf("%d days", *retention), nil
}

// getDBClusterPort returns the port
func getDBClusterPort(instance any) (string, error) {
	port := instance.(types.DBCluster).Port
	if port == nil {
		return "", nil
	}
	return fmt.Sprintf("%d", *port), nil
}

// getDBClusterPreferredBackupWindow returns the preferred backup window
func getDBClusterPreferredBackupWindow(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).PreferredBackupWindow), nil
}

// getDBClusterMaintenanceWindow returns the maintenance window
func getDBClusterMaintenanceWindow(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).PreferredMaintenanceWindow), nil
}

// getDBClusterPubliclyAccessible returns the publicly accessible setting
func getDBClusterPubliclyAccessible(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	if cluster.PubliclyAccessible == nil {
		return "Disabled", nil
	}
	if *cluster.PubliclyAccessible {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBClusterReadReplicaIdentifiers returns the read replica identifiers
func getDBClusterReadReplicaIdentifiers(instance any) (string, error) {
	replicas := instance.(types.DBCluster).ReadReplicaIdentifiers
	if len(replicas) == 0 {
		return "", nil
	}
	return strings.Join(replicas, ", "), nil
}

// getDBClusterReaderEndpoint returns the reader endpoint
func getDBClusterReaderEndpoint(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	endpoint := aws.ToString(cluster.ReaderEndpoint)
	if endpoint == "" {
		return "Not configured", nil
	}
	return endpoint, nil
}

// getDBClusterReplicationSourceIdentifier returns the replication source identifier
func getDBClusterReplicationSourceIdentifier(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).ReplicationSourceIdentifier), nil
}

// getDBClusterResourceID returns the resource ID
func getDBClusterResourceID(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).DbClusterResourceId), nil
}

// getDBClusterScalingConfiguration returns the scaling configuration
func getDBClusterScalingConfiguration(instance any) (string, error) {
	scaling := instance.(types.DBCluster).ScalingConfigurationInfo
	if scaling == nil {
		return "", nil
	}

	var config []string
	if scaling.MinCapacity != nil {
		config = append(config, fmt.Sprintf("Min: %d", *scaling.MinCapacity))
	}
	if scaling.MaxCapacity != nil {
		config = append(config, fmt.Sprintf("Max: %d", *scaling.MaxCapacity))
	}
	if scaling.AutoPause != nil {
		if *scaling.AutoPause {
			config = append(config, "Auto Pause: Enabled")
		} else {
			config = append(config, "Auto Pause: Disabled")
		}
	}

	if len(config) == 0 {
		return "", nil
	}
	return strings.Join(config, ", "), nil
}

// getDBClusterServerlessV2PlatformVersion returns the Serverless V2 platform version
func getDBClusterServerlessV2PlatformVersion(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).ServerlessV2PlatformVersion), nil
}

// getDBClusterServerlessV2ScalingConfig returns the Serverless V2 scaling configuration
func getDBClusterServerlessV2ScalingConfig(instance any) (string, error) {
	scaling := instance.(types.DBCluster).ServerlessV2ScalingConfiguration
	if scaling == nil {
		return "", nil
	}

	var config []string
	if scaling.MinCapacity != nil {
		config = append(config, fmt.Sprintf("Min: %.1f", *scaling.MinCapacity))
	}
	if scaling.MaxCapacity != nil {
		config = append(config, fmt.Sprintf("Max: %.1f", *scaling.MaxCapacity))
	}

	if len(config) == 0 {
		return "", nil
	}
	return strings.Join(config, ", "), nil
}

// getDBClusterStatus returns the cluster status
func getDBClusterStatus(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).Status), nil
}

// getDBClusterStorageEncrypted returns the storage encrypted setting
func getDBClusterStorageEncrypted(instance any) (string, error) {
	cluster := instance.(types.DBCluster)
	if cluster.StorageEncrypted == nil {
		return "", nil
	}
	if *cluster.StorageEncrypted {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBClusterStorageThroughput returns the storage throughput
func getDBClusterStorageThroughput(instance any) (string, error) {
	throughput := instance.(types.DBCluster).StorageThroughput
	if throughput == nil {
		return "", nil
	}
	return fmt.Sprintf("%d MiB/s", *throughput), nil
}

// getDBClusterStorageType returns the storage type
func getDBClusterStorageType(instance any) (string, error) {
	return aws.ToString(instance.(types.DBCluster).StorageType), nil
}

// getDBClusterVPCSecurityGroups returns the VPC security groups
func getDBClusterVPCSecurityGroups(instance any) (string, error) {
	groups := instance.(types.DBCluster).VpcSecurityGroups
	if len(groups) == 0 {
		return "", nil
	}

	var groupIds []string
	for _, group := range groups {
		groupIds = append(groupIds, aws.ToString(group.VpcSecurityGroupId))
	}
	return strings.Join(groupIds, ", "), nil
}

// Clusters context - stored globally to support role calculation
var clustersContext []types.DBCluster
