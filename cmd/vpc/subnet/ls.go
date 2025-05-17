package subnet

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

// subnetListFields returns the fields for the Subnet list table.
func subnetListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Subnet ID", Display: true, DefaultSort: true},
		{ID: "VPC ID", Display: true},
		{ID: "CIDR Block", Display: true},
		{ID: "Availability Zone", Display: true},
		{ID: "State", Display: true},
		{ID: "Available IPs", Display: true},
		{ID: "Default For AZ", Display: true},
	}
}

// lsCmd is the cobra command for listing Subnets.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all Subnets",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListSubnets(cmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Subnets in list format.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
}

// ListSubnets is the handler for the ls subcommand.
func ListSubnets(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	subnets, err := svc.GetSubnets(ctx, &ascTypes.GetSubnetsInput{})
	if err != nil {
		return fmt.Errorf("get subnets: %w", err)
	}

	fields := subnetListFields()
	opts := tableformat.RenderOptions{
		Title:  "Subnets",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(subnets),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return vpc.GetSubnetAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
