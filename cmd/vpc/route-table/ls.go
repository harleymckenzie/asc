package route_table

import (
	"context"
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

// routeTableListFields returns the fields for the Route Table list table.
func routeTableListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Route Table ID", Category: "VPC", Visible: true, SortBy: true, SortDirection: tablewriter.Asc},
		{Name: "Association Count", Category: "VPC", Visible: true},
		{Name: "Route Count", Category: "VPC", Visible: false},
		{Name: "Main", Category: "VPC", Visible: true},
		{Name: "VPC ID", Category: "VPC", Visible: true},
		{Name: "Owner", Category: "VPC", Visible: true},
	}
}

// routeTableRouteFields returns the fields for the Route Table route list table.
func routeTableRouteFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Destination", Category: "VPC", Visible: true},
		{Name: "Target", Category: "VPC", Visible: true},
		{Name: "Status", Category: "VPC", Visible: true},
		{Name: "Propagated", Category: "VPC", Visible: true},
	}
}

// lsCmd is the cobra command for listing Route Tables.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all Route Tables",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListRouteTables(cmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Route Tables in list format.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
}

// ListRouteTables is the handler for the ls subcommand.
func ListRouteTables(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create service: %w", err)
	}

	if len(args) > 0 {
		return ListRouteTableRules(cmd, args)
	} else {
		rts, err := svc.GetRouteTables(context.TODO(), &ascTypes.GetRouteTablesInput{})
		if err != nil {
			return fmt.Errorf("get route tables: %w", err)
		}

		table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
			Title: "Route Tables",
		})
		if list {
			table.SetRenderStyle("plain")
		}

		fields := routeTableListFields()
		fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(rts))

		headerRow := tablewriter.BuildHeaderRow(fields)
		table.AppendHeader(headerRow)
		table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(rts), fields, vpc.GetFieldValue, vpc.GetTagValue))
		table.SetFieldConfigs(fields, reverseSort)

		table.Render()
		return nil
	}
}

func ListRouteTableRules(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create service: %w", err)
	}

	rts, err := svc.GetRouteTables(context.TODO(), &ascTypes.GetRouteTablesInput{
		RouteTableIds: []string{args[0]},
	})
	if err != nil {
		return fmt.Errorf("get route tables: %w", err)
	}

	routes := rts[0].Routes

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: fmt.Sprintf("%s - Routes", args[0]),
	})
	if list {
		table.SetRenderStyle("plain")
	}

	fields := routeTableRouteFields()
	fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(routes))

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(routes), fields, vpc.GetFieldValue, vpc.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)

	table.Render()
	return nil
}
