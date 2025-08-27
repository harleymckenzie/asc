package route_table

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

// routeTableShowFields returns the fields for the Route Table detail table.
func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Route Table ID", Category: "VPC", Visible: true},
		{Name: "VPC ID", Category: "VPC", Visible: true},
		{Name: "Main", Category: "VPC", Visible: true},
		{Name: "Owner", Category: "VPC", Visible: true},
		{Name: "Association Count", Category: "VPC", Visible: true},
		{Name: "Route Count", Category: "VPC", Visible: true},
	}
}

// showCmd is the cobra command for showing Route Table details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a Route Table",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowRouteTable(cmd, args[0]))
	},
}

// NewShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "vertical")
}

// ShowRouteTable displays detailed information for a specified Route Table.
func ShowRouteTable(cmd *cobra.Command, id string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create vpc service: %w", err)
	}

	ctx := context.TODO()
	rts, err := svc.GetRouteTables(ctx, &ascTypes.GetRouteTablesInput{RouteTableIds: []string{id}})
	if err != nil {
		return fmt.Errorf("get route tables: %w", err)
	}
	if len(rts) == 0 {
		return fmt.Errorf("route table not found: %s", id)
	}
	rt := rts[0]

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:          "Route Table summary for " + id,
		Columns:        3,
		MaxColumnWidth: 70,
	})

	fields, err := cmdutil.PopulateFieldValues(rt, getShowFields(), vpc.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	// Layout = Horizontal or Grid
	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}
	table.AddSections(tablewriter.BuildSections(fields, layout))

	tags, err := awsutil.PopulateTagFields(rt.Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags from route table: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}
