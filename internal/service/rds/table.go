package rds

import (
	"fmt"
	"strings"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// Attribute is a struct that defines a field in a detailed table.
type Attribute struct {
	GetValue func(*types.DBInstance, []types.DBCluster) string
}

// GetAttributeValue returns the value for a given field and DBInstance.
func GetAttributeValue(fieldID string, instance any, clusters []types.DBCluster) (string, error) {
	inst, ok := instance.(types.DBInstance)
	if !ok {
		return "", fmt.Errorf("instance is not a types.DBInstance")
	}
	attr, ok := availableAttributes()[fieldID]
	if !ok || attr.GetValue == nil {
		return "", fmt.Errorf("error getting attribute %q", fieldID)
	}
	return attr.GetValue(&inst, clusters), nil
}

func availableAttributes() map[string]Attribute {
	return map[string]Attribute{
		"Cluster Identifier": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				if i.DBClusterIdentifier != nil {
					return format.StringOrEmpty(i.DBClusterIdentifier)
				}
				return "-"
			},
		},
		"Identifier": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.DBInstanceIdentifier)
			},
		},
		"Status": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.Status(aws.ToString(i.DBInstanceStatus))
			},
		},
		"Engine": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return string(*i.Engine)
			},
		},
		"Engine Version": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return string(*i.EngineVersion)
			},
		},
		"Size": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return string(*i.DBInstanceClass)
			},
		},
		"Role": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return getDBInstanceRole(*i, clusters)
			},
		},
		"Endpoint": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.Endpoint.Address)
			},
		},
		"ARN": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.DBInstanceArn)
			},
		},
		"VPC ID": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.DBSubnetGroup.VpcId)
			},
		},
		"Subnet Group": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.DBSubnetGroup.DBSubnetGroupName)
			},
		},
		"Subnet IDs": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				subnets := make([]string, len(i.DBSubnetGroup.Subnets))
				for i, subnet := range i.DBSubnetGroup.Subnets {
					subnets[i] = aws.ToString(subnet.SubnetIdentifier)
				}
				return strings.Join(subnets, ", ")
			},
		},
		"Public Access": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.BoolToLabel(i.PubliclyAccessible, "Yes", "No")
			},
		},
		"Network Type": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.DBSubnetGroup.DBSubnetGroupDescription)
			},
		},
		"VPC Security Group(s)": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				securityGroups := make([]string, len(i.VpcSecurityGroups))
				for i, securityGroup := range i.VpcSecurityGroups {
					securityGroups[i] = aws.ToString(securityGroup.VpcSecurityGroupId)
				}
				return strings.Join(securityGroups, ", ")
			},
		},
		"Certificate Authority": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.CACertificateIdentifier)
			},
		},
		"Certificate Authority Date": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				if i.CertificateDetails != nil {
					return format.TimeToStringOrEmpty(i.CertificateDetails.ValidTill)
				}
				return "-"
			},
		},
		"Engine Lifecycle": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.EngineLifecycleSupport)
			},
		},
		"DB Name": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.DBName)
			},
		},
		"Option Groups": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				if len(i.OptionGroupMemberships) > 0 {
					optionGroups := make([]string, len(i.OptionGroupMemberships))
					for i, optionGroup := range i.OptionGroupMemberships {
						optionGroups[i] = aws.ToString(optionGroup.OptionGroupName)
					}
					return strings.Join(optionGroups, ", ")
				}
				return "-"
			},
		},
		"DB Instance Parameter Group": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				if len(i.DBParameterGroups) > 0 {
					parameterGroups := make([]string, len(i.DBParameterGroups))
					for i, parameterGroup := range i.DBParameterGroups {
						parameterGroups[i] = aws.ToString(parameterGroup.DBParameterGroupName)
					}
					return strings.Join(parameterGroups, ", ")
				}
				return "-"
			},
		},
		"Cluster ID": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.DBClusterIdentifier)
			},
		},
		"Created Time": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.TimeToStringOrEmpty(i.InstanceCreateTime)
			},
		},
		"Instance Type": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.DBInstanceClass)
			},
		},
		// Get the minimum and maximum capacity from the clusters
		"Minimum ACU Capacity": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				for _, cluster := range clusters {
					if aws.ToString(cluster.DBClusterIdentifier) == aws.ToString(i.DBClusterIdentifier) {
						if cluster.ServerlessV2ScalingConfiguration != nil {
							return format.Float64ToStringOrEmpty(cluster.ServerlessV2ScalingConfiguration.MinCapacity)
						} else if cluster.ScalingConfigurationInfo != nil {
							return format.Int32ToStringOrEmpty(cluster.ScalingConfigurationInfo.MinCapacity)
						}
					}
				}
				return ""
			},
		},
		"Maximum ACU Capacity": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				for _, cluster := range clusters {
					if aws.ToString(cluster.DBClusterIdentifier) == aws.ToString(i.DBClusterIdentifier) {
						if cluster.ServerlessV2ScalingConfiguration != nil {
							return format.Float64ToStringOrEmpty(cluster.ServerlessV2ScalingConfiguration.MaxCapacity)
						} else if cluster.ScalingConfigurationInfo != nil {
							return format.Int32ToStringOrEmpty(cluster.ScalingConfigurationInfo.MaxCapacity)
						}
					}
				}
				return ""
			},
		},
		"Failover Priority": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.Int32ToStringOrEmpty(i.PromotionTier)
			},
		},
		"Encryption": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				if i.StorageEncrypted != nil && *i.StorageEncrypted {
					return "Yes"
				}
				return "No"
			},
		},
		"AWS KMS Key": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.KmsKeyId)
			},
		},
		"Storage Type": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.StorageType)
			},
		},
		"Monitoring Type": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return "Database Insights - " + string(i.DatabaseInsightsMode)
			},
		},
		"Performance Insights": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				if i.PerformanceInsightsEnabled != nil && *i.PerformanceInsightsEnabled {
					return "Yes"
				}
				return "No"
			},
		},
		"Enhanced Monitoring": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.Int32ToStringOrEmpty(i.MonitoringInterval)
			},
		},
		"Granularity": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.Int32ToStringOrEmpty(i.MonitoringInterval)
			},
		},
		"Monitoring Role": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.MonitoringRoleArn)
			},
		},
		"Auto Minor Version Upgrade": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				if i.AutoMinorVersionUpgrade != nil && *i.AutoMinorVersionUpgrade {
					return "Yes"
				}
				return "No"
			},
		},
		"Maintenance Window": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.StringOrEmpty(i.PreferredMaintenanceWindow)
			},
		},
		"Pending Maintenance": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				if i.PendingModifiedValues != nil {
					return string(i.PendingModifiedValues.AutomationMode)
				}
				return "No"
			},
		},
		"Pending Modifications": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				if i.PendingModifiedValues != nil {
					return aws.ToString(i.PendingModifiedValues.DBInstanceClass)
				}
				return "No"
			},
		},
	}
}
