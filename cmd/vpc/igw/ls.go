// ls.go defines the 'ls' subcommand for volume operations.
package igw

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

// Define columns for volumes
func vpcIGWListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Internet Gateway ID", Display: true, DefaultSort: true},
		{ID: "State", Display: true},
		{ID: "VPC ID", Display: true},
		{ID: "Owner", Display: true},
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

// ListVolumes is the handler for the ls subcommand.
func ListIGWs(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)

	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	igws, err := svc.GetIGWs(ctx, &ascTypes.GetIGWsInput{})
	if err != nil {
		return fmt.Errorf("get internet gateways: %w", err)
	}

	fields := vpcIGWListFields()
	opts := tableformat.RenderOptions{
		Title:  "Internet Gateways",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(igws),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return vpc.GetIGWAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
