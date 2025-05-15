// ls.go defines the 'ls' subcommand for security group operations.
package security_group

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	sortID      bool
	sortName    bool
	sortVPCID   bool
	sortOwnerID bool
	showDesc    bool
	showOwnerID bool
	reverseSort bool
)

func init() {
	NewLsFlags(lsCmd)
}

// ec2SecurityGroupListFields returns the fields for the security group list table.
func ec2SecurityGroupListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Group Name", Visible: true, Sort: sortName, DefaultSort: true},
		{ID: "Group ID", Visible: true, Sort: sortID},
		{ID: "Description", Visible: showDesc},
		{ID: "VPC ID", Visible: true, Sort: sortVPCID},
		{ID: "Owner ID", Visible: showOwnerID, Sort: sortOwnerID},
		{ID: "Ingress Count", Visible: true},
		{ID: "Egress Count", Visible: true},
		{ID: "Tag Count", Visible: false},
	}
}

func ec2SecurityGroupRulesFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Rule ID", Visible: true},
		{ID: "IP Version", Visible: true},
		{ID: "Type", Visible: true},
		{ID: "Protocol", Visible: true},
		{ID: "Port Range", Visible: true},
		{ID: "Source", Visible: true},      // Inbound rules only
		{ID: "Destination", Visible: true}, // Outbound rules only
		{ID: "Description", Visible: true},
	}
}

// lsCmd is the cobra command for listing security groups.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all security groups",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListSecurityGroups(cmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs security groups in list format.")
	cobraCmd.Flags().BoolVarP(&sortID, "sort-id", "i", false, "Sort by descending group ID.")
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending group name.")
	cobraCmd.Flags().
		BoolVarP(&showDesc, "show-description", "d", false, "Show the security group description column.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
	cobraCmd.Flags().BoolVarP(&sortVPCID, "sort-vpc-id", "v", false, "Sort by descending VPC ID.")
	cobraCmd.Flags().
		BoolVarP(&sortOwnerID, "sort-owner-id", "o", false, "Sort by descending owner ID.")
	cobraCmd.Flags().
		BoolVarP(&showOwnerID, "show-owner-id", "O", false, "Show the security group owner ID column.")
}

// ListSecurityGroups is the handler for the ls subcommand.
// If a security group name is provided, it will list the IP permissions for that security group.
// Otherwise, it will list all security groups.
func ListSecurityGroups(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	if len(args) > 0 {
		return ListSecurityGroupRules(cmd, args)
	} else {
		groups, err := svc.GetSecurityGroups(ctx, &ascTypes.GetSecurityGroupsInput{})
		if err != nil {
			return fmt.Errorf("get security groups: %w", err)
		}

		fields := ec2SecurityGroupListFields()
		opts := tableformat.RenderOptions{
			Title:  "Security Groups",
			Style:  "rounded",
			SortBy: tableformat.GetSortByField(fields, reverseSort),
		}

		tableformat.RenderTableList(&tableformat.ListTable{
			Instances: utils.SlicesToAny(groups),
			Fields:    fields,
			GetAttribute: func(fieldID string, instance any) (string, error) {
				return ec2.GetSecurityGroupAttributeValue(fieldID, instance)
			},
		}, opts)
		return nil
	}
}

// ListSecurityGroupRules lists the rules for the provided security group.
func ListSecurityGroupRules(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	rules, err := svc.GetSecurityGroupRules(ctx, &ascTypes.GetSecurityGroupRulesInput{
		SecurityGroupID: args[0],
	})
	if err != nil {
		return fmt.Errorf("get security groups: %w", err)
	}

	ingressRules := ec2.FilterSecurityGroupRules(rules, false)
	egressRules := ec2.FilterSecurityGroupRules(rules, true)

	fields := ec2SecurityGroupRulesFields()
	ingressOpts := tableformat.RenderOptions{
		Title:  fmt.Sprintf("%s - Inbound Rules", args[0]),
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}
	egressOpts := tableformat.RenderOptions{
		Title:  fmt.Sprintf("%s - Outbound Rules", args[0]),
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	// Print inbound and outbound rules separately
	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(ingressRules),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return ec2.GetSecurityGroupRuleAttributeValue(fieldID, instance)
		},
	}, ingressOpts)
	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(egressRules),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return ec2.GetSecurityGroupRuleAttributeValue(fieldID, instance)
		},
	}, egressOpts)
	return nil
}
