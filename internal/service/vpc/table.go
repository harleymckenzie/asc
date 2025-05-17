package vpc

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

type VPCAttribute struct {
	GetValue func(*types.Vpc) string
}

func GetVPCAttributeValue(fieldID string, instance any) (string, error) {
	vpc, ok := instance.(types.Vpc)
	if !ok {
		return "", fmt.Errorf("instance is not a types.Vpc")
	}
	attr, exists := vpcAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&vpc), nil
}

func vpcAttributes() map[string]VPCAttribute {
	return map[string]VPCAttribute{
		"VPC ID": {
			GetValue: func(vpc *types.Vpc) string {
				return format.StringOrEmpty(vpc.VpcId)
			},
		},
		"State": {
			GetValue: func(vpc *types.Vpc) string {
				return string(vpc.State)
			},
		},
		"IPv4 CIDR": {
			GetValue: func(vpc *types.Vpc) string {
				return format.StringOrEmpty(vpc.CidrBlock)
			},
		},
		"IPv6 CIDR": {
			GetValue: func(vpc *types.Vpc) string {
				if len(vpc.Ipv6CidrBlockAssociationSet) == 0 {
					return ""
				}
				return format.StringOrEmpty(vpc.Ipv6CidrBlockAssociationSet[0].Ipv6CidrBlock)
			},
		},
		"DHCP Option Set": {
			GetValue: func(vpc *types.Vpc) string {
				return format.StringOrEmpty(vpc.DhcpOptionsId)
			},
		},
		"Main Route Table": {
			// TODO: Implement lookup for main route table
			GetValue: func(vpc *types.Vpc) string {
				return "-"
			},
		},
		"Main Network ACL": {
			// TODO: Implement lookup for main network ACL
			GetValue: func(vpc *types.Vpc) string {
				return "-"
			},
		},
		"Owner ID": {
			GetValue: func(vpc *types.Vpc) string {
				return format.StringOrEmpty(vpc.OwnerId)
			},
		},
		"Tenancy": {
			GetValue: func(vpc *types.Vpc) string {
				return string(vpc.InstanceTenancy)
			},
		},
		"Default VPC": {
			GetValue: func(vpc *types.Vpc) string {
				return format.BoolToLabel(vpc.IsDefault, "Yes", "No")
			},
		},
	}
}

// GetNetworkAclAttributeValue returns the value for a field from a NetworkAcl.
func GetNetworkAclAttributeValue(fieldID string, instance any) (string, error) {
	nacl, ok := instance.(types.NetworkAcl)
	if !ok {
		return "", fmt.Errorf("instance is not a types.NetworkAcl")
	}
	attr, exists := networkAclAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&nacl), nil
}

type NetworkAclAttribute struct {
	GetValue func(*types.NetworkAcl) string
}

func networkAclAttributes() map[string]NetworkAclAttribute {
	return map[string]NetworkAclAttribute{
		"Network ACL ID": {
			GetValue: func(n *types.NetworkAcl) string { return format.StringOrEmpty(n.NetworkAclId) },
		},
		"VPC ID": {
			GetValue: func(n *types.NetworkAcl) string { return format.StringOrEmpty(n.VpcId) },
		},
		"Is Default": {
			GetValue: func(n *types.NetworkAcl) string { return format.BoolToLabel(n.IsDefault, "Yes", "No") },
		},
		"Entry Count": {
			GetValue: func(n *types.NetworkAcl) string { return fmt.Sprintf("%d", len(n.Entries)) },
		},
		"Association Count": {
			GetValue: func(n *types.NetworkAcl) string { return fmt.Sprintf("%d", len(n.Associations)) },
		},
	}
}

// GetNatGatewayAttributeValue returns the value for a field from a NatGateway.
func GetNatGatewayAttributeValue(fieldID string, instance any) (string, error) {
	nat, ok := instance.(types.NatGateway)
	if !ok {
		return "", fmt.Errorf("instance is not a types.NatGateway")
	}
	attr, exists := natGatewayAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&nat), nil
}

type NatGatewayAttribute struct {
	GetValue func(*types.NatGateway) string
}

func natGatewayAttributes() map[string]NatGatewayAttribute {
	return map[string]NatGatewayAttribute{
		"NAT Gateway ID": {
			GetValue: func(n *types.NatGateway) string { return format.StringOrEmpty(n.NatGatewayId) },
		},
		"VPC ID": {
			GetValue: func(n *types.NatGateway) string { return format.StringOrEmpty(n.VpcId) },
		},
		"Subnet ID": {
			GetValue: func(n *types.NatGateway) string { return format.StringOrEmpty(n.SubnetId) },
		},
		"State": {
			GetValue: func(n *types.NatGateway) string { return string(n.State) },
		},
		"Type": {
			GetValue: func(n *types.NatGateway) string { return string(n.ConnectivityType) },
		},
		"IP Addresses": {
			GetValue: func(n *types.NatGateway) string {
				ips := make([]string, 0, len(n.NatGatewayAddresses))
				for _, addr := range n.NatGatewayAddresses {
					ips = append(ips, format.StringOrEmpty(addr.PublicIp))
				}
				if len(ips) == 0 {
					return ""
				}
				return strings.Join(ips, ", ")
			},
		},
	}
}

// GetPrefixListAttributeValue returns the value for a field from a PrefixList.
func GetPrefixListAttributeValue(fieldID string, instance any) (string, error) {
	pl, ok := instance.(types.PrefixList)
	if !ok {
		return "", fmt.Errorf("instance is not a types.PrefixList")
	}
	attr, exists := prefixListAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&pl), nil
}

