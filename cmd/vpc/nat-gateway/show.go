package nat_gateway

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {
	NewShowFlags(showCmd)
}

// natGatewayShowFields returns the fields for the NAT Gateway detail table.
func natGatewayShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "NAT Gateway ID", Display: true},
		{ID: "VPC ID", Display: true},
		{ID: "Subnet ID", Display: true},
		{ID: "Connectivity", Display: true},
		{ID: "State", Display: true},
		{ID: "Primary Public IP", Display: true},
		{ID: "Primary Private IP", Display: true},
		{ID: "Created", Display: true},
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
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	if cmd.Flags().Changed("output") {
		if err := cmdutil.ValidateFlagChoice(cmd, "output", cmdutil.ValidLayouts); err != nil {
			return err
		}
	}

	nats, err := svc.GetNatGateways(ctx, &ascTypes.GetNatGatewaysInput{NatGatewayIds: []string{id}})
	if err != nil {
		return fmt.Errorf("get nat gateways: %w", err)
	}
	if len(nats) == 0 {
		return fmt.Errorf("nat gateway not found: %s", id)
	}
	nat := nats[0]

	fields := natGatewayShowFields()
	opts := tableformat.RenderOptions{
		Title: fmt.Sprintf("NAT Gateway Details\n(%s)", id),
		Style: "rounded",
		Layout: tableformat.DetailTableLayout{
			Type: cmdutil.GetLayout(cmd),
		},
	}

	return tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: nat,
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return vpc.GetNatGatewayAttributeValue(fieldID, instance)
		},
	}, opts)
}
