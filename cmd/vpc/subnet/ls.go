package subnet

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
	sortId      bool
)

func init() {
	NewLsFlags(lsCmd)
}

// getListFields returns the fields for the Subnet list table.
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Subnet ID", Category: "Subnet", Visible: true, DefaultSort: true, SortBy: sortId, SortDirection: tablewriter.Asc},
		{Name: "VPC ID", Category: "Subnet", Visible: true},
		{Name: "CIDR Block", Category: "Subnet", Visible: true},
		{Name: "Availability Zone", Category: "Subnet", Visible: true},
		{Name: "State", Category: "Subnet", Visible: true},
		{Name: "Available IPs", Category: "Subnet", Visible: true},
		{Name: "Default For AZ", Category: "Subnet", Visible: true},
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
	cobraCmd.Flags().BoolVarP(&sortId, "sort-id", "i", false, "Sort by descending subnet ID.")
}

// ListSubnets is the handler for the ls subcommand.
func ListSubnets(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create service: %w", err)
	}

	subnets, err := svc.GetSubnets(cmd.Context(), &ascTypes.GetSubnetsInput{})
	if err != nil {
		return fmt.Errorf("get subnets: %w", err)
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         "Subnets",
		PlainStyle:    list,
		Fields:        getListFields(),
		Tags:          cmdutil.Tags,
		Data:          utils.SlicesToAny(subnets),
		GetFieldValue: vpc.GetFieldValue,
		GetTagValue:   vpc.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}
