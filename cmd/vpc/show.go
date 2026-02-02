package vpc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
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

// showCmd is the cobra command for showing VPC details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a VPC",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(showVPC(cmd, args[0]))
	},
}

// NewShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "vertical")
}

func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "VPC ID", Category: "VPC", Visible: true},
		{Name: "State", Category: "VPC", Visible: true},
		{Name: "Tenancy", Category: "VPC", Visible: true},
		{Name: "Default VPC", Category: "VPC", Visible: true},
		{Name: "DHCP Option Set", Category: "VPC", Visible: true},
		{Name: "Main Route Table", Category: "VPC", Visible: true},
		{Name: "Main Network ACL", Category: "VPC", Visible: true},
		{Name: "IPv4 CIDR", Category: "VPC", Visible: true},
		{Name: "IPv6 CIDR", Category: "VPC", Visible: true},
		{Name: "Owner ID", Category: "VPC", Visible: true},
	}
}

// ShowVPC displays detailed information for a specified VPC.
func showVPC(cmd *cobra.Command, id string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create vpc service: %w", err)
	}

	ctx := cmd.Context()
	vpcs, err := svc.GetVPCs(ctx, &ascTypes.GetVPCsInput{})
	if err != nil {
		return fmt.Errorf("get vpc: %w", err)
	}
	var v types.Vpc
	found := false
	for _, candidate := range vpcs {
		if candidate.VpcId != nil && *candidate.VpcId == id {
			v = candidate
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("VPC not found: %s", id)
	}

	routeTables, _ := svc.GetRouteTables(ctx, &ascTypes.GetRouteTablesInput{RouteTableIds: nil})
	networkAcls, _ := svc.GetNACLs(ctx, &ascTypes.GetNACLsInput{NetworkAclIds: nil})

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:          "VPC summary for " + id,
		Columns:        3,
		MaxColumnWidth: 70,
	})

	// Custom field value getter that handles special cases
	getFieldValue := func(fieldName string, instance any) (string, error) {
		if fieldName == "Main Route Table" {
			return vpc.FindMainRouteTable(id, routeTables), nil
		}
		if fieldName == "Main Network ACL" {
			return vpc.FindMainNetworkACL(id, networkAcls), nil
		}
		return vpc.GetFieldValue(fieldName, instance)
	}

	fields, err := tablewriter.PopulateFieldValues(v, getShowFields(), getFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	// Layout = Horizontal or Grid
	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}
	table.AddSections(tablewriter.BuildSections(fields, layout))

	tags, err := awsutil.PopulateTagFields(v.Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags from VPC: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}
