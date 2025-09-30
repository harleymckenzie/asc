// The show command displays detailed information about an RDS instance.

package rds

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
		{Name: "Endpoint", Category: "Connectivity & security", Visible: true},
		{Name: "Port", Category: "Connectivity & security", Visible: true},
		{Name: "Availability Zone", Category: "Connectivity & security", Visible: true},
		{Name: "VPC ID", Category: "Connectivity & security", Visible: true},
		{Name: "Subnet Group", Category: "Connectivity & security", Visible: true},
		{Name: "Subnets", Category: "Connectivity & security", Visible: true},
		{Name: "Security Group(s)", Category: "Connectivity & security", Visible: true},
		{Name: "Network Type", Category: "Connectivity & security", Visible: true},
		{Name: "Publicly Accessible", Category: "Connectivity & security", Visible: true},
		{Name: "Certificate Authority", Category: "Connectivity & security", Visible: true},
		{Name: "Certificate Expiry Date", Category: "Connectivity & security", Visible: true},

		{Name: "Identifier", Category: "Configuration", Visible: true},
		{Name: "Engine Version", Category: "Configuration", Visible: true},
		{Name: "RDS Extended Support", Category: "Configuration", Visible: true},
		{Name: "DB Name", Category: "Configuration", Visible: true},
		{Name: "Option Group", Category: "Configuration", Visible: true},
		{Name: "Parameter Group", Category: "Configuration", Visible: true},
		{Name: "ARN", Category: "Configuration", Visible: true},
		{Name: "Resource ID", Category: "Configuration", Visible: true},
		{Name: "Created Time", Category: "Configuration", Visible: true},

		{Name: "Class", Category: "Instance Class", Visible: true},

		{Name: "Failover Priority", Category: "Availability", Visible: true},

		{Name: "Encryption", Category: "Primary Storage", Visible: true},
		{Name: "AWS KMS Key", Category: "Primary Storage", Visible: true},
		{Name: "Storage Type", Category: "Primary Storage", Visible: true},

		{Name: "Performance Insights", Category: "Monitoring", Visible: true},
		{Name: "Monitoring Interval", Category: "Monitoring", Visible: true},
		{Name: "Monitoring Role", Category: "Monitoring", Visible: true},

		{Name: "Auto Minor Version Upgrade", Category: "Maintenance & Backups", Visible: true},
		{Name: "Maintenance Window", Category: "Maintenance & Backups", Visible: true},
		{Name: "Pending Modifications", Category: "Maintenance & Backups", Visible: true},
	}
}

// Command variable
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a RDS instance",
	Aliases: []string{"describe"},
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowRDSInstance(cobraCmd, args))
	},
}

// Flag function
func NewShowFlags(showCmd *cobra.Command) {
	cmdutil.AddShowFlags(showCmd, "vertical")
}

// Show function
func ShowRDSInstance(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, rds.NewRDSService)
	if err != nil {
		return fmt.Errorf("create new RDS service: %w", err)
	}

	instance, err := svc.GetInstances(cmd.Context(), &ascTypes.GetInstancesInput{
		InstanceIdentifier: args[0],
	})
	if err != nil {
		return fmt.Errorf("get instance: %w", err)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("Database Details\n(%s)", args[0]),
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(instance[0], getShowFields(), rds.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))
	tags, err := populateRDSTagFields(instance[0].TagList)
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
