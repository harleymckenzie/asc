package subnet

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

// subnetShowFields returns the fields for the Subnet detail table.
func subnetShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Subnet ID", Display: true},
		{ID: "Subnet ARN", Display: true},
		{ID: "VPC ID", Display: true},
		// {ID: "IPv4 CIDR", Display: false}, // TODO: Implement in attributes table
		// {ID: "IPv6 CIDR", Display: false}, // TODO: Implement in attributes table
		{ID: "Availability Zone", Display: true},
		// {ID: "Availability Zone ID", Display: true}, // TODO: Implement in attributes table
		{ID: "Network ACL", Display: true},
		{ID: "Route Table", Display: true},
		{ID: "State", Display: true},
		{ID: "Owner", Display: true},
		{ID: "Default subnet", Display: true},
		{ID: "Auto-assign public IPv4 address", Display: true},
		{ID: "Auto-assign customer-owned IPv4 address", Display: true},
		{ID: "Customer-owned IPv4 pool", Display: true},
		{ID: "Outpost ID", Display: true},
		{ID: "Hostname type", Display: true},
		{ID: "DNS64", Display: true},
		{ID: "IPv6-only", Display: true},
		{ID: "Available IPs", Display: true},
	}
}

// showCmd is the cobra command for showing Subnet details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a Subnet",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowSubnet(cmd, args[0]))
	},
}

// NewShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "vertical")
}

// ShowSubnet displays detailed information for a specified Subnet.
func ShowSubnet(cmd *cobra.Command, id string) error {
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

	subnets, err := svc.GetSubnets(ctx, &ascTypes.GetSubnetsInput{SubnetIds: []string{id}})
	if err != nil {
		return fmt.Errorf("get subnets: %w", err)
	}
	if len(subnets) == 0 {
		return fmt.Errorf("Subnet not found: %s", id)
	}
	sub := subnets[0]

	fields := subnetShowFields()
	opts := tableformat.RenderOptions{
		Title: fmt.Sprintf("Subnet Details\n(%s)", id),
		Style: "rounded",
		Layout: tableformat.DetailTableLayout{
			Type: cmdutil.GetLayout(cmd),
		},
	}

	return tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: sub,
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			// Placeholder logic for extra fields
			switch fieldID {
			case "Subnet ARN":
				return "-", nil // TODO: Compute ARN
			case "Network ACL":
				return "-", nil // TODO: Lookup NACL
			case "Route Table":
				return "-", nil // TODO: Lookup Route Table
			case "Owner":
				return "-", nil // TODO: Lookup Owner
			case "Default subnet":
				return "No", nil // TODO: Compute
			case "Auto-assign public IPv4 address":
				return "-", nil // TODO: Lookup
			case "Auto-assign customer-owned IPv4 address":
				return "-", nil // TODO: Lookup
			case "Customer-owned IPv4 pool":
				return "-", nil // TODO: Lookup
			case "Outpost ID":
				return "-", nil // TODO: Lookup
			case "Hostname type":
				return "-", nil // TODO: Lookup
			case "DNS64":
				return "-", nil // TODO: Lookup
			case "IPv6-only":
				return "No", nil // TODO: Compute
			case "Available IPs":
				return fmt.Sprintf("%d", sub.AvailableIpAddressCount), nil
			}
			return vpc.GetSubnetAttributeValue(fieldID, instance)
		},
	}, opts)
}
