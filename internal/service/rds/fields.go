package rds

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type FieldValueGetter func(instance any) (string, error)

// RDS DB Instance field getters
var dbInstanceFieldValueGetters = map[string]FieldValueGetter{
	"Auto Minor Version Upgrade": getDBInstanceAutoMinorVersionUpgrade,
	"Availability Zone":          getDBInstanceAvailabilityZone,
	"ARN":                        getDBInstanceARN,
	"AWS KMS Key":                getDBInstanceKMSKey,
	"Certificate Authority":      getDBInstanceCertificateAuthority,
	"Certificate Expiry Date":    getDBInstanceCertificateExpiryDate,
	"Class":                      getDBInstanceClass,
	"Cluster Identifier":         getDBInstanceClusterID,
	"Created Time":               getDBInstanceCreatedTime,
	"DB Name":                    getDBInstanceDBName,
	"Encryption":                 getDBInstanceEncryption,
	"Endpoint":                   getDBInstanceEndpoint,
	"Engine":                     getDBInstanceEngine,
	"Engine Version":             getDBInstanceEngineVersion,
	"Failover Priority":          getDBInstanceFailoverPriority,
	"Identifier":                 getDBInstanceID,
	"Maintenance Window":         getDBInstanceMaintenanceWindow,
	"Monitoring Interval":        getDBInstanceMonitoringInterval,
	"Monitoring Role":            getDBInstanceMonitoringRole,
	"Network Type":               getDBInstanceNetworkType,
	"Option Group":               getDBInstanceOptionGroup,
	"Parameter Group":            getDBInstanceParameterGroup,
	"Pending Modifications":      getDBInstancePendingModifications,
	"Performance Insights":       getDBInstancePerformanceInsights,
	"Port":                       getDBInstancePort,
	"Publicly Accessible":        getDBInstancePubliclyAccessible,
	"RDS Extended Support":       getDBInstanceRDSExtendedSupport,
	"Resource ID":                getDBInstanceResourceID,
	"Role":                       getDBInstanceRoleField,
	"Security Group(s)":          getDBInstanceSecurityGroups,
	"Size":                       getDBInstanceSize,
	"Status":                     getDBInstanceStatus,
	"Storage Type":               getDBInstanceStorageType,
	"Subnet Group":               getDBInstanceSubnetGroup,
	"Subnets":                    getDBInstanceSubnets,
	"VPC ID":                     getDBInstanceVPCID,
}

// GetFieldValue returns the value of a field for the given instance.
func GetFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.DBInstance:
		return getDBInstanceFieldValue(fieldName, v)
	case types.DBCluster:
		return getDBClusterFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// getDBInstanceFieldValue returns the value of a field for an RDS DB instance
func getDBInstanceFieldValue(fieldName string, instance types.DBInstance) (string, error) {
	if getter, exists := dbInstanceFieldValueGetters[fieldName]; exists {
		return getter(instance)
	}
	return "", fmt.Errorf("field %s not found in dbInstanceFieldValueGetters", fieldName)
}

// GetTagValue returns the value of a tag for a DB instance or DB cluster
func GetTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.DBInstance:
		for _, tag := range v.TagList {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	case types.DBCluster:
		for _, tag := range v.TagList {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	default:
		return "", fmt.Errorf("unsupported instance type for tags: %T", instance)
	}
	return "", nil
}

// -----------------------------------------------------------------------------
// DB Instance field getters
// -----------------------------------------------------------------------------

// getDBInstanceClusterID returns the cluster identifier if this instance is part of a cluster
func getDBInstanceClusterID(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.DBClusterIdentifier != nil {
		return aws.ToString(dbInstance.DBClusterIdentifier), nil
	}
	return "-", nil
}

// getDBInstanceID returns the database instance identifier
func getDBInstanceID(instance any) (string, error) {
	return aws.ToString(instance.(types.DBInstance).DBInstanceIdentifier), nil
}

// getDBInstanceStatus returns the current status of the database instance
func getDBInstanceStatus(instance any) (string, error) {
	return format.Status(aws.ToString(instance.(types.DBInstance).DBInstanceStatus)), nil
}

// getDBInstanceEngine returns the database engine (e.g., mysql, postgres, aurora)
func getDBInstanceEngine(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.Engine == nil {
		return "", nil
	}
	return string(*dbInstance.Engine), nil
}

// getDBInstanceEngineVersion returns the version of the database engine
func getDBInstanceEngineVersion(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.EngineVersion == nil {
		return "", nil
	}
	return string(*dbInstance.EngineVersion), nil
}

