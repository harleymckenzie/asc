// ls.go defines the 'ls' subcommand for volume operations.
package igw

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

// Variables
var (
	list           bool
	sortType       bool
	sortState      bool
	sortSize       bool
	sortAttachTime bool
	sortCreatedAt  bool
	showKMS        bool
	showCreatedAt  bool
	showAttachTime bool
	reverseSort    bool
)

// Init function
func init() {
	NewLsFlags(lsCmd)
}

// Define columns for Internet Gateways
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Internet Gateway ID", Category: "IGW", Visible: true, SortBy: true, SortDirection: tablewriter.Asc},
		{Name: "State", Category: "IGW", Visible: true, SortBy: sortState, SortDirection: tablewriter.Asc},
		{Name: "VPC ID", Category: "IGW", Visible: true},
		{Name: "Owner", Category: "IGW", Visible: true},
	}
}

// lsCmd is the cobra command for listing volumes.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all Internet Gateways",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListIGWs(cobraCmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs volumes in list format.")
	cobraCmd.Flags().BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending volume type.")
	cobraCmd.Flags().BoolVarP(&showKMS, "show-kms", "K", false, "Show the KMS Key ID column.")
	cobraCmd.Flags().BoolVarP(&sortState, "sort-state", "S", false, "Sort by descending volume state.")
	cobraCmd.Flags().BoolVarP(&sortAttachTime, "sort-attach-time", "a", false, "Sort by descending attach time.")
	cobraCmd.Flags().BoolVarP(&sortSize, "sort-size", "s", false, "Sort by descending size.")
	cobraCmd.Flags().BoolVarP(&sortCreatedAt, "sort-created-at", "t", false, "Sort by descending creation time.")
	cobraCmd.Flags().BoolVarP(&showAttachTime, "show-attach-time", "A", false, "Show the attach time column.")
	cobraCmd.Flags().BoolVarP(&showCreatedAt, "show-created-at", "C", false, "Show the creation time column.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
}

// ListIGWs is the handler for the ls subcommand.
func ListIGWs(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	igws, err := svc.GetIGWs(cmd.Context(), &ascTypes.GetIGWsInput{})
	if err != nil {
		return fmt.Errorf("get internet gateways: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Internet Gateways",
	})
	if list {
		table.SetRenderStyle("plain")
	}

	fields := getListFields()
	fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(igws))

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(igws), fields, vpc.GetIGWFieldValue, vpc.GetIGWTagValue))
	table.SortBy(fields, reverseSort)
	table.Render()
	return nil
}
