package cluster

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

func init() {
	NewShowFlags(showCmd)
}

func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "General", Visible: true},
		{Name: "ARN", Category: "General", Visible: true},
		{Name: "Status", Category: "General", Visible: true},

		{Name: "Active Services", Category: "Statistics", Visible: true},
		{Name: "Running Tasks", Category: "Statistics", Visible: true},
		{Name: "Pending Tasks", Category: "Statistics", Visible: true},
		{Name: "Registered Instances", Category: "Statistics", Visible: true},

		{Name: "Capacity Providers", Category: "Capacity", Visible: true},
		{Name: "Default Strategy", Category: "Capacity", Visible: true},
	}
}

var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an ECS cluster",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowCluster(cmd, args[0]))
	},
}

func NewShowFlags(cmd *cobra.Command) {
	cmdutil.AddShowFlags(cmd, "vertical")
}

func ShowCluster(cmd *cobra.Command, clusterName string) error {
	svc, err := cmdutil.CreateService(cmd, ecs.NewECSService)
	if err != nil {
		return fmt.Errorf("create ECS service: %w", err)
	}

	clusters, err := svc.DescribeClusters(cmd.Context(), &ascTypes.DescribeClustersInput{
		ClusterARNs: []string{clusterName},
	})
	if err != nil {
		return fmt.Errorf("describe cluster: %w", err)
	}

	if len(clusters) == 0 {
		return fmt.Errorf("cluster %s not found", clusterName)
	}

	cluster := clusters[0]

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("ECS Cluster Details\n(%s)", aws.ToString(cluster.ClusterName)),
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(cluster, getShowFields(), ecs.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))

	tags, err := populateECSTagFields(cluster.Tags)
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
