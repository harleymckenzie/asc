package nacl

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

// naclListFields returns the fields for the NACL list table.
func naclListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Network ACL ID", Display: true, DefaultSort: true},
		{ID: "VPC ID", Display: true},
		{ID: "Is Default", Display: true},
		{ID: "Entry Count", Display: true},
		{ID: "Association Count", Display: true},
	}
}

// lsCmd is the cobra command for listing NACLs.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all Network ACLs",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListNACLs(cmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs NACLs in list format.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
}

// ListNACLs is the handler for the ls subcommand.
func ListNACLs(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	nacls, err := svc.GetNACLs(ctx, &ascTypes.GetNACLsInput{})
	if err != nil {
		return fmt.Errorf("get network acls: %w", err)
	}

	fields := naclListFields()
	opts := tableformat.RenderOptions{
		Title:  "Network ACLs",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(nacls),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return vpc.GetNetworkAclAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