// getDBInstanceSize returns the instance class/size of the database
func getDBInstanceSize(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.DBInstanceClass == nil {
		return "", nil
	}
	return string(*dbInstance.DBInstanceClass), nil
}

// getDBInstanceRoleField returns the role of the instance within a cluster (Primary, Replica, etc.)
func getDBInstanceRoleField(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	return calculateDBInstanceRole(dbInstance, clustersContext), nil
}

// getDBInstanceEndpoint returns the endpoint address of the database instance
func getDBInstanceEndpoint(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.Endpoint == nil || dbInstance.Endpoint.Address == nil {
		return "", nil
	}
	return aws.ToString(dbInstance.Endpoint.Address), nil
}

// getDBInstanceAutoMinorVersionUpgrade returns whether auto minor version upgrade is enabled
func getDBInstanceAutoMinorVersionUpgrade(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if *dbInstance.AutoMinorVersionUpgrade {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBInstanceAvailabilityZone returns the availability zone of the database instance
func getDBInstanceAvailabilityZone(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	return aws.ToString(dbInstance.AvailabilityZone), nil
}

// getDBInstanceARN returns the ARN of the database instance
func getDBInstanceARN(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	return aws.ToString(dbInstance.DBInstanceArn), nil
}

// getDBInstanceKMSKey returns the KMS key ID used for encryption
func getDBInstanceKMSKey(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	return aws.ToString(dbInstance.KmsKeyId), nil
}

// getDBInstanceCertificateAuthority returns the certificate authority identifier
func getDBInstanceCertificateAuthority(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.CertificateDetails == nil {
		return "", nil
	}
	return aws.ToString(dbInstance.CertificateDetails.CAIdentifier), nil
}

// getDBInstanceCertificateExpiryDate returns the certificate expiry date
func getDBInstanceCertificateExpiryDate(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.CertificateDetails == nil {
		return "", nil
	}
	return format.TimeToStringOrEmpty(dbInstance.CertificateDetails.ValidTill), nil
}

// getDBInstanceClass returns the instance class of the database
func getDBInstanceClass(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.DBInstanceClass == nil {
		return "", nil
	}
	return string(*dbInstance.DBInstanceClass), nil
}

// getDBInstanceCreatedTime returns the creation time of the database instance
func getDBInstanceCreatedTime(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	return format.TimeToStringOrEmpty(dbInstance.InstanceCreateTime), nil
}

// getDBInstanceDBName returns the database name
func getDBInstanceDBName(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	return aws.ToString(dbInstance.DBName), nil
}

// getDBInstanceEncryption returns whether storage encryption is enabled
func getDBInstanceEncryption(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if *dbInstance.StorageEncrypted {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBInstanceFailoverPriority returns the failover priority (promotion tier)
func getDBInstanceFailoverPriority(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.PromotionTier != nil {
		return strconv.Itoa(int(*dbInstance.PromotionTier)), nil
	}
	return "-", nil
}

// getDBInstanceMaintenanceWindow returns the preferred maintenance window
func getDBInstanceMaintenanceWindow(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	return aws.ToString(dbInstance.PreferredMaintenanceWindow), nil
}

// getDBInstanceMonitoringInterval returns the monitoring interval in seconds
func getDBInstanceMonitoringInterval(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.MonitoringInterval == nil {
		return "Disabled", nil
	}
	if *dbInstance.MonitoringInterval == 0 {
		return "Disabled", nil
	}
	return strconv.Itoa(int(*dbInstance.MonitoringInterval)), nil
}

// getDBInstanceMonitoringRole returns the monitoring role ARN
func getDBInstanceMonitoringRole(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	roleArn := aws.ToString(dbInstance.MonitoringRoleArn)
	if roleArn == "" {
		return "Not configured", nil
	}
	return aws.ToString(dbInstance.MonitoringRoleArn), nil
}

// getDBInstanceNetworkType returns the network type (VPC ID)
func getDBInstanceNetworkType(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.DBSubnetGroup == nil {
		return "Not configured", nil
	}
	return aws.ToString(dbInstance.DBSubnetGroup.VpcId), nil
}

// getDBInstanceOptionGroup returns the option group names
func getDBInstanceOptionGroup(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if len(dbInstance.OptionGroupMemberships) == 0 {
		return "-", nil
	}

	var optionGroups []string
	for _, og := range dbInstance.OptionGroupMemberships {
		optionGroups = append(optionGroups, aws.ToString(og.OptionGroupName))
	}

	// Join with commas for cleaner display
	return strings.Join(optionGroups, ", "), nil
}

// getDBInstanceParameterGroup returns the parameter group names
func getDBInstanceParameterGroup(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if len(dbInstance.DBParameterGroups) == 0 {
		return "-", nil
	}

	var parameterGroups []string
	for _, pg := range dbInstance.DBParameterGroups {
		parameterGroups = append(parameterGroups, aws.ToString(pg.DBParameterGroupName))
	}

	// Join with commas for cleaner display
	return strings.Join(parameterGroups, ", "), nil
}

// getDBInstancePendingModifications returns pending modifications information
func getDBInstancePendingModifications(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	pending := dbInstance.PendingModifiedValues

	// If no pending modifications, return "None"
	if pending == nil {
		return "None", nil
	}

	modifications := extractPendingModifications(pending)

	// If no specific modifications found, return "None"
	if len(modifications) == 0 {
		return "None", nil
	}

	// Return formatted list of modifications
	return fmt.Sprintf("%d pending: %s", len(modifications), fmt.Sprintf("%v", modifications)), nil
}

// extractPendingModifications extracts and formats pending modifications from the PendingModifiedValues struct
func extractPendingModifications(pending *types.PendingModifiedValues) []string {
	var modifications []string

	// Storage modifications
	if pending.AllocatedStorage != nil {
		modifications = append(modifications, fmt.Sprintf("Allocated Storage: %d", *pending.AllocatedStorage))
	}
	if pending.StorageType != nil {
		modifications = append(modifications, fmt.Sprintf("Storage Type: %s", *pending.StorageType))
	}
	if pending.StorageThroughput != nil {
		modifications = append(modifications, fmt.Sprintf("Storage Throughput: %d", *pending.StorageThroughput))
	}

	// Instance modifications
	if pending.DBInstanceClass != nil {
		modifications = append(modifications, fmt.Sprintf("Instance Class: %s", *pending.DBInstanceClass))
	}
	if pending.DBInstanceIdentifier != nil {
		modifications = append(modifications, fmt.Sprintf("Instance ID: %s", *pending.DBInstanceIdentifier))
	}
	if pending.EngineVersion != nil {
		modifications = append(modifications, fmt.Sprintf("Engine Version: %s", *pending.EngineVersion))
	}

	// Security modifications
	if pending.CACertificateIdentifier != nil {
		modifications = append(modifications, fmt.Sprintf("CA Certificate: %s", *pending.CACertificateIdentifier))
	}
	if pending.IAMDatabaseAuthenticationEnabled != nil {
		modifications = append(modifications, fmt.Sprintf("IAM Auth: %t", *pending.IAMDatabaseAuthenticationEnabled))
	}
	if pending.MasterUserPassword != nil {
		modifications = append(modifications, "Master Password: [CHANGED]")
	}

	// Performance modifications
	if pending.Iops != nil {
		modifications = append(modifications, fmt.Sprintf("IOPS: %d", *pending.Iops))
	}
	if len(pending.ProcessorFeatures) > 0 {
		modifications = append(modifications, fmt.Sprintf("Processor Features: %v", pending.ProcessorFeatures))
	}
	if pending.MultiAZ != nil {
		modifications = append(modifications, fmt.Sprintf("Multi-AZ: %t", *pending.MultiAZ))
	}
	if pending.MultiTenant != nil {
		modifications = append(modifications, fmt.Sprintf("Multi-Tenant: %t", *pending.MultiTenant))
	}

	// Networking modifications
	if pending.Port != nil {
		modifications = append(modifications, fmt.Sprintf("Port: %d", *pending.Port))
	}
	if pending.DBSubnetGroupName != nil {
		modifications = append(modifications, fmt.Sprintf("Subnet Group: %s", *pending.DBSubnetGroupName))
	}

	// Backup modifications
	if pending.BackupRetentionPeriod != nil {
		modifications = append(modifications, fmt.Sprintf("Backup Retention: %d days", *pending.BackupRetentionPeriod))
	}

	// Monitoring modifications
	if pending.PendingCloudwatchLogsExports != nil && len(pending.PendingCloudwatchLogsExports.LogTypesToEnable) > 0 {
		modifications = append(modifications, fmt.Sprintf("CloudWatch Logs: %v", pending.PendingCloudwatchLogsExports.LogTypesToEnable))
	}
	if pending.DedicatedLogVolume != nil {
		modifications = append(modifications, fmt.Sprintf("Dedicated Log Volume: %t", *pending.DedicatedLogVolume))
	}

	// Licensing modifications
	if pending.LicenseModel != nil {
		modifications = append(modifications, fmt.Sprintf("License Model: %s", *pending.LicenseModel))
	}

	return modifications
}

// getDBInstancePerformanceInsights returns whether Performance Insights is enabled
func getDBInstancePerformanceInsights(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if *dbInstance.PerformanceInsightsEnabled {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBInstancePort returns the port number of the database instance
func getDBInstancePort(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.Endpoint == nil || dbInstance.Endpoint.Port == nil {
		return "Not configured", nil
	}
	return strconv.Itoa(int(*dbInstance.Endpoint.Port)), nil
}

// getDBInstancePubliclyAccessible returns whether the instance is publicly accessible
func getDBInstancePubliclyAccessible(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if *dbInstance.PubliclyAccessible {
		return "Enabled", nil
	}
	return "Disabled", nil
}

// getDBInstanceRDSExtendedSupport returns the RDS Extended Support status
func getDBInstanceRDSExtendedSupport(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	return string(*dbInstance.EngineLifecycleSupport), nil
}

// getDBInstanceResourceID returns the resource ID of the database instance
func getDBInstanceResourceID(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	return aws.ToString(dbInstance.DbiResourceId), nil
}

// getDBInstanceSecurityGroups returns the security group names
func getDBInstanceSecurityGroups(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)

	var securityGroups []string

	// Check VPC security groups (modern, used by Aurora)
	for _, sg := range dbInstance.VpcSecurityGroups {
		securityGroups = append(securityGroups, aws.ToString(sg.VpcSecurityGroupId))
	}

	// Check DB security groups (legacy, used by older instances)
	for _, sg := range dbInstance.DBSecurityGroups {
		securityGroups = append(securityGroups, aws.ToString(sg.DBSecurityGroupName))
	}

	if len(securityGroups) == 0 {
		return "-", nil
	}

	// Join with commas for cleaner display
	return strings.Join(securityGroups, ", "), nil
}

// getDBInstanceStorageType returns the storage type of the database instance
func getDBInstanceStorageType(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	return aws.ToString(dbInstance.StorageType), nil
}

// getDBInstanceSubnetGroup returns the subnet group name
func getDBInstanceSubnetGroup(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.DBSubnetGroup != nil {
		return aws.ToString(dbInstance.DBSubnetGroup.DBSubnetGroupName), nil
	}
	return "-", nil
}

// getDBInstanceSubnets returns the subnet identifiers
func getDBInstanceSubnets(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.DBSubnetGroup == nil || len(dbInstance.DBSubnetGroup.Subnets) == 0 {
		return "Not configured", nil
	}
	return aws.ToString(dbInstance.DBSubnetGroup.Subnets[0].SubnetIdentifier), nil
}

// getDBInstanceVPCID returns the VPC ID of the database instance
func getDBInstanceVPCID(instance any) (string, error) {
	dbInstance := instance.(types.DBInstance)
	if dbInstance.DBSubnetGroup == nil {
		return "Not configured", nil
	}
	return aws.ToString(dbInstance.DBSubnetGroup.VpcId), nil
}

// -----------------------------------------------------------------------------
// Helper functions
// -----------------------------------------------------------------------------

// calculateDBInstanceRole calculates the role of the RDS instance within a cluster
func calculateDBInstanceRole(instance types.DBInstance, clusters []types.DBCluster) string {
	// If ReadReplicaSourceDBInstanceIdentifier is set, then this is a replica
	if instance.ReadReplicaSourceDBInstanceIdentifier != nil {
		return "Replica"
	}

	// If ReadReplicaDBInstanceIdentifiers is set, then this is a primary
	if len(instance.ReadReplicaDBInstanceIdentifiers) > 0 {
		return "Primary"
	}

	// If not part of a cluster, return None
	if instance.DBClusterIdentifier == nil {
		return "None"
	}

	// Check cluster membership to determine if this is a writer or reader
	for _, cluster := range clusters {
		for _, member := range cluster.DBClusterMembers {
			if aws.ToString(member.DBInstanceIdentifier) == aws.ToString(instance.DBInstanceIdentifier) {
				if member.IsClusterWriter != nil && *member.IsClusterWriter {
					return "Writer"
				}
				return "Reader"
			}
		}
	}

	return "Unknown"
}
