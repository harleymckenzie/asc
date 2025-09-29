package nat_gateway

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

// natGatewayShowFields returns the fields for the NAT Gateway detail table.
func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "NAT Gateway ID", Category: "VPC", Visible: true},
		{Name: "VPC ID", Category: "VPC", Visible: true},
		{Name: "Subnet ID", Category: "VPC", Visible: true},
		{Name: "Connectivity", Category: "VPC", Visible: true},
		{Name: "State", Category: "VPC", Visible: true},
		{Name: "Primary Public IP", Category: "VPC", Visible: true},
		{Name: "Primary Private IP", Category: "VPC", Visible: true},
		{Name: "Created", Category: "VPC", Visible: true},
	}
}

// showCmd is the cobra command for showing NAT Gateway details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a NAT Gateway",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowNatGateway(cmd, args[0]))
	},
}

// NewShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "vertical")
}

// ShowNatGateway displays detailed information for a specified NAT Gateway.
func ShowNatGateway(cmd *cobra.Command, id string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create service: %w", err)
	}

	if cmd.Flags().Changed("output") {
		if err := cmdutil.ValidateFlagChoice(cmd, "output", cmdutil.ValidLayouts); err != nil {
			return err
		}
	}

	nats, err := svc.GetNatGateways(context.TODO(), &ascTypes.GetNatGatewaysInput{NatGatewayIds: []string{id}})
	if err != nil {
		return fmt.Errorf("get nat gateways: %w", err)
	}
	if len(nats) == 0 {
		return fmt.Errorf("nat gateway not found: %s", id)
	}
	nat := nats[0]

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:          fmt.Sprintf("NAT Gateway Details (%s)", id),
		Columns:        3,
		MaxColumnWidth: 70,
	})

	fields, err := tablewriter.PopulateFieldValues(nat, getShowFields(), vpc.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))
	tags, err := awsutil.PopulateTagFields(nat.Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}
