package prefix_list

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
	sortId      bool
)

func init() {
	NewLsFlags(lsCmd)
}

// getListFields returns the fields for the Prefix List list table.
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Prefix List ID", Category: "Details", Visible: true, DefaultSort: true, SortBy: sortId, SortDirection: tablewriter.Asc},
		{Name: "Prefix List Name", Category: "Details", Visible: true},
		{Name: "Max Entries", Category: "Details", Visible: false},
		{Name: "Address Family", Category: "Details", Visible: true},
		{Name: "State", Category: "Details", Visible: true},
		{Name: "Version", Category: "Details", Visible: false},
		{Name: "Prefix List ARN", Category: "Details", Visible: false},
		{Name: "Owner", Category: "Details", Visible: true},
	}
}

func prefixListEntriesFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "CIDR", Category: "VPC", Visible: true},
		{Name: "Description", Category: "VPC", Visible: true},
	}
}

// lsCmd is the cobra command for listing Prefix Lists.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all Prefix Lists",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListPrefixLists(cmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Prefix Lists in list format.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
	cobraCmd.Flags().BoolVarP(&sortId, "sort-id", "i", false, "Sort by descending prefix list ID.")
}

// ListPrefixLists is the handler for the ls subcommand.
func ListPrefixLists(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create vpc service: %w", err)
	}

	if len(args) > 0 {
		return ListPrefixListEntries(cmd, args)
	}

	pls, err := svc.GetManagedPrefixLists(context.TODO(), &ascTypes.GetManagedPrefixListsInput{})
	if err != nil {
		return fmt.Errorf("get prefix lists: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Prefix Lists",
	})
	if list {
		table.SetRenderStyle("plain")
	}

	fields := getListFields()
	fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(pls))

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(pls), fields, vpc.GetFieldValue, vpc.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)

	table.Render()
	return nil
}

func ListPrefixListEntries(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create vpc service: %w", err)
	}

	pl, err := svc.GetPrefixLists(context.TODO(), &ascTypes.GetPrefixListsInput{
		PrefixListIds: []string{args[0]},
	})
	if err != nil {
		return fmt.Errorf("get prefix list: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: fmt.Sprintf("%s - Entries", args[0]),
	})

	fields := prefixListEntriesFields()
	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(pl), fields, vpc.GetFieldValue, vpc.GetTagValue))

	table.Render()
	return nil
}