type PrefixListAttribute struct {
	GetValue func(*types.PrefixList) string
}

func prefixListAttributes() map[string]PrefixListAttribute {
	return map[string]PrefixListAttribute{
		"Prefix List ID": {
			GetValue: func(p *types.PrefixList) string { return format.StringOrEmpty(p.PrefixListId) },
		},
		"Name": {
			GetValue: func(p *types.PrefixList) string { return format.StringOrEmpty(p.PrefixListName) },
		},
		"CIDRs": {
			GetValue: func(p *types.PrefixList) string {
				if len(p.Cidrs) == 0 {
					return ""
				}
				return strings.Join(p.Cidrs, ", ")
			},
		},
	}
}

// GetRouteTableAttributeValue returns the value for a field from a RouteTable.
func GetRouteTableAttributeValue(fieldID string, instance any) (string, error) {
	rt, ok := instance.(types.RouteTable)
	if !ok {
		return "", fmt.Errorf("instance is not a types.RouteTable")
	}
	attr, exists := routeTableAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&rt), nil
}

type RouteTableAttribute struct {
	GetValue func(*types.RouteTable) string
}

func routeTableAttributes() map[string]RouteTableAttribute {
	return map[string]RouteTableAttribute{
		"Route Table ID": {
			GetValue: func(r *types.RouteTable) string { return format.StringOrEmpty(r.RouteTableId) },
		},
		"VPC ID": {
			GetValue: func(r *types.RouteTable) string { return format.StringOrEmpty(r.VpcId) },
		},
		"Association Count": {
			GetValue: func(r *types.RouteTable) string { return fmt.Sprintf("%d", len(r.Associations)) },
		},
		"Route Count": {
			GetValue: func(r *types.RouteTable) string { return fmt.Sprintf("%d", len(r.Routes)) },
		},
	}
}

// GetSubnetAttributeValue returns the value for a field from a Subnet.
func GetSubnetAttributeValue(fieldID string, instance any) (string, error) {
	sub, ok := instance.(types.Subnet)
	if !ok {
		return "", fmt.Errorf("instance is not a types.Subnet")
	}
	attr, exists := subnetAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&sub), nil
}

type SubnetAttribute struct {
	GetValue func(*types.Subnet) string
}

func subnetAttributes() map[string]SubnetAttribute {
	return map[string]SubnetAttribute{
		"Subnet ID": {
			GetValue: func(s *types.Subnet) string { return format.StringOrEmpty(s.SubnetId) },
		},
		"VPC ID": {
			GetValue: func(s *types.Subnet) string { return format.StringOrEmpty(s.VpcId) },
		},
		"CIDR Block": {
			GetValue: func(s *types.Subnet) string { return format.StringOrEmpty(s.CidrBlock) },
		},
		"Availability Zone": {
			GetValue: func(s *types.Subnet) string { return format.StringOrEmpty(s.AvailabilityZone) },
		},
		"State": {
			GetValue: func(s *types.Subnet) string { return string(s.State) },
		},
		"Available IPs": {
			GetValue: func(s *types.Subnet) string { return fmt.Sprintf("%d", s.AvailableIpAddressCount) },
		},
		"Default For AZ": {
			GetValue: func(s *types.Subnet) string { return format.BoolToLabel(s.DefaultForAz, "Yes", "No") },
		},
	}
}

// GetIGWAttributeValue returns the value for a field from an InternetGateway.
func GetIGWAttributeValue(fieldID string, instance any) (string, error) {
	igw, ok := instance.(types.InternetGateway)
	if !ok {
		return "", fmt.Errorf("instance is not a types.InternetGateway")
	}
	attr, exists := igwAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&igw), nil
}

type IGWAttribute struct {
	GetValue func(*types.InternetGateway) string
}

func igwAttributes() map[string]IGWAttribute {
	return map[string]IGWAttribute{
		"Internet Gateway ID": {
			GetValue: func(i *types.InternetGateway) string { return format.StringOrEmpty(i.InternetGatewayId) },
		},
		"VPC Attachments": {
			GetValue: func(i *types.InternetGateway) string {
				if len(i.Attachments) == 0 {
					return ""
				}
				var vpcs []string
				for _, a := range i.Attachments {
					vpcs = append(vpcs, format.StringOrEmpty(a.VpcId))
				}
				return strings.Join(vpcs, ", ")
			},
		},
		"State": {
			GetValue: func(i *types.InternetGateway) string {
				if len(i.Attachments) == 0 {
					return "-"
				}
				return string(i.Attachments[0].State)
			},
		},
	}
}

// Helper: Find the main route table for a VPC from a list of route tables
func FindMainRouteTable(vpcID string, routeTables []types.RouteTable) string {
	for _, rt := range routeTables {
		if rt.VpcId != nil && *rt.VpcId == vpcID {
			for _, assoc := range rt.Associations {
				if assoc.Main != nil && *assoc.Main {
					return format.StringOrEmpty(rt.RouteTableId)
				}
			}
		}
	}
	return "-"
}

// Helper: Find the main network ACL for a VPC from a list of network ACLs
func FindMainNetworkACL(vpcID string, acls []types.NetworkAcl) string {
	for _, acl := range acls {
		if acl.VpcId != nil && *acl.VpcId == vpcID && acl.IsDefault != nil && *acl.IsDefault {
			return format.StringOrEmpty(acl.NetworkAclId)
		}
	}
	return "-"
}
