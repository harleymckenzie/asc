package nacl

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {
	NewShowFlags(showCmd)
}

// getShowFields returns the fields for the NACL detail table.
func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Network ACL ID", Category: "Details", Visible: true},
		{Name: "VPC ID", Category: "Details", Visible: true},
		{Name: "Default", Category: "Details", Visible: true},
		{Name: "Owner", Category: "Details", Visible: true},
		{Name: "Associated with", Category: "Association", Visible: true},
		{Name: "Inbound Rules", Category: "Rules", Visible: true},
		{Name: "Outbound Rules", Category: "Rules", Visible: true},
	}
}

// showCmd is the cobra command for showing NACL details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a Network ACL",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowNACL(cmd, args[0]))
	},
}

// NewShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "vertical")
}

// ShowNACL displays detailed information for a specified NACL.
func ShowNACL(cmd *cobra.Command, id string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create vpc service: %w", err)
	}

	nacls, err := svc.GetNACLs(context.TODO(), &ascTypes.GetNACLsInput{
		NetworkAclIds: []string{id},
	})
	if err != nil {
		return fmt.Errorf("get network acls: %w", err)
	}

	if len(nacls) == 0 {
		return fmt.Errorf("network ACL not found: %s", id)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:          "Network ACL Details (" + id + ")",
		Columns:        3,
		MaxColumnWidth: 70,
	})

	fields, err := tablewriter.PopulateFieldValues(nacls[0], getShowFields(), vpc.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))
	tags, err := awsutil.PopulateTagFields(nacls[0].Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}
