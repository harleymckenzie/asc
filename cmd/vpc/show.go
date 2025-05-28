package vpc

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

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
func NewShowFlags(cobraCmd *cobra.Command) {}

func vpcShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "VPC ID", Display: true},
		{ID: "State", Display: true},
		{ID: "Tenancy", Display: true},
		{ID: "Default VPC", Display: true},
		{ID: "DHCP Option Set", Display: true},
		{ID: "Main Route Table", Display: true},
		{ID: "Main Network ACL", Display: true},
		{ID: "IPv4 CIDR", Display: true},
		{ID: "IPv6 CIDR", Display: true},
		{ID: "Owner ID", Display: true},
	}
}

// ShowVPC displays detailed information for a specified VPC.
func showVPC(cmd *cobra.Command, id string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)

	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

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

	fields := vpcShowFields()
	opts := tableformat.RenderOptions{
		Title: id + " VPC Details",
		Style: "rounded",
		Layout: tableformat.DetailTableLayout{
			Type:           "horizontal",
			ColumnsPerRow:  3,
			ColumnMinWidth: 20,
		},
	}

	return tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: v,
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			if fieldID == "Main Route Table" {
				return vpc.FindMainRouteTable(id, routeTables), nil
			}
			if fieldID == "Main Network ACL" {
				return vpc.FindMainNetworkACL(id, networkAcls), nil
			}
			return vpc.GetVPCAttributeValue(fieldID, instance)
		},
	}, opts)
}
