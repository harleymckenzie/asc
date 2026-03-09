package organizations

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/harleymckenzie/asc/internal/service/organizations"
	ascTypes "github.com/harleymckenzie/asc/internal/service/organizations/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	cmdutil.AddShowFlags(showCmd, "horizontal")
}

var showCmd = &cobra.Command{
	Use:     "show [id]",
	Short:   "Show details for the organization, an account, or an OU",
	Aliases: []string{"describe"},
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(showOrganizations(cmd, args))
	},
}

func showOrganizations(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, organizations.NewOrganizationsService)
	if err != nil {
		return fmt.Errorf("create new Organizations service: %w", err)
	}

	// No argument: show organization details
	if len(args) == 0 {
		return showOrganization(cmd, svc)
	}

	id := args[0]

	// Detect type by ID prefix
	if strings.HasPrefix(id, "ou-") {
		return showOU(cmd, svc, id)
	}
	return showAccount(cmd, svc, id)
}

func showOrganization(cmd *cobra.Command, svc *organizations.OrganizationsService) error {
	org, err := svc.GetOrganization(cmd.Context())
	if err != nil {
		return fmt.Errorf("describe organization: %w", err)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   "Organization Details",
		Columns: 3,
	})

	var orgFields []tablewriter.Field
	orgFields = append(orgFields,
		tablewriter.Field{Name: "ID", Category: "Organization", Value: aws.ToString(org.Id), Visible: true},
		tablewriter.Field{Name: "ARN", Category: "Organization", Value: aws.ToString(org.Arn), Visible: true},
		tablewriter.Field{Name: "Master Account ID", Category: "Organization", Value: aws.ToString(org.MasterAccountId), Visible: true},
		tablewriter.Field{Name: "Master Account Email", Category: "Organization", Value: aws.ToString(org.MasterAccountEmail), Visible: true},
		tablewriter.Field{Name: "Feature Set", Category: "Organization", Value: string(org.FeatureSet), Visible: true},
	)

	// Add available policy types
	for _, pt := range org.AvailablePolicyTypes {
		orgFields = append(orgFields, tablewriter.Field{
			Name:     string(pt.Type),
			Category: "Available Policy Types",
			Value:    string(pt.Status),
			Visible:  true,
		})
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}
	table.AddSections(tablewriter.BuildSections(orgFields, layout))
	table.Render()
	return nil
}

func showAccount(cmd *cobra.Command, svc *organizations.OrganizationsService, accountID string) error {
	account, err := svc.GetAccount(cmd.Context(), &ascTypes.GetAccountInput{AccountID: accountID})
	if err != nil {
		return fmt.Errorf("describe account: %w", err)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("Account Details\n(%s)", accountID),
		Columns: 3,
	})

	joinedTimestamp := ""
	if account.JoinedTimestamp != nil {
		joinedTimestamp = account.JoinedTimestamp.Format("2006-01-02 15:04:05")
	}

	accountFields := []tablewriter.Field{
		{Name: "ID", Category: "Account", Value: aws.ToString(account.Id), Visible: true},
		{Name: "Name", Category: "Account", Value: aws.ToString(account.Name), Visible: true},
		{Name: "Email", Category: "Account", Value: aws.ToString(account.Email), Visible: true},
		{Name: "ARN", Category: "Account", Value: aws.ToString(account.Arn), Visible: true},
		{Name: "Status", Category: "Account", Value: string(account.Status), Visible: true},
		{Name: "Joined Method", Category: "Account", Value: string(account.JoinedMethod), Visible: true},
		{Name: "Joined", Category: "Account", Value: joinedTimestamp, Visible: true},
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}
	table.AddSections(tablewriter.BuildSections(accountFields, layout))

	// Add tags
	tags, err := svc.GetTagsForResource(cmd.Context(), accountID)
	if err != nil {
		return fmt.Errorf("list tags: %w", err)
	}
	if len(tags) > 0 {
		table.AddSection(tablewriter.BuildSection("Tags", populateTagFields(tags), tablewriter.Horizontal))
	}

	table.Render()
	return nil
}

func showOU(cmd *cobra.Command, svc *organizations.OrganizationsService, ouID string) error {
	ou, err := svc.GetOU(cmd.Context(), &ascTypes.GetOUInput{OUID: ouID})
	if err != nil {
		return fmt.Errorf("describe OU: %w", err)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("Organizational Unit Details\n(%s)", ouID),
		Columns: 3,
	})

	ouFields := []tablewriter.Field{
		{Name: "ID", Category: "Organizational Unit", Value: aws.ToString(ou.Id), Visible: true},
		{Name: "Name", Category: "Organizational Unit", Value: aws.ToString(ou.Name), Visible: true},
		{Name: "ARN", Category: "Organizational Unit", Value: aws.ToString(ou.Arn), Visible: true},
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}
	table.AddSections(tablewriter.BuildSections(ouFields, layout))

	// Add child OUs
	childOUs, err := svc.GetOUsForParent(cmd.Context(), ouID)
	if err != nil {
		return fmt.Errorf("list child OUs: %w", err)
	}
	if len(childOUs) > 0 {
		var ouFieldList []tablewriter.Field
		for _, child := range childOUs {
			ouFieldList = append(ouFieldList, tablewriter.Field{
				Name:  aws.ToString(child.Id),
				Value: aws.ToString(child.Name),
			})
		}
		table.AddSection(tablewriter.BuildSection("Child OUs", ouFieldList, tablewriter.Horizontal))
	}

	// Add accounts
	accounts, err := svc.GetAccountsForParent(cmd.Context(), ouID)
	if err != nil {
		return fmt.Errorf("list accounts: %w", err)
	}
	if len(accounts) > 0 {
		var accountFieldList []tablewriter.Field
		for _, account := range accounts {
			name := aws.ToString(account.Name)
			if account.Status == types.AccountStatusSuspended {
				name += " (suspended)"
			}
			accountFieldList = append(accountFieldList, tablewriter.Field{
				Name:  aws.ToString(account.Id),
				Value: name,
			})
		}
		table.AddSection(tablewriter.BuildSection("Accounts", accountFieldList, tablewriter.Horizontal))
	}

	// Add tags
	tags, err := svc.GetTagsForResource(cmd.Context(), ouID)
	if err != nil {
		return fmt.Errorf("list tags: %w", err)
	}
	if len(tags) > 0 {
		table.AddSection(tablewriter.BuildSection("Tags", populateTagFields(tags), tablewriter.Horizontal))
	}

	table.Render()
	return nil
}

// populateTagFields converts Organizations tags to tablewriter fields.
func populateTagFields(tags []types.Tag) []tablewriter.Field {
	var fields []tablewriter.Field
	for _, tag := range tags {
		if tag.Key != nil && tag.Value != nil {
			fields = append(fields, tablewriter.Field{
				Name:  aws.ToString(tag.Key),
				Value: aws.ToString(tag.Value),
			})
		}
	}
	return fields
}
