package nat_gateway

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
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
func natGatewayListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "NAT Gateway ID", Category: "VPC", Visible: true, SortBy: true, SortDirection: tablewriter.Asc},
		{Name: "Connectivity", Category: "VPC", Visible: true},
		{Name: "State", Category: "VPC", Visible: true},
		{Name: "VPC ID", Category: "VPC", Visible: true},
		{Name: "Subnet ID", Category: "VPC", Visible: true},
		{Name: "Primary Public IP", Category: "VPC", Visible: true},
		{Name: "Primary Private IP", Category: "VPC", Visible: false},
		{Name: "Created", Category: "VPC", Visible: false},
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
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create service: %w", err)
	}

	nats, err := svc.GetNatGateways(cmd.Context(), &ascTypes.GetNatGatewaysInput{})
	if err != nil {
		return fmt.Errorf("get nat gateways: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "NAT Gateways",
	})
	if list {
		table.SetRenderStyle("plain")
	}

	fields := natGatewayListFields()
	fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(nats))

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(nats), fields, vpc.GetFieldValue, vpc.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)

	table.Render()
	return nil
}
