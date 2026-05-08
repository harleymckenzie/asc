package user

import (
	"fmt"
	"strings"

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
		{Name: "User ID", Category: "General", Visible: true},
		{Name: "ARN", Category: "General", Visible: true},
		{Name: "Path", Category: "General", Visible: true},
		{Name: "Created Date", Category: "General", Visible: true},

		{Name: "Password Last Used", Category: "Activity", Visible: true},
	}
}

var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an IAM user",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowUser(cmd, args[0]))
	},
}

func NewShowFlags(cmd *cobra.Command) {
	cmdutil.AddShowFlags(cmd, "vertical")
}

func ShowUser(cmd *cobra.Command, userName string) error {
	svc, err := cmdutil.CreateService(cmd, iam.NewIAMService)
	if err != nil {
		return fmt.Errorf("create iam service: %w", err)
	}

	u, err := svc.GetUser(cmd.Context(), userName)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("User Details\n(%s)", aws.ToString(u.UserName)),
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(*u, getShowFields(), iam.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))

	groups, err := svc.GetGroupsForUser(cmd.Context(), userName)
	if err != nil {
		return fmt.Errorf("get groups for user: %w", err)
	}
	if len(groups) > 0 {
		var groupNames []string
		for _, g := range groups {
			groupNames = append(groupNames, aws.ToString(g.GroupName))
		}
		table.AddSection(tablewriter.BuildSection("Groups", []tablewriter.Field{
			{Name: "Groups", Value: strings.Join(groupNames, ", ")},
		}, tablewriter.Horizontal))
	}

	attachedPolicies, err := svc.GetAttachedUserPolicies(cmd.Context(), userName)
	if err != nil {
		return fmt.Errorf("get attached user policies: %w", err)
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

	tags := populateIAMTagFields(u.Tags)
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
