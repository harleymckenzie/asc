// The show command displays detailed information about an RDS cluster.

package cluster

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/harleymckenzie/asc/internal/service/rds"
	ascTypes "github.com/harleymckenzie/asc/internal/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {
	NewShowFlags(showCmd)
}

// Column functions
func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		// Connectivity & Security
		{Name: "Endpoint", Category: "Connectivity & Security", Visible: true},
		{Name: "Reader Endpoint", Category: "Connectivity & Security", Visible: true},
		{Name: "Custom Endpoints", Category: "Connectivity & Security", Visible: true},
		{Name: "Port", Category: "Connectivity & Security", Visible: true},
		{Name: "Availability Zones", Category: "Connectivity & Security", Visible: true},
		{Name: "Subnet Group", Category: "Connectivity & Security", Visible: true},
		{Name: "Security Groups", Category: "Connectivity & Security", Visible: true},
		{Name: "Publicly Accessible", Category: "Connectivity & Security", Visible: true},
		{Name: "Certificate Authority", Category: "Connectivity & Security", Visible: false},
		{Name: "Certificate Expiry Date", Category: "Connectivity & Security", Visible: false},

		// Configuration
		{Name: "Cluster Identifier", Category: "Configuration", Visible: true},
		{Name: "Engine Version", Category: "Configuration", Visible: true},
		{Name: "Resource ID", Category: "Configuration", Visible: true},
		{Name: "DB Cluster ARN", Category: "Configuration", Visible: true},
		{Name: "Network Type", Category: "Configuration", Visible: true},
		{Name: "DB Cluster Instance Class", Category: "Configuration", Visible: false},
		{Name: "Parameter Group", Category: "Configuration", Visible: true},
		{Name: "Deletion Protection", Category: "Configuration", Visible: true},

		// Authentication
		{Name: "IAM Database Authentication", Category: "Authentication", Visible: true},
		{Name: "Master Username", Category: "Authentication", Visible: true},

		// Availability
		{Name: "Multi AZ", Category: "Availability", Visible: true},

		// Encryption
		{Name: "Encryption", Category: "Encryption", Visible: true},
		{Name: "KMS Key ID", Category: "Encryption", Visible: true},

		// Monitoring
		{Name: "Performance Insights", Category: "Monitoring", Visible: true},
		{Name: "Monitoring Interval", Category: "Monitoring", Visible: true},
		{Name: "Monitoring Role", Category: "Monitoring", Visible: true},

		// Maintenance & Backups
		{Name: "Auto Minor Version Upgrade", Category: "Maintenance & Backups", Visible: true},
		{Name: "Preferred Maintenance Window", Category: "Maintenance & Backups", Visible: true},
		{Name: "Pending Modifications", Category: "Maintenance & Backups", Visible: true},
		{Name: "Backup Retention Period", Category: "Maintenance & Backups", Visible: true},
		{Name: "Preferred Backup Window", Category: "Maintenance & Backups", Visible: true},
		{Name: "Copy Tags To Snapshot", Category: "Maintenance & Backups", Visible: true},
		{Name: "Earliest Restorable Time", Category: "Maintenance & Backups", Visible: true},
		{Name: "Latest Restorable Time", Category: "Maintenance & Backups", Visible: true},
	}
}

// Command variable
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a RDS cluster",
	Aliases: []string{"describe"},
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowRDSCluster(cobraCmd, args))
	},
}

// Flag function
func NewShowFlags(showCmd *cobra.Command) {
	cmdutil.AddShowFlags(showCmd, "vertical")
}

// Show function
func ShowRDSCluster(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, rds.NewRDSService)
	if err != nil {
		return fmt.Errorf("create new RDS service: %w", err)
	}

	cluster, err := svc.GetClusters(cmd.Context(), &ascTypes.GetClustersInput{
		ClusterIdentifier: args[0],
	})
	if err != nil {
		return fmt.Errorf("get cluster: %w", err)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("RDS Cluster Details\n(%s)", args[0]),
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(cluster[0], getShowFields(), rds.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))
	tags, err := populateRDSTagFields(cluster[0].TagList)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}

// populateRDSTagFields converts RDS tags to tablewriter fields
func populateRDSTagFields(tags []types.Tag) ([]tablewriter.Field, error) {
	var fields []tablewriter.Field
	for _, tag := range tags {
		if tag.Key != nil && tag.Value != nil {
			fields = append(fields, tablewriter.Field{
				Category: "Tag",
				Name:     aws.ToString(tag.Key),
				Value:    aws.ToString(tag.Value),
			})
		}
	}
	return fields, nil
}
