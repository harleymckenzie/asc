package rds

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// RDS DB snapshot field getters (instance snapshots)
var dbSnapshotFieldValueGetters = map[string]FieldValueGetter{
	"Snapshot Identifier": getDBSnapshotIdentifier,
	"Source Identifier":   getDBSnapshotSourceIdentifier,
	"Type":                getDBSnapshotType,
	"Engine":              getDBSnapshotEngine,
	"Engine Version":      getDBSnapshotEngineVersion,
	"Status":              getDBSnapshotStatus,
	"Storage":             getDBSnapshotStorage,
	"Encryption":          getDBSnapshotEncryption,
	"Created Time":        getDBSnapshotCreatedTime,
	"Availability Zone":   getDBSnapshotAvailabilityZone,
	"Port":                getDBSnapshotPort,
	"VPC ID":              getDBSnapshotVPCID,
}

// RDS DB cluster snapshot field getters
var dbClusterSnapshotFieldValueGetters = map[string]FieldValueGetter{
	"Snapshot Identifier": getDBClusterSnapshotIdentifier,
	"Source Identifier":   getDBClusterSnapshotSourceIdentifier,
	"Type":                getDBClusterSnapshotType,
	"Engine":              getDBClusterSnapshotEngine,
	"Engine Version":      getDBClusterSnapshotEngineVersion,
	"Status":              getDBClusterSnapshotStatus,
	"Storage":             getDBClusterSnapshotStorage,
	"Encryption":          getDBClusterSnapshotEncryption,
	"Created Time":        getDBClusterSnapshotCreatedTime,
	"Availability Zone":   getDBClusterSnapshotAvailabilityZones,
	"Port":                getDBClusterSnapshotPort,
	"VPC ID":              getDBClusterSnapshotVPCID,
}

// getDBSnapshotFieldValue returns the value of a field for an RDS DB instance snapshot
func getDBSnapshotFieldValue(fieldName string, snapshot types.DBSnapshot) (string, error) {
	if getter, exists := dbSnapshotFieldValueGetters[fieldName]; exists {
		return getter(snapshot)
	}
	return "", fmt.Errorf("field %s not found in dbSnapshotFieldValueGetters", fieldName)
}

// getDBClusterSnapshotFieldValue returns the value of a field for an RDS DB cluster snapshot
func getDBClusterSnapshotFieldValue(fieldName string, snapshot types.DBClusterSnapshot) (string, error) {
	if getter, exists := dbClusterSnapshotFieldValueGetters[fieldName]; exists {
		return getter(snapshot)
	}
	return "", fmt.Errorf("field %s not found in dbClusterSnapshotFieldValueGetters", fieldName)
}

// -----------------------------------------------------------------------------
// DB instance snapshot field getters
// -----------------------------------------------------------------------------

func getDBSnapshotIdentifier(instance any) (string, error) {
	return aws.ToString(instance.(types.DBSnapshot).DBSnapshotIdentifier), nil
}

func getDBSnapshotSourceIdentifier(instance any) (string, error) {
	return aws.ToString(instance.(types.DBSnapshot).DBInstanceIdentifier), nil
}

func getDBSnapshotType(instance any) (string, error) {
	return aws.ToString(instance.(types.DBSnapshot).SnapshotType), nil
}

func getDBSnapshotEngine(instance any) (string, error) {
	return aws.ToString(instance.(types.DBSnapshot).Engine), nil
}

func getDBSnapshotEngineVersion(instance any) (string, error) {
	return aws.ToString(instance.(types.DBSnapshot).EngineVersion), nil
}

func getDBSnapshotStatus(instance any) (string, error) {
	return format.Status(aws.ToString(instance.(types.DBSnapshot).Status)), nil
}

func getDBSnapshotStorage(instance any) (string, error) {
	snapshot := instance.(types.DBSnapshot)
	if snapshot.AllocatedStorage == nil || *snapshot.AllocatedStorage == 0 {
		return "", nil
	}
	return strconv.Itoa(int(*snapshot.AllocatedStorage)) + " GB", nil
}

func getDBSnapshotEncryption(instance any) (string, error) {
	return format.BoolToLabel(instance.(types.DBSnapshot).Encrypted, "Enabled", "Disabled"), nil
}

func getDBSnapshotCreatedTime(instance any) (string, error) {
	return format.TimeToStringOrEmpty(instance.(types.DBSnapshot).SnapshotCreateTime), nil
}

func getDBSnapshotAvailabilityZone(instance any) (string, error) {
	return aws.ToString(instance.(types.DBSnapshot).AvailabilityZone), nil
}

func getDBSnapshotPort(instance any) (string, error) {
	snapshot := instance.(types.DBSnapshot)
	if snapshot.Port == nil || *snapshot.Port == 0 {
		return "", nil
	}
	return strconv.Itoa(int(*snapshot.Port)), nil
}

func getDBSnapshotVPCID(instance any) (string, error) {
	return aws.ToString(instance.(types.DBSnapshot).VpcId), nil
}

// -----------------------------------------------------------------------------
// DB cluster snapshot field getters
// -----------------------------------------------------------------------------

func getDBClusterSnapshotIdentifier(instance any) (string, error) {
	return aws.ToString(instance.(types.DBClusterSnapshot).DBClusterSnapshotIdentifier), nil
}

func getDBClusterSnapshotSourceIdentifier(instance any) (string, error) {
	return aws.ToString(instance.(types.DBClusterSnapshot).DBClusterIdentifier), nil
}

func getDBClusterSnapshotType(instance any) (string, error) {
	return aws.ToString(instance.(types.DBClusterSnapshot).SnapshotType), nil
}

func getDBClusterSnapshotEngine(instance any) (string, error) {
	return aws.ToString(instance.(types.DBClusterSnapshot).Engine), nil
}

func getDBClusterSnapshotEngineVersion(instance any) (string, error) {
	return aws.ToString(instance.(types.DBClusterSnapshot).EngineVersion), nil
}

func getDBClusterSnapshotStatus(instance any) (string, error) {
	return format.Status(aws.ToString(instance.(types.DBClusterSnapshot).Status)), nil
}

func getDBClusterSnapshotStorage(instance any) (string, error) {
	snapshot := instance.(types.DBClusterSnapshot)
	if snapshot.AllocatedStorage == nil || *snapshot.AllocatedStorage == 0 {
		return "", nil
	}
	return strconv.Itoa(int(*snapshot.AllocatedStorage)) + " GB", nil
}

func getDBClusterSnapshotEncryption(instance any) (string, error) {
	return format.BoolToLabel(instance.(types.DBClusterSnapshot).StorageEncrypted, "Enabled", "Disabled"), nil
}

func getDBClusterSnapshotCreatedTime(instance any) (string, error) {
	return format.TimeToStringOrEmpty(instance.(types.DBClusterSnapshot).SnapshotCreateTime), nil
}

func getDBClusterSnapshotAvailabilityZones(instance any) (string, error) {
	return strings.Join(instance.(types.DBClusterSnapshot).AvailabilityZones, "\n"), nil
}

func getDBClusterSnapshotPort(instance any) (string, error) {
	snapshot := instance.(types.DBClusterSnapshot)
	if snapshot.Port == nil || *snapshot.Port == 0 {
		return "", nil
	}
	return strconv.Itoa(int(*snapshot.Port)), nil
}

func getDBClusterSnapshotVPCID(instance any) (string, error) {
	return aws.ToString(instance.(types.DBClusterSnapshot).VpcId), nil
}
