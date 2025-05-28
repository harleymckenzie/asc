// show.go displays detailed information about a security group.
package security_group

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

// ec2SecurityGroupShowFields returns the fields for the security group detail table.
func ec2SecurityGroupShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Group Name", Display: true},
		{ID: "Group ID", Display: true},
		{ID: "Description", Display: true},
		{ID: "VPC ID", Display: true},
		{ID: "Owner ID", Display: true},
		{ID: "Ingress Count", Display: false},
		{ID: "Egress Count", Display: false},
		{ID: "Tag Count", Display: false},
	}
}

// showCmd is the cobra command for showing security group details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a security group",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowSecurityGroup(cmd, args[0]))
	},
}

// newShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {}

// ShowSecurityGroup displays detailed information for a specified security group.
func ShowSecurityGroup(cmd *cobra.Command, arg string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	groups, err := svc.GetSecurityGroups(
		ctx,
		&ascTypes.GetSecurityGroupsInput{GroupIDs: []string{arg}},
	)
	if err != nil {
		return fmt.Errorf("get security groups: %w", err)
	}
	if len(groups) == 0 {
		return fmt.Errorf("Security group not found: %s", arg)
	}

	fields := ec2SecurityGroupShowFields()
	opts := tableformat.RenderOptions{
		Title: "Security Group Details",
		Style: "rounded",
		Layout: tableformat.DetailTableLayout{
			Type: "vertical",
			ColumnsPerRow: 2,
		},
	}

	return tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: groups[0],
		Fields:   fields,
		GetAttribute: func(fieldID string, group any) (string, error) {
			return ec2.GetSecurityGroupAttributeValue(fieldID, group)
		},
	}, opts)
}
