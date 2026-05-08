package policy

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
		{Name: "Policy ID", Category: "General", Visible: true},
		{Name: "ARN", Category: "General", Visible: true},
		{Name: "Path", Category: "General", Visible: true},
		{Name: "Created Date", Category: "General", Visible: true},
		{Name: "Updated Date", Category: "General", Visible: true},
		{Name: "Description", Category: "General", Visible: true},

		{Name: "Default Version ID", Category: "Versioning", Visible: true},

		{Name: "Attachment Count", Category: "Usage", Visible: true},
		{Name: "Permissions Boundary Count", Category: "Usage", Visible: true},
	}
}

var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an IAM policy",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	Example: "  asc iam policy show MyPolicyName\n" +
		"  asc iam policy show arn:aws:iam::123456789012:policy/MyPolicy",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowPolicy(cmd, args[0]))
	},
}

func NewShowFlags(cmd *cobra.Command) {
	cmdutil.AddShowFlags(cmd, "vertical")
}

func ShowPolicy(cmd *cobra.Command, identifier string) error {
	svc, err := cmdutil.CreateService(cmd, iam.NewIAMService)
	if err != nil {
		return fmt.Errorf("create iam service: %w", err)
	}

	p, err := svc.GetPolicy(cmd.Context(), identifier)
	if err != nil {
		return fmt.Errorf("get policy: %w", err)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("Policy Details\n(%s)", aws.ToString(p.PolicyName)),
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(*p, getShowFields(), iam.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))

	tags := populateIAMTagFields(p.Tags)
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
