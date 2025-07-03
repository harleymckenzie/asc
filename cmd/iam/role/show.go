// show.go displays detailed information about an IAM role.
package role

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/harleymckenzie/asc/internal/service/iam"
	ascTypes "github.com/harleymckenzie/asc/internal/service/iam/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

// Compose a struct to hold both the role and its inline policies
type roleWithPolicies struct {
	Role           any
	InlinePolicies string
	ManagedPolicies string
}

// iamRoleShowFields returns the fields for the role detail table.
func iamRoleShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Name", Display: true},
		{ID: "Role ID", Display: true},
		{ID: "Last Activity", Display: true},
		{ID: "Arn", Display: true},
		{ID: "Creation Time", Display: true},
		{ID: "Description", Display: true},
		{ID: "Path", Display: true},
		{ID: "Max Session Duration", Display: true},
		{ID: "Permissions Boundary", Display: true},
		{ID: "Assume Role Policy Document", Display: true},
		{ID: "Inline Policies", Display: true},
		{ID: "Managed Policies", Display: true},
	}
}

// showCmd is the cobra command for showing AMI details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an IAM role",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowIAMRole(cmd, args[0]))
	},
}

// NewShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {}

// ShowIAMRole displays detailed information for a specified IAM role.
func ShowIAMRole(cmd *cobra.Command, arg string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := iam.NewIAMService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new IAM service: %w", err)
	}

	roles, err := svc.GetRoles(ctx, &ascTypes.GetRolesInput{RoleName: arg})
	if err != nil {
		return fmt.Errorf("get roles: %w", err)
	}
	if len(roles) == 0 {
		return fmt.Errorf("role not found: %s", arg)
	}

	rolePolicies, err := svc.GetRoleInlinePolicies(ctx, &ascTypes.GetRolePoliciesInput{RoleName: arg})
	if err != nil {
		return fmt.Errorf("get role policies: %w", err)
	}

	inlinePolicies := ""
	if len(rolePolicies) > 0 {
		for i, name := range rolePolicies {
			if i > 0 {
				inlinePolicies += "\n"
			}
			inlinePolicies += aws.ToString(name.PolicyName)
		}
	}

	managedPolicies, err := svc.GetRoleManagedPolicies(ctx, &ascTypes.GetRoleManagedPoliciesInput{RoleName: arg})
	if err != nil {
		return fmt.Errorf("get role managed policies: %w", err)
	}

	managedPoliciesStr := ""
	if len(managedPolicies) > 0 {
		for i, name := range managedPolicies {
			if i > 0 {
				managedPoliciesStr += "\n"
			}
			managedPoliciesStr += aws.ToString(name.PolicyName)
		}
	}

	fields := iamRoleShowFields()
	opts := tableformat.RenderOptions{
		Title: fmt.Sprintf("Role Details\n(%s)", aws.ToString(roles[0].RoleName)),
		Style: "rounded",
		Layout: tableformat.DetailTableLayout{
			Type: "vertical",
		},
		SortBy: tableformat.GetSortByField(fields, false),
	}

	instance := roleWithPolicies{
		Role:           roles[0],
		InlinePolicies: inlinePolicies,
		ManagedPolicies: managedPoliciesStr,
	}

	return tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: instance,
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			if fieldID == "Inline Policies" {
				if rwp, ok := instance.(roleWithPolicies); ok {
					return rwp.InlinePolicies, nil
				}
				return "", nil
			}
			if fieldID == "Managed Policies" {
				if rwp, ok := instance.(roleWithPolicies); ok {
					return rwp.ManagedPolicies, nil
				}
				return "", nil
			}
			if rwp, ok := instance.(roleWithPolicies); ok {
				return iam.GetRoleAttributeValue(fieldID, rwp.Role)
			}
			return iam.GetRoleAttributeValue(fieldID, instance)
		},
	}, opts)
}
