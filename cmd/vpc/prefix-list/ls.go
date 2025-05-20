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
		{ID: "Prefix List Name", Display: true},
		{ID: "Max Entries", Display: false},
		{ID: "Address Family", Display: true},
		{ID: "State", Display: true},
		{ID: "Version", Display: false},
		{ID: "ARN", Display: false},
		{ID: "Owner", Display: true},
	}
}

func prefixListEntriesFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "CIDR", Display: true},
		{ID: "Description", Display: true},
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

	if len(args) > 0 {
		return ListPrefixListEntries(cmd, args)
	} else {
		pls, err := svc.GetManagedPrefixLists(ctx, &ascTypes.GetManagedPrefixListsInput{})
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
}

func ListPrefixListEntries(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}
	
	pl, err := svc.GetPrefixLists(ctx, &ascTypes.GetPrefixListsInput{
		PrefixListIds: []string{args[0]},
	})
	if err != nil {
		return fmt.Errorf("get prefix list: %w", err)
	}

	fields := prefixListEntriesFields()
	opts := tableformat.RenderOptions{
		Title:  fmt.Sprintf("%s - Entries", args[0]),
		Style:  "rounded",
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(pl),
		Fields:    fields,
	}, opts)
	return nil
}
