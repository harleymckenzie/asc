package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// volumeFieldValueGetters maps field names to their getter functions for a
// types.Volume instance.  This file is intentionally self-contained so it can
// live on after we remove table.go.
var volumeFieldValueGetters = map[string]FieldValueGetter{
	// Core details
	"Volume ID":              getVolumeID,
	"Type":                   getVolumeType,
	"Size":                   getVolumeSize,
	"Size Raw":               getVolumeSizeRaw,
	"State":                  getVolumeState,
	"IOPS":                   getVolumeIOPS,
	"Throughput":             getVolumeThroughput,
	"Fast Snapshot Restored": getVolumeFSR,
	"Availability Zone":      getVolumeAZ,
	"Created":                getVolumeCreated,
	"Multi-Attach Enabled":   getVolumeMultiAttach,
	"Encryption":             getVolumeEncryption,
	"KMS Key ID":             getVolumeKMSID,

	// Associations
	"Snapshot ID":           getVolumeSnapshotID,
	"Associated Resource":   getVolumeAssociatedResource,
	"Attach Time":           getVolumeAttachTime,
	"Delete on Termination": getVolumeDeleteOnTermination,
	"Device":                getVolumeDevice,
	"Instance ID":           getVolumeInstanceID,
	"Attachment State":      getVolumeAttachmentState,
}

// getVolumeFieldValue returns the requested field value or an error if the
// field is unknown.
func getVolumeFieldValue(fieldName string, volume types.Volume) (string, error) {
	if getter, ok := volumeFieldValueGetters[fieldName]; ok {
		return getter(volume)
	}
	return "", fmt.Errorf("field %s not found in volumeFieldValueGetters", fieldName)
}

// -----------------------------------------------------------------------------
// Individual field getters
// -----------------------------------------------------------------------------

func getVolumeID(v any) (string, error) { return aws.ToString(v.(types.Volume).VolumeId), nil }

func getVolumeType(v any) (string, error) { return string(v.(types.Volume).VolumeType), nil }

func getVolumeSize(v any) (string, error) {
	size := format.Int32ToStringOrEmpty(v.(types.Volume).Size)
	if size == "" {
		return "", nil
	}
	return fmt.Sprintf("%s GiB", size), nil
}

func getVolumeSizeRaw(v any) (string, error) {
	return format.Int32ToStringOrEmpty(v.(types.Volume).Size), nil
}

func getVolumeState(v any) (string, error) { return format.Status(string(v.(types.Volume).State)), nil }

func getVolumeIOPS(v any) (string, error) {
	return format.Int32ToStringOrEmpty(v.(types.Volume).Iops), nil
}

func getVolumeThroughput(v any) (string, error) {
	vol := v.(types.Volume)
	if vol.Throughput == nil {
		return "-", nil
	}
	return fmt.Sprintf("%s MiB/s", format.Int32ToStringOrEmpty(vol.Throughput)), nil
}

func getVolumeFSR(v any) (string, error) {
	vol := v.(types.Volume)
	if vol.FastRestored == nil {
		return "", nil
	}
	return format.BoolToLabel(vol.FastRestored, "Yes", "No"), nil
}

func getVolumeAZ(v any) (string, error) { return aws.ToString(v.(types.Volume).AvailabilityZone), nil }

func getVolumeCreated(v any) (string, error) {
	return format.TimeToStringOrEmpty(v.(types.Volume).CreateTime), nil
}

func getVolumeMultiAttach(v any) (string, error) {
	return format.BoolToLabel(v.(types.Volume).MultiAttachEnabled, "Yes", "No"), nil
}

func getVolumeEncryption(v any) (string, error) {
	return format.BoolToLabel(v.(types.Volume).Encrypted, "Encrypted", "Not encrypted"), nil
}

func getVolumeKMSID(v any) (string, error) { return aws.ToString(v.(types.Volume).KmsKeyId), nil }

func getVolumeSnapshotID(v any) (string, error) {
	return aws.ToString(v.(types.Volume).SnapshotId), nil
}

func getVolumeAssociatedResource(v any) (string, error) {
	attachments := v.(types.Volume).Attachments
	if len(attachments) == 0 {
		return "", nil
	}
	return aws.ToString(attachments[0].InstanceId), nil
}

func getVolumeAttachTime(v any) (string, error) {
	attachments := v.(types.Volume).Attachments
	if len(attachments) == 0 {
		return "", nil
	}
	return format.TimeToStringOrEmpty(attachments[0].AttachTime), nil
}

func getVolumeDeleteOnTermination(v any) (string, error) {
	attachments := v.(types.Volume).Attachments
	if len(attachments) == 0 {
		return "", nil
	}
	return format.BoolToLabel(attachments[0].DeleteOnTermination, "Yes", "No"), nil
}

func getVolumeDevice(v any) (string, error) {
	attachments := v.(types.Volume).Attachments
	if len(attachments) == 0 {
		return "", nil
	}
	return aws.ToString(attachments[0].Device), nil
}

func getVolumeInstanceID(v any) (string, error) {
	attachments := v.(types.Volume).Attachments
	if len(attachments) == 0 {
		return "", nil
	}
	return aws.ToString(attachments[0].InstanceId), nil
}

func getVolumeAttachmentState(v any) (string, error) {
	attachments := v.(types.Volume).Attachments
	if len(attachments) == 0 {
		return "", nil
	}
	return format.Status(string(attachments[0].State)), nil
}
