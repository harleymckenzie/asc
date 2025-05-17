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

// natGatewayShowFields returns the fields for the NAT Gateway detail table.
func natGatewayShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "NAT Gateway ID", Display: true},
		{ID: "NAT Gateway ARN", Display: true},
		{ID: "VPC ID", Display: true},
		{ID: "Subnet ID", Display: true},
		{ID: "Connectivity type", Display: true},
		{ID: "State", Display: true},
		{ID: "Primary public IPv4 address", Display: true},
		{ID: "Primary private IPv4 address", Display: true},
		{ID: "Created", Display: true},
		{ID: "Owner", Display: true},
		{ID: "IP Addresses", Display: true},
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
func NewShowFlags(cobraCmd *cobra.Command) {}

// ShowNatGateway displays detailed information for a specified NAT Gateway.
func ShowNatGateway(cmd *cobra.Command, id string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	nats, err := svc.GetNatGateways(ctx, &ascTypes.GetNatGatewaysInput{NatGatewayIDs: []string{id}})
	if err != nil {
		return fmt.Errorf("get nat gateways: %w", err)
	}
	if len(nats) == 0 {
		return fmt.Errorf("NAT Gateway not found: %s", id)
	}
	nat := nats[0]

	fields := natGatewayShowFields()
	opts := tableformat.RenderOptions{
		Title: fmt.Sprintf("NAT Gateway Details\n(%s)", id),
		Style: "rounded",
	}

	return tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: nat,
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			// Placeholder logic for extra fields
			switch fieldID {
			case "NAT Gateway ARN":
				return "-", nil // TODO: Lookup ARN
			case "Connectivity type":
				return string(nat.ConnectivityType), nil
			case "Primary public IPv4 address":
				return "-", nil // TODO: Lookup primary public IP
			case "Primary private IPv4 address":
				return "-", nil // TODO: Lookup primary private IP
			case "Created":
				return "-", nil // TODO: Format created time
			case "Owner":
				return "-", nil // TODO: Lookup Owner
			}
			return vpc.GetNatGatewayAttributeValue(fieldID, instance)
		},
	}, opts)
}
