package route_table

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/format"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

// routeTableShowFields returns the fields for the Route Table detail table.
func routeTableShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Route Table ID", Display: true},
		{ID: "VPC ID", Display: true},
		{ID: "Main", Display: true},
		{ID: "Owner ID", Display: true},
		{ID: "Explicit subnet associations", Display: true},
		{ID: "Edge associations", Display: true},
		{ID: "Association Count", Display: true},
		{ID: "Route Count", Display: true},
		{ID: "Routes", Display: true},
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
func NewShowFlags(cobraCmd *cobra.Command) {}

// ShowRouteTable displays detailed information for a specified Route Table.
func ShowRouteTable(cmd *cobra.Command, id string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	rts, err := svc.GetRouteTables(ctx, &ascTypes.GetRouteTablesInput{RouteTableIDs: []string{id}})
	if err != nil {
		return fmt.Errorf("get route tables: %w", err)
	}
	if len(rts) == 0 {
		return fmt.Errorf("Route Table not found: %s", id)
	}
	rt := rts[0]

	fields := routeTableShowFields()
	opts := tableformat.RenderOptions{
		Title: fmt.Sprintf("Route Table Details\n(%s)", id),
		Style: "rounded",
	}

	return tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: rt,
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			// Placeholder logic for extra fields
			switch fieldID {
			case "Main":
				for _, assoc := range rt.Associations {
					if assoc.Main != nil && *assoc.Main {
						return "Yes", nil
					}
				}
				return "No", nil
			case "Owner ID":
				return format.StringOrEmpty(rt.OwnerId), nil
			case "Explicit subnet associations":
				return "-", nil // TODO: List explicit subnet associations
			case "Edge associations":
				return "-", nil // TODO: List edge associations
			case "Routes":
				return "-", nil // TODO: Format route list
			}
			return vpc.GetRouteTableAttributeValue(fieldID, instance)
		},
	}, opts)
}
