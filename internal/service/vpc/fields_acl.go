package vpc

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type NACLFieldValueGetter func(instance any) (string, error)
type NACLRuleFieldValueGetter func(instance any) (string, error)

// Network ACL field getters
var naclFieldValueGetters = map[string]NACLFieldValueGetter{
	"Network ACL ID":  getNACLID,
	"Associated with": getNACLAssociations,
	"Default":         getNACLIsDefault,
	"VPC ID":          getNACLVPCID,
	"Inbound Rules":   getNACLInboundRulesCount,
	"Outbound Rules":  getNACLOutboundRulesCount,
	"Owner":           getNACLOwner,
}

// Network ACL Rule field getters
var naclRuleFieldValueGetters = map[string]NACLRuleFieldValueGetter{
	"Rule number": getNACLRuleNumber,
	"Type":        getNACLRuleType,
	"Protocol":    getNACLRuleProtocol,
	"Port Range":  getNACLRulePortRange,
	"Source":      getNACLRuleSource,
	"Allow/Deny":  getNACLRuleAction,
}

// getNACLFieldValue returns the value of a field for a Network ACL
func getNACLFieldValue(fieldName string, nacl types.NetworkAcl) (string, error) {
	if getter, exists := naclFieldValueGetters[fieldName]; exists {
		value, err := getter(nacl)
		if err != nil {
			return "", fmt.Errorf("failed to get field value for %s: %w", fieldName, err)
		}
		return value, nil
	}
	return "", fmt.Errorf("field %s not found in naclFieldValueGetters", fieldName)
}

// getNACLRuleFieldValue returns the value of a field for a Network ACL Rule
func getNACLRuleFieldValue(fieldName string, rule types.NetworkAclEntry) (string, error) {
	if getter, exists := naclRuleFieldValueGetters[fieldName]; exists {
		return getter(rule)
	}
	return "", fmt.Errorf("field %s not found in naclRuleFieldValueGetters", fieldName)
}

// -----------------------------------------------------------------------------
// Network ACL field getters
// -----------------------------------------------------------------------------

func getNACLID(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.NetworkAcl).NetworkAclId), nil
}

func getNACLAssociations(instance any) (string, error) {
	nacl := instance.(types.NetworkAcl)
	return fmt.Sprintf("%d Subnets", len(nacl.Associations)), nil
}

func getNACLIsDefault(instance any) (string, error) {
	nacl := instance.(types.NetworkAcl)
	return format.BoolToLabel(nacl.IsDefault, "Yes", "No"), nil
}

func getNACLVPCID(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.NetworkAcl).VpcId), nil
}

func getNACLInboundRulesCount(instance any) (string, error) {
	nacl := instance.(types.NetworkAcl)
	// Use the function from table.go
	count := 0
	for _, entry := range nacl.Entries {
		if entry.Egress == nil || !*entry.Egress {
			count++
		}
	}
	return fmt.Sprintf("%d Inbound rules", count), nil
}

func getNACLOutboundRulesCount(instance any) (string, error) {
	nacl := instance.(types.NetworkAcl)
	// Use the function from table.go
	count := 0
	for _, entry := range nacl.Entries {
		if entry.Egress != nil && *entry.Egress {
			count++
		}
	}
	return fmt.Sprintf("%d Outbound rules", count), nil
}

func getNACLOwner(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.NetworkAcl).OwnerId), nil
}

// -----------------------------------------------------------------------------
// Network ACL Rule field getters
// -----------------------------------------------------------------------------

func getNACLRuleNumber(instance any) (string, error) {
	rule := instance.(types.NetworkAclEntry)
	ruleNumber := format.Int32ToStringOrDefault(rule.RuleNumber, "-")
	if ruleNumber == "32767" {
		return "*", nil
	}
	return ruleNumber, nil
}

func getNACLRuleType(instance any) (string, error) {
	rule := instance.(types.NetworkAclEntry)
	if rule.Protocol == nil {
		return "", nil
	}
	if *rule.Protocol == "-1" {
		return "All traffic", nil
	}
	return strings.ToUpper(*rule.Protocol), nil
}

func getNACLRuleProtocol(instance any) (string, error) {
	rule := instance.(types.NetworkAclEntry)
	if rule.Protocol == nil {
		return "", nil
	}
	return *rule.Protocol, nil
}

func getNACLRulePortRange(instance any) (string, error) {
	rule := instance.(types.NetworkAclEntry)
	if rule.PortRange == nil {
		return "All", nil
	}
	if rule.PortRange.From != nil && rule.PortRange.To != nil {
		if *rule.PortRange.From == *rule.PortRange.To {
			return fmt.Sprintf("%d", *rule.PortRange.From), nil
		}
		return fmt.Sprintf("%d-%d", *rule.PortRange.From, *rule.PortRange.To), nil
	}
	return "All", nil
}

func getNACLRuleSource(instance any) (string, error) {
	rule := instance.(types.NetworkAclEntry)
	if rule.CidrBlock != nil {
		return *rule.CidrBlock, nil
	}
	if rule.Ipv6CidrBlock != nil {
		return *rule.Ipv6CidrBlock, nil
	}
	return "", nil
}

func getNACLRuleAction(instance any) (string, error) {
	rule := instance.(types.NetworkAclEntry)
	if rule.RuleAction == types.RuleActionAllow {
		return "ALLOW", nil
	}
	return "DENY", nil
}

// -----------------------------------------------------------------------------
// Helper functions
// -----------------------------------------------------------------------------

// isInboundRule checks if a Network ACL entry is an inbound rule
func isInboundRule(entry *types.NetworkAclEntry) bool {
	return entry.Egress == nil || !*entry.Egress
}

// countNACLRules counts the number of inbound or outbound rules in an ACL
func countNACLRules(acl *types.NetworkAcl, inbound bool) int {
	count := 0
	for _, entry := range acl.Entries {
		if isInboundRule(&entry) == inbound {
			count++
		}
	}
	return count
}
