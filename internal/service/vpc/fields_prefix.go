package vpc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type PrefixFieldValueGetter func(instance any) (string, error)

// Prefix List field getters
var prefixFieldValueGetters = map[string]PrefixFieldValueGetter{
	"Prefix List ID":   getPrefixListID,
	"Prefix List Name": getPrefixListName,
	"Max Entries":      getPrefixListMaxEntries,
	"Address Family":   getPrefixListAddressFamily,
	"State":            getPrefixListState,
	"Version":          getPrefixListVersion,
	"Prefix List ARN":  getPrefixListARN,
	"Owner":            getPrefixListOwner,
}

// GetPrefixFieldValue returns the value of a field for the given Prefix List instance.
func GetPrefixFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.ManagedPrefixList:
		return getPrefixFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// getPrefixFieldValue returns the value of a field for a Prefix List
func getPrefixFieldValue(fieldName string, prefix types.ManagedPrefixList) (string, error) {
	if getter, exists := prefixFieldValueGetters[fieldName]; exists {
		return getter(prefix)
	}
	return "", fmt.Errorf("field %s not found in prefixFieldValueGetters", fieldName)
}

// GetPrefixTagValue returns the value of a tag for the given instance.
func GetPrefixTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.ManagedPrefixList:
		for _, tag := range v.Tags {
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
// Prefix List field getters
// -----------------------------------------------------------------------------

func getPrefixListID(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.ManagedPrefixList).PrefixListId), nil
}

func getPrefixListName(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.ManagedPrefixList).PrefixListName), nil
}

func getPrefixListMaxEntries(instance any) (string, error) {
	prefix := instance.(types.ManagedPrefixList)
	return format.Int32ToStringOrDefault(prefix.MaxEntries, "-"), nil
}

func getPrefixListAddressFamily(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.ManagedPrefixList).AddressFamily), nil
}

func getPrefixListState(instance any) (string, error) {
	return format.Status(string(instance.(types.ManagedPrefixList).State)), nil
}

func getPrefixListVersion(instance any) (string, error) {
	prefix := instance.(types.ManagedPrefixList)
	return format.Int64ToStringOrDefault(prefix.Version, "-"), nil
}

func getPrefixListARN(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.ManagedPrefixList).PrefixListArn), nil
}

func getPrefixListOwner(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.ManagedPrefixList).OwnerId), nil
}
