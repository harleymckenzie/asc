package ec2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var securityGroupFieldValueGetters = map[string]FieldValueGetter{
	"Group ID":      getSecurityGroupID,
	"Group Name":    getSecurityGroupName,
	"Description":   getSecurityGroupDescription,
	"VPC ID":        getSecurityGroupVPCID,
	"Owner ID":      getSecurityGroupOwnerID,
	"Ingress Count": getSecurityGroupIngressCount,
	"Egress Count":  getSecurityGroupEgressCount,
	"Tag Count":     getSecurityGroupTagCount,
}

var securityGroupRuleFieldValueGetters = map[string]FieldValueGetter{
	"Rule ID":     getSecurityGroupRuleID,
	"IP Version":  getSecurityGroupRuleIPVersion,
	"Type":        getSecurityGroupRuleType,
	"Protocol":    getSecurityGroupRuleProtocol,
	"Port Range":  getSecurityGroupRulePortRange,
	"Source":      getSecurityGroupRuleSource,
	"Destination": getSecurityGroupRuleDestination,
	"Description": getSecurityGroupRuleDescription,
}

// getSecurityGroupFieldValue returns the value of a field for an EC2 security group
func getSecurityGroupFieldValue(fieldName string, group any) (string, error) {
	if getter, exists := securityGroupFieldValueGetters[fieldName]; exists {
		value, err := getter(group)
		if err != nil {
			return "", fmt.Errorf("failed to get field value for %s: %w", fieldName, err)
		}
		return value, nil
	}
	return "", fmt.Errorf("field %s not found in security group fieldValueGetters", fieldName)
}

func getSecurityGroupRuleFieldValue(fieldName string, rule any) (string, error) {
	if getter, exists := securityGroupRuleFieldValueGetters[fieldName]; exists {
		value, err := getter(rule)
		if err != nil {
			return "", fmt.Errorf("failed to get field value for %s: %w", fieldName, err)
		}
		return value, nil
	}
	return "", fmt.Errorf("field %s not found in security group rule fieldValueGetters", fieldName)
}

// Individual field value getters

// Security Group Field Value Getters
func getSecurityGroupID(group any) (string, error) {
	return aws.ToString(group.(types.SecurityGroup).GroupId), nil
}

func getSecurityGroupName(group any) (string, error) {
	return aws.ToString(group.(types.SecurityGroup).GroupName), nil
}

func getSecurityGroupDescription(group any) (string, error) {
	return aws.ToString(group.(types.SecurityGroup).Description), nil
}

func getSecurityGroupVPCID(group any) (string, error) {
	return aws.ToString(group.(types.SecurityGroup).VpcId), nil
}

func getSecurityGroupOwnerID(group any) (string, error) {
	return aws.ToString(group.(types.SecurityGroup).OwnerId), nil
}

func getSecurityGroupIngressCount(group any) (string, error) {
	return strconv.Itoa(len(group.(types.SecurityGroup).IpPermissions)), nil
}

func getSecurityGroupEgressCount(group any) (string, error) {
	return strconv.Itoa(len(group.(types.SecurityGroup).IpPermissionsEgress)), nil
}

func getSecurityGroupTagCount(group any) (string, error) {
	return strconv.Itoa(len(group.(types.SecurityGroup).Tags)), nil
}

// Security Group Rule Field Value Getters
func getSecurityGroupRuleID(rule any) (string, error) {
	return aws.ToString(rule.(types.SecurityGroupRule).SecurityGroupRuleId), nil
}

func getSecurityGroupRuleIPVersion(rule any) (string, error) {
	if rule.(types.SecurityGroupRule).CidrIpv4 != nil && *rule.(types.SecurityGroupRule).CidrIpv4 != "" {
		return "IPv4", nil
	}
	if rule.(types.SecurityGroupRule).CidrIpv6 != nil && *rule.(types.SecurityGroupRule).CidrIpv6 != "" {
		return "IPv6", nil
	}
	return "-", nil
}

