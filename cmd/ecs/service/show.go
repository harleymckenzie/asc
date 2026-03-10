package service

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/harleymckenzie/asc/internal/service/ecs"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ecs/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

var showCluster string

func init() {
	NewShowFlags(showCmd)
}

func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "General", Visible: true},
		{Name: "ARN", Category: "General", Visible: true},
		{Name: "Status", Category: "General", Visible: true},
		{Name: "Cluster", Category: "General", Visible: true},
		{Name: "Launch Type", Category: "General", Visible: true},
		{Name: "Platform Version", Category: "General", Visible: true},
		{Name: "Created Date", Category: "General", Visible: true},

		{Name: "Task Definition", Category: "Task", Visible: true},
		{Name: "Desired Count", Category: "Task", Visible: true},
		{Name: "Running Count", Category: "Task", Visible: true},
		{Name: "Pending Count", Category: "Task", Visible: true},
		{Name: "Scheduling", Category: "Task", Visible: true},
		{Name: "Deployment Config", Category: "Task", Visible: true},

		{Name: "Network Mode", Category: "Network", Visible: true},
		{Name: "Subnets", Category: "Network", Visible: true},
		{Name: "Security Groups", Category: "Network", Visible: true},
		{Name: "Public IP", Category: "Network", Visible: true},

		{Name: "Load Balancers", Category: "Load Balancing", Visible: true},
		{Name: "Role ARN", Category: "IAM", Visible: true},
	}
}

var showCmd = &cobra.Command{
	Use:     "show [service-name]",
	Short:   "Show detailed information about an ECS service",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowService(cmd, args[0]))
	},
}

func NewShowFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&showCluster, "cluster", "c", "", "Cluster name or ARN (required).")
	cmd.MarkFlagRequired("cluster")
	cmdutil.AddShowFlags(cmd, "vertical")
}

func ShowService(cmd *cobra.Command, serviceName string) error {
	svc, err := cmdutil.CreateService(cmd, ecs.NewECSService)
	if err != nil {
		return fmt.Errorf("create ECS service: %w", err)
	}

	services, err := svc.DescribeServices(cmd.Context(), &ascTypes.DescribeServicesInput{
		Cluster:  showCluster,
		Services: []string{serviceName},
	})
	if err != nil {
		return fmt.Errorf("describe service: %w", err)
	}

	if len(services) == 0 {
		return fmt.Errorf("service %s not found", serviceName)
	}

	service := services[0]

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("ECS Service Details\n(%s)", aws.ToString(service.ServiceName)),
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(service, getShowFields(), ecs.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))

	tags, err := populateECSTagFields(service.Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}

// populateECSTagFields converts ECS tags to tablewriter fields
func populateECSTagFields(tags []types.Tag) ([]tablewriter.Field, error) {
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
