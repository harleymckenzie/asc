// The show command displays detailed information about an VPC internet gateway.

package igw

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

// Column functions
func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Internet Gateway ID", Category: "VPC", Visible: true},
		{Name: "State", Category: "VPC", Visible: true},
		{Name: "VPC ID", Category: "VPC", Visible: true},
		{Name: "Owner", Category: "VPC", Visible: true},
	}
}

// Command variable
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an internet gateway",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowVPCIGW(cmd, args[0]))
	},
}

// Flag function
func NewShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "vertical")
}

// ShowVPCIGW is the function for showing VPC internet gateways
func ShowVPCIGW(cmd *cobra.Command, id string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create vpc service: %w", err)
	}

	ctx := context.TODO()
	igws, err := svc.GetIGWs(ctx, &ascTypes.GetIGWsInput{IGWIds: []string{id}})
	if err != nil {
		return fmt.Errorf("get internet gateways: %w", err)
	}
	if len(igws) == 0 {
		return fmt.Errorf("internet gateway not found: %s", id)
	}
	igw := igws[0]

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:          "Internet Gateway summary for " + id,
		Columns:        3,
		MaxColumnWidth: 70,
	})

	fields, err := tablewriter.PopulateFieldValues(igw, getShowFields(), vpc.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	// Layout = Horizontal or Grid
	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}
	table.AddSections(tablewriter.BuildSections(fields, layout))

	tags, err := awsutil.PopulateTagFields(igw.Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags from internet gateway: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}
