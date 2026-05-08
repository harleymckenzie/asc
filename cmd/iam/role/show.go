package role

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/harleymckenzie/asc/internal/service/iam"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	NewShowFlags(showCmd)
}

func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Role ID", Category: "General", Visible: true},
		{Name: "ARN", Category: "General", Visible: true},
		{Name: "Path", Category: "General", Visible: true},
		{Name: "Created Date", Category: "General", Visible: true},
		{Name: "Description", Category: "General", Visible: true},

		{Name: "Max Session Duration", Category: "Configuration", Visible: true},

		{Name: "Last Used Date", Category: "Activity", Visible: true},
		{Name: "Last Used Region", Category: "Activity", Visible: true},
	}
}

var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an IAM role",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowRole(cmd, args[0]))
	},
}

func NewShowFlags(cmd *cobra.Command) {
	cmdutil.AddShowFlags(cmd, "vertical")
}

func ShowRole(cmd *cobra.Command, roleName string) error {
	svc, err := cmdutil.CreateService(cmd, iam.NewIAMService)
	if err != nil {
		return fmt.Errorf("create iam service: %w", err)
	}

	role, err := svc.GetRole(cmd.Context(), roleName)
	if err != nil {
		return fmt.Errorf("get role: %w", err)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("Role Details\n(%s)", aws.ToString(role.RoleName)),
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(*role, getShowFields(), iam.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))

	attachedPolicies, err := svc.GetAttachedRolePolicies(cmd.Context(), roleName)
	if err != nil {
		return fmt.Errorf("get attached role policies: %w", err)
	}
	if len(attachedPolicies) > 0 {
		var policyFields []tablewriter.Field
		for _, p := range attachedPolicies {
			policyFields = append(policyFields, tablewriter.Field{
				Name:  aws.ToString(p.PolicyName),
				Value: aws.ToString(p.PolicyArn),
			})
		}
		table.AddSection(tablewriter.BuildSection("Attached Policies", policyFields, tablewriter.Horizontal))
	}

	tags := populateIAMTagFields(role.Tags)
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}

func populateIAMTagFields(tags []iamTypes.Tag) []tablewriter.Field {
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
	return fields
}
