package prefix_list

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

// prefixListShowFields returns the fields for the Prefix List detail table.
func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Prefix List ID", Category: "VPC", Visible: true},
		{Name: "Prefix List ARN", Category: "VPC", Visible: true},
		{Name: "Prefix List Name", Category: "VPC", Visible: true},
		{Name: "State", Category: "VPC", Visible: true},
		{Name: "Version", Category: "VPC", Visible: true},
		{Name: "Max Entries", Category: "VPC", Visible: true},
		{Name: "Address Family", Category: "VPC", Visible: true},
		{Name: "Owner", Category: "VPC", Visible: true},
	}
}

// showCmd is the cobra command for showing Prefix List details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a Prefix List",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowPrefixList(cmd, args[0]))
	},
}

// NewShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "horizontal")
}

// ShowPrefixList displays detailed information for a specified Prefix List.
func ShowPrefixList(cmd *cobra.Command, id string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create vpc service: %w", err)
	}

	ctx := context.TODO()
	pls, err := svc.GetManagedPrefixLists(ctx, &ascTypes.GetManagedPrefixListsInput{
		PrefixListIds: []string{id},
	})
	if err != nil {
		return fmt.Errorf("get prefix lists: %w", err)
	}
	if len(pls) == 0 {
		return fmt.Errorf("prefix list not found: %s", id)
	}
	pl := pls[0]

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:          "Prefix List summary for " + id,
		Columns:        3,
		MaxColumnWidth: 70,
	})

	fields, err := tablewriter.PopulateFieldValues(pl, getShowFields(), vpc.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	// Layout = Horizontal or Grid
	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}
	table.AddSections(tablewriter.BuildSections(fields, layout))

	tags, err := awsutil.PopulateTagFields(pl.Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags from prefix list: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}
