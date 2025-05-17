package prefix_list

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

// prefixListListFields returns the fields for the Prefix List list table.
func prefixListListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Prefix List ID", Display: true, DefaultSort: true},
		{ID: "Name", Display: true},
		{ID: "CIDRs", Display: true},
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
}

// ListPrefixLists is the handler for the ls subcommand.
func ListPrefixLists(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	pls, err := svc.GetPrefixLists(ctx, &ascTypes.GetPrefixListsInput{})
	if err != nil {
		return fmt.Errorf("get prefix lists: %w", err)
	}

	fields := prefixListListFields()
	opts := tableformat.RenderOptions{
		Title:  "Prefix Lists",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(pls),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return vpc.GetPrefixListAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