func getSecurityGroupRuleType(rule any) (string, error) {
	if rule.(types.SecurityGroupRule).FromPort != nil && rule.(types.SecurityGroupRule).ToPort != nil {
		if *rule.(types.SecurityGroupRule).FromPort == *rule.(types.SecurityGroupRule).ToPort {
			switch *rule.(types.SecurityGroupRule).FromPort {
			case -1:
				return "All traffic", nil
			case 22:
				return "SSH", nil
			case 25:
				return "SMTP", nil
			case 53:
				return "DNS", nil
			case 80:
				return "HTTP", nil
			case 110:
				return "POP3", nil
			case 143:
				return "IMAP", nil
			case 389:
				return "LDAP", nil
			case 443:
				return "HTTPS", nil
			case 445:
				return "SMB", nil
			case 465:
				return "SMTPS", nil
			case 993:
				return "IMAPS", nil
			case 995:
				return "POP3S", nil
			case 1433:
				return "MSSQL", nil
			case 2049:
				return "NFS", nil
			case 3306:
				return "MySQL/Aurora", nil
			case 3389:
				return "RDP", nil
			case 5439:
				return "Redshift", nil
			case 5432:
				return "PostgreSQL", nil
			case 1521:
				return "Oracle RDS", nil
			case 5985:
				return "WinRM-HTTP", nil
			case 5986:
				return "WinRM-HTTPS", nil
			case 20049:
				return "Elastic Graphics", nil
			case 9042:
				return "CQLSH / Cassandra", nil
			default:
				// Handle custom protocol types based on IpProtocol
				if rule.(types.SecurityGroupRule).IpProtocol != nil {
					return fmt.Sprintf("Custom (%s)", strings.ToUpper(*rule.(types.SecurityGroupRule).IpProtocol)), nil
				} else {
					return "Custom Protocol", nil
				}
			}
		}
	}
	return "", nil
}

func getSecurityGroupRuleProtocol(rule any) (string, error) {
	if rule.(types.SecurityGroupRule).IpProtocol == nil {
		return "", nil
	}
	if *rule.(types.SecurityGroupRule).IpProtocol == "-1" {
		return "All", nil
	}
	return strings.ToUpper(*rule.(types.SecurityGroupRule).IpProtocol), nil
}

func getSecurityGroupRulePortRange(rule any) (string, error) {
	if *rule.(types.SecurityGroupRule).FromPort == -1 {
		return "All", nil
	}
	if rule.(types.SecurityGroupRule).FromPort != nil && rule.(types.SecurityGroupRule).ToPort != nil {
		if *rule.(types.SecurityGroupRule).FromPort == *rule.(types.SecurityGroupRule).ToPort {
			return strconv.Itoa(int(*rule.(types.SecurityGroupRule).FromPort)), nil
		}
		return fmt.Sprintf("%d-%d", *rule.(types.SecurityGroupRule).FromPort, *rule.(types.SecurityGroupRule).ToPort), nil
	}
	return "", nil
}

func getSecurityGroupRuleSource(rule any) (string, error) {
	if rule.(types.SecurityGroupRule).IsEgress != nil && *rule.(types.SecurityGroupRule).IsEgress {
		return aws.ToString(rule.(types.SecurityGroupRule).CidrIpv4), nil
	}
	if rule.(types.SecurityGroupRule).CidrIpv4 != nil && *rule.(types.SecurityGroupRule).CidrIpv4 != "" {
		return aws.ToString(rule.(types.SecurityGroupRule).CidrIpv4), nil
	}
	if rule.(types.SecurityGroupRule).CidrIpv6 != nil && *rule.(types.SecurityGroupRule).CidrIpv6 != "" {
		return aws.ToString(rule.(types.SecurityGroupRule).CidrIpv6), nil
	}
	if rule.(types.SecurityGroupRule).ReferencedGroupInfo != nil && rule.(types.SecurityGroupRule).ReferencedGroupInfo.GroupId != nil {
		return aws.ToString(rule.(types.SecurityGroupRule).ReferencedGroupInfo.GroupId), nil
	}
	return "", nil
}

func getSecurityGroupRuleDestination(rule any) (string, error) {
	if rule.(types.SecurityGroupRule).IsEgress != nil && *rule.(types.SecurityGroupRule).IsEgress {
		return aws.ToString(rule.(types.SecurityGroupRule).CidrIpv4), nil
	}
	if rule.(types.SecurityGroupRule).CidrIpv6 != nil && *rule.(types.SecurityGroupRule).CidrIpv6 != "" {
		return aws.ToString(rule.(types.SecurityGroupRule).CidrIpv6), nil
	}
	if rule.(types.SecurityGroupRule).ReferencedGroupInfo != nil && rule.(types.SecurityGroupRule).ReferencedGroupInfo.GroupId != nil {
		return aws.ToString(rule.(types.SecurityGroupRule).ReferencedGroupInfo.GroupId), nil
	}
	return "", nil
}

func getSecurityGroupRuleDescription(rule any) (string, error) {
	return aws.ToString(rule.(types.SecurityGroupRule).Description), nil
}
