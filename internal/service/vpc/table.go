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

type NetworkAclAttribute struct {
	GetValue func(*types.NetworkAcl) string
}

type NACLRuleAttribute struct {
	GetValue func(*types.NetworkAclEntry) string
}

type RouteTableAttribute struct {
	GetValue func(*types.RouteTable) string
}

type SubnetAttribute struct {
	GetValue func(*types.Subnet) string
}

type NatGatewayAttribute struct {
	GetValue func(*types.NatGateway) string
}

type PrefixListAttribute struct {
	GetValue func(*types.ManagedPrefixList) string
}

type IGWAttribute struct {
	GetValue func(*types.InternetGateway) string
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
				return format.Status(string(vpc.State))
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

func networkAclAttributes() map[string]NetworkAclAttribute {
	return map[string]NetworkAclAttribute{
		"Network ACL ID": {
			GetValue: func(n *types.NetworkAcl) string { return format.StringOrEmpty(n.NetworkAclId) },
		},
		"Associated with": {
			GetValue: func(n *types.NetworkAcl) string { return fmt.Sprintf("%d Subnets", len(n.Associations)) },
		},
		"Default": {
			GetValue: func(n *types.NetworkAcl) string { return format.BoolToLabel(n.IsDefault, "Yes", "No") },
		},
		"VPC ID": {
			GetValue: func(n *types.NetworkAcl) string { return format.StringOrEmpty(n.VpcId) },
		},
		"Inbound Rules": {
			GetValue: func(n *types.NetworkAcl) string { return fmt.Sprintf("%d Inbound rules", countNACLRules(n, true)) },
		},
		"Outbound Rules": {
			GetValue: func(n *types.NetworkAcl) string { return fmt.Sprintf("%d Outbound rules", countNACLRules(n, false)) },
		},
		"Owner": {
			GetValue: func(n *types.NetworkAcl) string { return format.StringOrEmpty(n.OwnerId) },
		},
	}
}

// GetNACLRuleAttributeValue returns the value for a field from a NetworkAclRule.
func GetNACLRuleAttributeValue(fieldID string, instance any) (string, error) {
	rule, ok := instance.(types.NetworkAclEntry)
	if !ok {
		return "", fmt.Errorf("instance is not a types.NetworkAclEntry")
	}
	attr, exists := naclRuleAttributes()[fieldID]
	if !exists {
		return "", fmt.Errorf("attribute %q does not exist", fieldID)
	}
	if attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q: GetValue is nil", fieldID)
	}
	return attr.GetValue(&rule), nil
}

