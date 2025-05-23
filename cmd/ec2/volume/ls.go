// ls.go defines the 'ls' subcommand for volume operations.
package volume

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
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
func ec2VolumeListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Volume ID", Display: true, DefaultSort: true},
		{ID: "Type", Display: true, Sort: sortType},
		{ID: "Size", Display: true},
		{ID: "Size Raw", Hidden: true, Sort: sortSize, SortDirection: "desc"}, // This is used for sorting, as Size is a combination of numbers and letters
		{ID: "IOPS", Display: true},
		{ID: "Throughput", Display: true},
		{ID: "Snapshot ID", Display: true},
		{ID: "State", Display: true},
		{ID: "Created", Display: showCreatedAt, Sort: sortCreatedAt, SortDirection: "desc"},
		{ID: "Attach Time", Display: showAttachTime, Sort: sortAttachTime, SortDirection: "desc"},
		{ID: "Availability Zone", Display: false},
		{ID: "Encryption", Display: true},
		{ID: "Fast Snapshot Restored", Display: true},
		{ID: "Multi-Attach Enabled", Display: false},
		{ID: "KMS Key ID", Display: showKMS},
	}
}

// lsCmd is the cobra command for listing volumes.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all volumes",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListVolumes(cobraCmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().
		BoolVarP(&list, "list", "l", false, "Outputs volumes in list format.")
	cobraCmd.Flags().BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending volume type.")
	cobraCmd.Flags().BoolVarP(&showKMS, "show-kms", "K", false, "Show the KMS Key ID column.")
	cobraCmd.Flags().
		BoolVarP(&sortState, "sort-state", "S", false, "Sort by descending volume state.")
	cobraCmd.Flags().
		BoolVarP(&sortAttachTime, "sort-attach-time", "a", false, "Sort by descending attach time.")
	cobraCmd.Flags().
		BoolVarP(&sortSize, "sort-size", "s", false, "Sort by descending size.")
	cobraCmd.Flags().
		BoolVarP(&sortCreatedAt, "sort-created-at", "t", false, "Sort by descending creation time.")
	cobraCmd.Flags().
		BoolVarP(&showAttachTime, "show-attach-time", "A", false, "Show the attach time column.")
	cobraCmd.Flags().
		BoolVarP(&showCreatedAt, "show-created-at", "C", false, "Show the creation time column.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
}

// ListVolumes is the handler for the ls subcommand.
func ListVolumes(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)

	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	volumes, err := svc.GetVolumes(ctx, &ascTypes.GetVolumesInput{})
	if err != nil {
		return fmt.Errorf("get volumes: %w", err)
	}

	fields := ec2VolumeListFields()
	opts := tableformat.RenderOptions{
		Title:  "Volumes",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(volumes),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return ec2.GetVolumeAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
