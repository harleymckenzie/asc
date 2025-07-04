// The show command displays detailed information about an EC2 instance.

package rds

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/rds"
	ascTypes "github.com/harleymckenzie/asc/internal/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {
	addShowFlags(showCmd)
}

// Column functions
func rdsShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Instance Details", Header: true},
		{ID: "Identifier", Display: true},
		{ID: "ARN", Display: true},
		{ID: "Status", Display: true},
		{ID: "Role", Display: true},
		{ID: "Engine", Display: true},
		{ID: "Engine Version", Display: true},
		{ID: "VPC ID", Display: true},
		{ID: "Subnet Group", Display: true},
		{ID: "Subnet IDs", Display: true},
		{ID: "Network Type", Display: true},
		{ID: "VPC Security Group(s)", Display: true},
		{ID: "Public Access", Display: true},
		{ID: "Certificate Authority", Display: true},
		{ID: "Certificate Authority Date", Display: true},

		{ID: "Configuration", Header: true},
		{ID: "Cluster ID", Display: true},
		{ID: "Engine Lifecycle", Display: true},
		{ID: "DB Name", Display: true},
		{ID: "Option Groups", Display: true},
		{ID: "Created Time", Display: true},
		{ID: "DB Instance Parameter Group", Display: true},
		{ID: "Instance Type", Display: true},
		{ID: "Minimum ACU Capacity", Display: true},
		{ID: "Maximum ACU Capacity", Display: true},
		{ID: "Failover Priority", Display: true},

		{ID: "Storage", Header: true},
		{ID: "Encryption", Display: true},
		{ID: "AWS KMS Key", Display: true},
		{ID: "Storage Type", Display: true},

		{ID: "Monitoring", Header: true},
		{ID: "Monitoring Type", Display: true},
		{ID: "Performance Insights", Display: true},
		{ID: "Enhanced Monitoring", Display: true},
		{ID: "Granularity", Display: true},
		{ID: "Monitoring Role", Display: true},

		{ID: "Maintenance", Header: true},
		{ID: "Auto Minor Version Upgrade", Display: true},
		{ID: "Maintenance Window", Display: true},
		{ID: "Pending Maintenance", Display: true},
		{ID: "Pending Modifications", Display: true},
	}
}

// Command variable
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an RDS instance",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowRDSInstance(cobraCmd, args))
	},
}

// Flag function
func addShowFlags(cmd *cobra.Command) {}

// ShowRDSInstance is the function for showing RDS instances
func ShowRDSInstance(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := rds.NewRDSService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new RDS service: %w", err)
	}

	instance, err := svc.GetInstances(ctx, &ascTypes.GetInstancesInput{
		DBInstanceIdentifier: args[0],
	})
	if err != nil {
		return fmt.Errorf("get instances: %w", err)
	}

	clusters, err := svc.GetClusters(ctx)
	if err != nil {
		return fmt.Errorf("get clusters: %w", err)
	}

	fields := rdsShowFields()
	opts := tableformat.RenderOptions{
		Title: "Instance summary for " + *instance[0].DBInstanceIdentifier,
		Style: "rounded",
		Layout: tableformat.DetailTableLayout{
			Type:          "horizontal",
			ColumnsPerRow: 3,
		},
	}

	err = tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: instance[0],
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return rds.GetAttributeValue(fieldID, instance, clusters)
		},
	}, opts)
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}
	return nil
}
