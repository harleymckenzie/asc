package efs

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

var fileSystemFieldValueGetters = map[string]FieldValueGetter{
	"Name":              getFileSystemName,
	"File System ID":    getFileSystemID,
	"State":             getFileSystemState,
	"Size (Bytes)":      getFileSystemSize,
	"Mount Targets":     getFileSystemMountTargets,
	"Performance Mode":  getFileSystemPerformanceMode,
	"Throughput Mode":   getFileSystemThroughputMode,
	"Encrypted":              getFileSystemEncrypted,
	"Availability Zone":      getFileSystemAvailabilityZone,
	"Creation Time":          getFileSystemCreationTime,
	"ARN":                    getFileSystemARN,
	"Owner ID":               getFileSystemOwnerID,
	"KMS Key ID":             getFileSystemKMSKeyID,
	"Provisioned Throughput": getFileSystemProvisionedThroughput,
}

func getFileSystemName(instance any) (string, error) {
	return aws.ToString(instance.(types.FileSystemDescription).Name), nil
}

func getFileSystemID(instance any) (string, error) {
	return aws.ToString(instance.(types.FileSystemDescription).FileSystemId), nil
}

func getFileSystemState(instance any) (string, error) {
	return format.Status(string(instance.(types.FileSystemDescription).LifeCycleState)), nil
}

func getFileSystemSize(instance any) (string, error) {
	fs := instance.(types.FileSystemDescription)
	if fs.SizeInBytes != nil {
		return fmt.Sprintf("%d", fs.SizeInBytes.Value), nil
	}
	return "", nil
}

func getFileSystemMountTargets(instance any) (string, error) {
	return strconv.Itoa(int(instance.(types.FileSystemDescription).NumberOfMountTargets)), nil
}

func getFileSystemPerformanceMode(instance any) (string, error) {
	return string(instance.(types.FileSystemDescription).PerformanceMode), nil
}

func getFileSystemThroughputMode(instance any) (string, error) {
	return string(instance.(types.FileSystemDescription).ThroughputMode), nil
}

func getFileSystemEncrypted(instance any) (string, error) {
	encrypted := instance.(types.FileSystemDescription).Encrypted
	if encrypted != nil {
		if *encrypted {
			return "Yes", nil
		}
		return "No", nil
	}
	return "", nil
}

func getFileSystemAvailabilityZone(instance any) (string, error) {
	return aws.ToString(instance.(types.FileSystemDescription).AvailabilityZoneName), nil
}

func getFileSystemCreationTime(instance any) (string, error) {
	fs := instance.(types.FileSystemDescription)
	if fs.CreationTime != nil {
		return fs.CreationTime.Format(time.RFC3339), nil
	}
	return "", nil
}

func getFileSystemARN(instance any) (string, error) {
	return aws.ToString(instance.(types.FileSystemDescription).FileSystemArn), nil
}

func getFileSystemOwnerID(instance any) (string, error) {
	return aws.ToString(instance.(types.FileSystemDescription).OwnerId), nil
}

func getFileSystemKMSKeyID(instance any) (string, error) {
	return aws.ToString(instance.(types.FileSystemDescription).KmsKeyId), nil
}

func getFileSystemProvisionedThroughput(instance any) (string, error) {
	fs := instance.(types.FileSystemDescription)
	if fs.ProvisionedThroughputInMibps != nil {
		return fmt.Sprintf("%.1f MiBps", *fs.ProvisionedThroughputInMibps), nil
	}
	return "", nil
}
