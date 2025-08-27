// show.go displays detailed information about a security group.
package security_group

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {
	NewShowFlags(showCmd)
}

func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Group Name", Category: "Security Group Details", Visible: true},
		{Name: "Group ID", Category: "Security Group Details", Visible: true},
		{Name: "Description", Category: "Security Group Details", Visible: true},
		{Name: "VPC ID", Category: "Security Group Details", Visible: true},
		{Name: "Owner ID", Category: "Security Group Details", Visible: true},
		{Name: "Ingress Count", Category: "Security Group Details", Visible: false},
		{Name: "Egress Count", Category: "Security Group Details", Visible: false},
		{Name: "Tag Count", Category: "Security Group Details", Visible: false},
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
func NewShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "vertical")
}

// ShowSecurityGroup displays detailed information for a specified security group.
func ShowSecurityGroup(cmd *cobra.Command, arg string) error {
	svc, err := cmdutil.CreateService(cmd, ec2.NewEC2Service)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	groups, err := svc.GetSecurityGroups(context.TODO(), &ascTypes.GetSecurityGroupsInput{
		GroupIDs: []string{arg},
	})
	if err != nil {
		return fmt.Errorf("get security groups: %w", err)
	}
	if len(groups) == 0 {
		return fmt.Errorf("security group not found: %s", arg)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:          "Security Group details for " + *groups[0].GroupName,
		Columns:        3,
		MaxColumnWidth: 70,
	})
	fields, err := cmdutil.PopulateFieldValues(groups[0], getShowFields(), ec2.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	// Layout = Horizontal or Grid
	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}
	table.AddSections(tablewriter.BuildSections(fields, layout))
	tags, err := awsutil.PopulateTagFields(groups[0].Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags from instance: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}
