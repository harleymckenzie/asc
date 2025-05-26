package route_table

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

// routeTableListFields returns the fields for the Route Table list table.
func routeTableListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Route Table ID", Display: true, DefaultSort: true},
		{ID: "Association Count", Display: true},
		{ID: "Route Count", Display: false},
		{ID: "Main", Display: true},
		{ID: "VPC ID", Display: true},
		{ID: "Owner", Display: true},
	}
}

// routeTableRouteFields returns the fields for the Route Table route list table.
func routeTableRouteFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Destination", Display: true},
		{ID: "Target", Display: true},
		{ID: "Status", Display: true},
		{ID: "Propagated", Display: true},
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
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	if len(args) > 0 {
		return ListRouteTableRules(cmd, args)
	} else {
		rts, err := svc.GetRouteTables(ctx, &ascTypes.GetRouteTablesInput{})
		if err != nil {
			return fmt.Errorf("get route tables: %w", err)
		}

		fields := routeTableListFields()
		opts := tableformat.RenderOptions{
			Title:  "Route Tables",
			Style:  "rounded",
			SortBy: tableformat.GetSortByField(fields, reverseSort),
		}

		if list {
			opts.Style = "list"
		}

		tableformat.RenderTableList(&tableformat.ListTable{
			Instances: utils.SlicesToAny(rts),
			Fields:    fields,
			GetAttribute: func(fieldID string, instance any) (string, error) {
				return vpc.GetRouteTableAttributeValue(fieldID, instance)
			},
		}, opts)
		return nil
	}
}

func ListRouteTableRules(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	rts, err := svc.GetRouteTables(ctx, &ascTypes.GetRouteTablesInput{
		RouteTableIds: []string{args[0]},
	})
	if err != nil {
		return fmt.Errorf("get route tables: %w", err)
	}

	routes := rts[0].Routes

	fields := routeTableRouteFields()
	opts := tableformat.RenderOptions{
		Title:  fmt.Sprintf("%s - Routes", args[0]),
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(routes),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return vpc.GetRouteTableRouteAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
