package rds

import (
	"fmt"
	"strconv"
	"time"

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
		"Auto Minor Version Upgrade": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return strconv.FormatBool(*i.AutoMinorVersionUpgrade)
			},
		},
		"Availability Zone": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.AvailabilityZone)
			},
		},
		"ARN": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.DBInstanceArn)
			},
		},
		"AWS KMS Key": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.KmsKeyId)
			},
		},
		"Certificate Authority": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.CertificateDetails.CAIdentifier)
			},
		},
		"Certificate Expiry Date": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return i.CertificateDetails.ValidTill.Format(time.DateTime)
			},
		},
		"Class": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return string(*i.DBInstanceClass)
			},
		},
		"Cluster Identifier": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				if i.DBClusterIdentifier != nil {
					return aws.ToString(i.DBClusterIdentifier)
				}
				return "-"
			},
		},
		"Created Time": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return i.InstanceCreateTime.Format(time.DateTime)
			},
		},
		"DB Name": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.DBName)
			},
		},
		"Encryption": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return strconv.FormatBool(*i.StorageEncrypted)
			},
		},
		"Endpoint": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.Endpoint.Address)
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
		"Failover Priority": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return strconv.Itoa(int(*i.PromotionTier))
			},
		},
		"Identifier": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.DBInstanceIdentifier)
			},
		},
		"Maintenance Window": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.PreferredMaintenanceWindow)
			},
		},
		"Monitoring Interval": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return strconv.Itoa(int(*i.MonitoringInterval))
			},
		},
		"Monitoring Role": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.MonitoringRoleArn)
			},
		},
		"Network Type": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return string(*i.DBSubnetGroup.VpcId)
			},
		},
		"Option Group": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.OptionGroupMemberships[0].OptionGroupName)
			},
		},
		"Parameter Group": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.DBParameterGroups[0].DBParameterGroupName)
			},
		},
		"Pending Modifications": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return fmt.Sprintf("%v", i.PendingModifiedValues)
			},
		},
		"Performance Insights": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return strconv.FormatBool(*i.PerformanceInsightsEnabled)
			},
		},
		"Port": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return strconv.Itoa(int(*i.Endpoint.Port))
			},
		},
		"Publicly Accessible": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return strconv.FormatBool(*i.PubliclyAccessible)
			},
		},
		"RDS Extended Support": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.EngineLifecycleSupport)
			},
		},
		"Resource ID": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.DbiResourceId)
			},
		},
		"Role": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return calculateDBInstanceRole(*i, clusters)
			},
		},
		"Security Group(s)": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.DBSecurityGroups[0].DBSecurityGroupName)
			},
		},
		"Status": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return format.Status(aws.ToString(i.DBInstanceStatus))
			},
		},
		"Storage Type": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.StorageType)
			},
		},
		"Subnet Group": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.DBSubnetGroup.DBSubnetGroupName)
			},
		},
		"Subnets": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.DBSubnetGroup.Subnets[0].SubnetIdentifier)
			},
		},
		"VPC ID": {
			GetValue: func(i *types.DBInstance, clusters []types.DBCluster) string {
				return aws.ToString(i.DBSubnetGroup.VpcId)
			},
		},
	}
}
