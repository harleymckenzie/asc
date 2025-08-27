package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

var imageFieldValueGetters = map[string]FieldValueGetter{
	"Allowed Image":             getImageAllowedImage,
	"AMI ID":                    getImageAMIID,
	"AMI Name":                  getImageAMIName,
	"Architecture":              getImageArchitecture,
	"Block Devices":             getImageBlockDevices,
	"Boot Mode":                 getImageBootMode,
	"Creation Date":             getImageCreationDate,
	"Deprecation Time":          getImageDeprecationTime,
	"Deregistration Protection": getImageDeregistrationProtection,
	"Description":               getImageDescription,
	"Image Type":                getImageType,
	"Kernel ID":                 getImageKernelID,
	"Owner":                     getImageOwner,
	"Platform":                  getImagePlatform,
	"Product Codes":             getImageProductCodes,
	"RAM Disk ID":               getImageRAMDiskID,
	"Root Device Name":          getImageRootDeviceName,
	"Root Device Type":          getImageRootDeviceType,
	"Source":                    getImageSource,
	"Source AMI ID":             getImageSourceAMIID,
	"Source AMI Region":         getImageSourceAMIRegion,
	"State Reason":              getImageStateReason,
	"Status":                    getImageStatus,
	"Usage Operation":           getImageUsageOperation,
	"Virtualization":            getImageVirtualizationType,
	"Visibility":                getImageVisibility,
}

// getImageFieldValue returns the value of a field for an EC2 image
func getImageFieldValue(fieldName string, image types.Image) (string, error) {
	if getter, exists := imageFieldValueGetters[fieldName]; exists {
		value, err := getter(image)
		if err != nil {
			return "", fmt.Errorf("failed to get field value for %s: %w", fieldName, err)
		}
		return value, nil
	}
	return "", fmt.Errorf("field %s not found in image fieldValueGetters", fieldName)
}

// Individual field value getters

func getImageAllowedImage(image any) (string, error) {
	return format.BoolToLabel(image.(types.Image).ImageAllowed, "Yes", "No"), nil
}

func getImageAMIID(image any) (string, error) {
	return aws.ToString(image.(types.Image).ImageId), nil
}

func getImageAMIName(image any) (string, error) {
	return aws.ToString(image.(types.Image).Name), nil
}

func getImageArchitecture(image any) (string, error) {
	return string(image.(types.Image).Architecture), nil
}

func getImageBlockDevices(image any) (string, error) {
	return getBlockDevices(image.(types.Image).BlockDeviceMappings), nil
}

func getImageBootMode(image any) (string, error) {
	return string(image.(types.Image).BootMode), nil
}

func getImageCreationDate(image any) (string, error) {
	if image.(types.Image).CreationDate == nil {
		return "", nil
	}
	return format.StringOrEmpty(image.(types.Image).CreationDate), nil
}

func getImageDeprecationTime(image any) (string, error) {
	if image.(types.Image).DeprecationTime == nil {
		return "", nil
	}
	return aws.ToString(image.(types.Image).DeprecationTime), nil
}

func getImageDeregistrationProtection(image any) (string, error) {
	return aws.ToString(image.(types.Image).DeregistrationProtection), nil
}

func getImageDescription(image any) (string, error) {
	return aws.ToString(image.(types.Image).Description), nil
}

func getImageType(image any) (string, error) {
	return string(image.(types.Image).ImageType), nil
}

func getImageKernelID(image any) (string, error) {
	return aws.ToString(image.(types.Image).KernelId), nil
}

func getImageOwner(image any) (string, error) {
	return aws.ToString(image.(types.Image).OwnerId), nil
}

func getImagePlatform(image any) (string, error) {
	return aws.ToString(image.(types.Image).PlatformDetails), nil
}

func getImageProductCodes(image any) (string, error) {
	return getProductCodes(image.(types.Image).ProductCodes), nil
}

func getImageRAMDiskID(image any) (string, error) {
	return aws.ToString(image.(types.Image).RamdiskId), nil
}

func getImageRootDeviceName(image any) (string, error) {
	return aws.ToString(image.(types.Image).RootDeviceName), nil
}

func getImageRootDeviceType(image any) (string, error) {
	return string(image.(types.Image).RootDeviceType), nil
}

func getImageSource(image any) (string, error) {
	return aws.ToString(image.(types.Image).ImageLocation), nil
}

func getImageSourceAMIID(image any) (string, error) {
	return aws.ToString(image.(types.Image).SourceImageId), nil
}

func getImageSourceAMIRegion(image any) (string, error) {
	return aws.ToString(image.(types.Image).SourceImageRegion), nil
}

func getImageStateReason(image any) (string, error) {
	if image.(types.Image).StateReason == nil || image.(types.Image).StateReason.Message == nil {
		return "", nil
	}
	return aws.ToString(image.(types.Image).StateReason.Message), nil
}

func getImageStatus(image any) (string, error) {
	return format.StatusOrDefault(string(image.(types.Image).State), ""), nil
}

func getImageUsageOperation(image any) (string, error) {
	return aws.ToString(image.(types.Image).UsageOperation), nil
}

func getImageVirtualizationType(image any) (string, error) {
	return string(image.(types.Image).VirtualizationType), nil
}

func getImageVisibility(image any) (string, error) {
	if image.(types.Image).Public == nil {
		return "", nil
	}
	return format.BoolToLabel(image.(types.Image).Public, "Public", "Private"), nil
}
