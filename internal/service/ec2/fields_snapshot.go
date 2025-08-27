package ec2

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// snapshotFieldValueGetters provides quick lookup for common snapshot fields.
// Each getter delegates to the central GetSnapshotAttributeValue ensuring a
// single source of truth in table.go.
var snapshotFieldValueGetters = map[string]FieldValueGetter{
	"Snapshot ID":            getSnapshotID,
	"Volume Size":            getSnapshotVolumeSize,
	"Description":            getSnapshotDescription,
	"Tier":                   getSnapshotTier,
	"State":                  getSnapshotState,
	"Started":                getSnapshotStarted,
	"Progress":               getSnapshotProgress,
	"Encryption":             getSnapshotEncryption,
	"Data Transfer Progress": getSnapshotDataTransferProgress,
	"KMS Key ID":             getSnapshotKMSKeyID,
	"Owner ID":               getSnapshotOwnerID,
	"Owner Alias":            getSnapshotOwnerAlias,
	"Source Volume":          getSnapshotSourceVolume,
	"Volume ID":              getSnapshotVolumeID,
	"Storage Tier":           getSnapshotStorageTier,
	"Restore Expiry Time":    getSnapshotRestoreExpiryTime,
}

// getSnapshotFieldValue returns the value of a field for an EC2 snapshot
func getSnapshotFieldValue(fieldName string, snapshot types.Snapshot) (string, error) {
	if getter, exists := snapshotFieldValueGetters[fieldName]; exists {
		return getter(snapshot)
	}
	return "", fmt.Errorf("field %s not found in snapshotFieldValueGetters", fieldName)
}

// Individual field value getters

func getSnapshotID(snapshot any) (string, error) {
	return aws.ToString(snapshot.(types.Snapshot).SnapshotId), nil
}

func getSnapshotVolumeSize(snapshot any) (string, error) {
	size := snapshot.(types.Snapshot).VolumeSize
	if size == nil {
		return "", nil
	}
	return fmt.Sprintf("%d GiB", *size), nil
}

func getSnapshotDescription(snapshot any) (string, error) {
	return aws.ToString(snapshot.(types.Snapshot).Description), nil
}

func getSnapshotTier(snapshot any) (string, error) {
	return string(snapshot.(types.Snapshot).StorageTier), nil
}

func getSnapshotState(snapshot any) (string, error) {
	return format.Status(string(snapshot.(types.Snapshot).State)), nil
}

func getSnapshotStarted(snapshot any) (string, error) {
	return snapshot.(types.Snapshot).StartTime.Format(time.RFC3339), nil
}

func getSnapshotProgress(snapshot any) (string, error) {
	return format.Status(aws.ToString(snapshot.(types.Snapshot).Progress)), nil
}

func getSnapshotEncryption(snapshot any) (string, error) {
	if snapshot.(types.Snapshot).Encrypted != nil && *snapshot.(types.Snapshot).Encrypted {
		return "Encrypted", nil
	}
	return "Not encrypted", nil
}

func getSnapshotDataTransferProgress(snapshot any) (string, error) {
	return aws.ToString(snapshot.(types.Snapshot).Progress), nil
}

func getSnapshotKMSKeyID(snapshot any) (string, error) {
	return aws.ToString(snapshot.(types.Snapshot).KmsKeyId), nil
}

func getSnapshotOwnerID(snapshot any) (string, error) {
	return aws.ToString(snapshot.(types.Snapshot).OwnerId), nil
}

func getSnapshotOwnerAlias(snapshot any) (string, error) {
	return aws.ToString(snapshot.(types.Snapshot).OwnerAlias), nil
}

func getSnapshotSourceVolume(snapshot any) (string, error) {
	return aws.ToString(snapshot.(types.Snapshot).VolumeId), nil
}

func getSnapshotVolumeID(snapshot any) (string, error) {
	return aws.ToString(snapshot.(types.Snapshot).VolumeId), nil
}

func getSnapshotStorageTier(snapshot any) (string, error) {
	return string(snapshot.(types.Snapshot).StorageTier), nil
}

func getSnapshotRestoreExpiryTime(snapshot any) (string, error) {
	if snapshot.(types.Snapshot).RestoreExpiryTime != nil {
		return snapshot.(types.Snapshot).RestoreExpiryTime.Format(time.RFC3339), nil
	}
	return "", nil
}
