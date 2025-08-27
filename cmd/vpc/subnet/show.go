package subnet

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

// getShowFields returns the fields for the Subnet detail table.
func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Subnet ID", Category: "Subnet", Visible: true},
		{Name: "VPC ID", Category: "Subnet", Visible: true},
		{Name: "IPv4 CIDR", Category: "Subnet", Visible: true},
		{Name: "Availability Zone", Category: "Subnet", Visible: true},
		{Name: "State", Category: "Subnet", Visible: true},
		{Name: "Available IPs", Category: "Subnet", Visible: true},
		{Name: "Default subnet", Category: "Subnet", Visible: true},
		{Name: "Auto-assign public IPv4 address", Category: "Subnet", Visible: true},
		{Name: "Owner", Category: "Subnet", Visible: true},
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
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create vpc service: %w", err)
	}

	ctx := context.TODO()
	subnets, err := svc.GetSubnets(ctx, &ascTypes.GetSubnetsInput{SubnetIds: []string{id}})
	if err != nil {
		return fmt.Errorf("get subnets: %w", err)
	}
	if len(subnets) == 0 {
		return fmt.Errorf("subnet not found: %s", id)
	}
	subnet := subnets[0]

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:          "Subnet summary for " + id,
		Columns:        3,
		MaxColumnWidth: 70,
	})

	fields, err := cmdutil.PopulateFieldValues(subnet, getShowFields(), vpc.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	// Layout = Horizontal or Grid
	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}
	table.AddSections(tablewriter.BuildSections(fields, layout))

	tags, err := awsutil.PopulateTagFields(subnet.Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags from subnet: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}
