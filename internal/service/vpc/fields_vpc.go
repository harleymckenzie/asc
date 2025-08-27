package vpc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// VPC field getters
var vpcFieldValueGetters = map[string]FieldValueGetter{
	"VPC ID":           getVPCID,
	"State":            getVPCState,
	"IPv4 CIDR":        getVPCIPv4CIDR,
	"IPv6 CIDR":        getVPCIPv6CIDR,
	"DHCP Option Set":  getVPCDHCPOptions,
	"Main Route Table": getVPCMainRouteTable,
	"Main Network ACL": getVPCMainNetworkACL,
	"Tenancy":          getVPCTenancy,
	"Default VPC":      getVPCIsDefault,
}

// GetVPCFieldValue returns the value of a field for the given VPC instance.
func GetVPCFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.Vpc:
		return getVPCFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// getVPCFieldValue returns the value of a field for a VPC
func getVPCFieldValue(fieldName string, vpc types.Vpc) (string, error) {
	if getter, exists := vpcFieldValueGetters[fieldName]; exists {
		return getter(vpc)
	}
	return "", fmt.Errorf("field %s not found in vpcFieldValueGetters", fieldName)
}

// -----------------------------------------------------------------------------
// VPC field getters
// -----------------------------------------------------------------------------

func getVPCID(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.Vpc).VpcId), nil
}

func getVPCState(instance any) (string, error) {
	return format.Status(string(instance.(types.Vpc).State)), nil
}

func getVPCIPv4CIDR(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.Vpc).CidrBlock), nil
}

func getVPCIPv6CIDR(instance any) (string, error) {
	vpc := instance.(types.Vpc)
	if len(vpc.Ipv6CidrBlockAssociationSet) == 0 {
		return "", nil
	}
	return format.StringOrEmpty(vpc.Ipv6CidrBlockAssociationSet[0].Ipv6CidrBlock), nil
}

func getVPCDHCPOptions(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.Vpc).DhcpOptionsId), nil
}

func getVPCMainRouteTable(instance any) (string, error) {
	// TODO: Implement lookup for main route table
	return "-", nil
}

func getVPCMainNetworkACL(instance any) (string, error) {
	// TODO: Implement lookup for main network ACL
	return "-", nil
}

func getVPCTenancy(instance any) (string, error) {
	vpc := instance.(types.Vpc)
	if vpc.InstanceTenancy == "" {
		return "", nil
	}
	return string(vpc.InstanceTenancy), nil
}

func getVPCIsDefault(instance any) (string, error) {
	vpc := instance.(types.Vpc)
	if vpc.IsDefault == nil {
		return "", nil
	}
	return format.BoolToLabel(vpc.IsDefault, "Yes", "No"), nil
}
