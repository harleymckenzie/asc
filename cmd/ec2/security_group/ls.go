// ls.go defines the 'ls' subcommand for security group operations.
package security_group

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/pkg/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/harleymckenzie/asc/pkg/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	sortID      bool
	sortName    bool
	showDesc    bool
	reverseSort bool
)

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs security groups in list format.")
	cobraCmd.Flags().BoolVarP(&sortID, "sort-id", "i", false, "Sort by descending group ID.")
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending group name.")
	cobraCmd.Flags().
		BoolVarP(&showDesc, "show-description", "d", false, "Show the security group description column.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
}

// ec2SecurityGroupListFields returns the fields for the security group list table.
func ec2SecurityGroupListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Group ID", Visible: true, Sort: sortID},
		{ID: "Group Name", Visible: true, Sort: sortName},
		{ID: "Description", Visible: showDesc},
		{ID: "VPC ID", Visible: true},
		{ID: "Owner ID", Visible: true},
		{ID: "Ingress Count", Visible: true},
		{ID: "Egress Count", Visible: true},
		{ID: "Tag Count", Visible: true},
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

// ListSecurityGroups is the handler for the ls subcommand.
func ListSecurityGroups(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	groups, err := svc.GetSecurityGroups(ctx, &ascTypes.GetSecurityGroupsInput{})
	if err != nil {
		return fmt.Errorf("get security groups: %w", err)
	}

	fields := ec2SecurityGroupListFields()
	opts := tableformat.RenderOptions{
		Title:  "Security Groups",
		Style:  "list",
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