func naclRuleAttributes() map[string]NACLRuleAttribute {
	return map[string]NACLRuleAttribute{
		"Rule number": {
			GetValue: func(r *types.NetworkAclEntry) string {
				ruleNumber := format.Int32ToStringOrDefault(r.RuleNumber, "-")
				if ruleNumber == "32767" {
					return "*"
				}
				return ruleNumber
			},
		},
		"Type": {
			GetValue: func(r *types.NetworkAclEntry) string {
				if r.Protocol == nil {
					return ""
				}
				if *r.Protocol == "-1" {
					return "All traffic"
				}
				return strings.ToUpper(*r.Protocol)
			},
		},
		"Protocol": {
			GetValue: func(r *types.NetworkAclEntry) string {
				if r.Protocol == nil {
					return ""
				}
				if *r.Protocol == "-1" {
					return "All"
				}
				return strings.ToUpper(*r.Protocol)
			},
		},
		"Port range": {
			GetValue: func(r *types.NetworkAclEntry) string {
				if r.PortRange == nil {
					return "All"
				}
				return fmt.Sprintf("%d-%d", r.PortRange.From, r.PortRange.To)
			},
		},
		"Source": {
			GetValue: func(r *types.NetworkAclEntry) string { return format.StringOrEmpty(r.CidrBlock) },
		},
		"Destination": {
			GetValue: func(r *types.NetworkAclEntry) string { return format.StringOrEmpty(r.CidrBlock) },
		},
		"Allow/Deny": {
			GetValue: func(r *types.NetworkAclEntry) string {
				return format.Status(getNACLRuleAction(r))
			},
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

func natGatewayAttributes() map[string]NatGatewayAttribute {
	return map[string]NatGatewayAttribute{
		"NAT Gateway ID": {
			GetValue: func(n *types.NatGateway) string { return format.StringOrEmpty(n.NatGatewayId) },
		},
		"Connectivity": {
			GetValue: func(n *types.NatGateway) string {
				if n.ConnectivityType == "public" {
					return "Public"
				}
				return "Private"
			},
		},
		"VPC ID": {
			GetValue: func(n *types.NatGateway) string { return format.StringOrEmpty(n.VpcId) },
		},
		"Subnet ID": {
			GetValue: func(n *types.NatGateway) string { return format.StringOrEmpty(n.SubnetId) },
		},
		"State": {
			GetValue: func(n *types.NatGateway) string { return format.Status(string(n.State)) },
		},
		"Primary Public IP": {
			GetValue: func(n *types.NatGateway) string { return getNatGatewayPrimaryPublicIP(n) },
		},
		"Primary Private IP": {
			GetValue: func(n *types.NatGateway) string { return getNatGatewayPrimaryPrivateIP(n) },
		},
		"Created": {
			GetValue: func(n *types.NatGateway) string {
				return format.TimeToStringOrEmpty(n.CreateTime)
			},
		},
	}
}

// GetPrefixListAttributeValue returns the value for a field from a PrefixList.
func GetPrefixListAttributeValue(fieldID string, instance any) (string, error) {
	pl, ok := instance.(types.ManagedPrefixList)
	if !ok {
		return "", fmt.Errorf("instance is not a types.ManagedPrefixList")
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

func prefixListAttributes() map[string]PrefixListAttribute {
	return map[string]PrefixListAttribute{
		"Prefix List ID": {
			GetValue: func(p *types.ManagedPrefixList) string { return format.StringOrEmpty(p.PrefixListId) },
		},
		"Prefix List Name": {
			GetValue: func(p *types.ManagedPrefixList) string { return format.StringOrEmpty(p.PrefixListName) },
		},
		"Max Entries": {
			GetValue: func(p *types.ManagedPrefixList) string {
				return format.Int32ToStringOrDefault(p.MaxEntries, "-")
			},
		},
		"Address Family": {
			GetValue: func(p *types.ManagedPrefixList) string {
				return format.StringOrEmpty(p.AddressFamily)
			},
		},
		"State": {
			GetValue: func(p *types.ManagedPrefixList) string {
				return format.Status(string(p.State))
			},
		},
		"Version": {
			GetValue: func(p *types.ManagedPrefixList) string {
				return format.Int64ToStringOrDefault(p.Version, "-")
			},
		},
		"Prefix List ARN": {
			GetValue: func(p *types.ManagedPrefixList) string {
				return format.StringOrEmpty(p.PrefixListArn)
			},
		},
		"Owner": {
			GetValue: func(p *types.ManagedPrefixList) string {
				return format.StringOrEmpty(p.OwnerId)
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
			GetValue: func(s *types.Subnet) string { return format.Status(string(s.State)) },
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

func igwAttributes() map[string]IGWAttribute {
	return map[string]IGWAttribute{
		"Internet Gateway ID": {
			GetValue: func(i *types.InternetGateway) string { return format.StringOrEmpty(i.InternetGatewayId) },
		},
		"VPC ID": {
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
				if i.Attachments[0].State == "available" {
					return format.Status("Attached")
				}
				return format.Status(string(i.Attachments[0].State))
			},
		},
		"Owner": {
			GetValue: func(i *types.InternetGateway) string {
				return format.StringOrEmpty(i.OwnerId)
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

// Helper: Determine whether an ACL entry is an inbound or outbound rule
func isInboundRule(entry *types.NetworkAclEntry) bool {
	return entry.Egress == nil || !*entry.Egress
}

// Helper: Count the number of inbound or outbound rules in an ACL
func countNACLRules(acl *types.NetworkAcl, inbound bool) int {
	count := 0
	for _, entry := range acl.Entries {
		if isInboundRule(&entry) == inbound {
			count++
		}
	}
	return count
}

func getNACLRuleAction(entry *types.NetworkAclEntry) string {
	if entry.RuleAction == "allow" {
		return "Allow"
	}
	return "Deny"
}

func getNatGatewayPrimaryPublicIP(nat *types.NatGateway) string {
	if len(nat.NatGatewayAddresses) == 0 {
		return ""
	}
	return format.StringOrEmpty(nat.NatGatewayAddresses[0].PublicIp)
}

func getNatGatewayPrimaryPrivateIP(nat *types.NatGateway) string {
	if len(nat.NatGatewayAddresses) == 0 {
		return ""
	}
	return format.StringOrEmpty(nat.NatGatewayAddresses[0].PrivateIp)
}
