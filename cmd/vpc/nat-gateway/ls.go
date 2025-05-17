package nat_gateway

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	reverseSort bool
)

func init() {
	NewLsFlags(lsCmd)
}

// natGatewayListFields returns the fields for the NAT Gateway list table.
func natGatewayListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "NAT Gateway ID", Display: true, DefaultSort: true},
		{ID: "VPC ID", Display: true},
		{ID: "Subnet ID", Display: true},
		{ID: "State", Display: true},
		{ID: "Type", Display: true},
		{ID: "IP Addresses", Display: true},
	}
}

// lsCmd is the cobra command for listing NAT Gateways.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all NAT Gateways",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListNatGateways(cmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs NAT Gateways in list format.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
}

// ListNatGateways is the handler for the ls subcommand.
func ListNatGateways(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	nats, err := svc.GetNatGateways(ctx, &ascTypes.GetNatGatewaysInput{})
	if err != nil {
		return fmt.Errorf("get nat gateways: %w", err)
	}

	fields := natGatewayListFields()
	opts := tableformat.RenderOptions{
		Title:  "NAT Gateways",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(nats),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return vpc.GetNatGatewayAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
