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
	sortID         bool
	sortType       bool
	sortSize       bool
	sortState      bool
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
		{ID: "Volume ID", Visible: true, Sort: sortID},
		{ID: "Type", Visible: true, Sort: sortType},
		{ID: "Size", Visible: true, Sort: sortSize},
		{ID: "IOPS", Visible: true},
		{ID: "Throughput", Visible: true},
		{ID: "Snapshot ID", Visible: true},
		{ID: "State", Visible: true},
		{ID: "Created", Visible: showCreatedAt, Sort: sortCreatedAt},
		{ID: "Attach Time", Visible: showAttachTime, Sort: sortAttachTime, DefaultSort: true},
		{ID: "Availability Zone", Visible: false},
		{ID: "Encryption", Visible: true},
		{ID: "Fast Snapshot Restored", Visible: true},
		{ID: "Multi-Attach Enabled", Visible: false},
		{ID: "KMS Key ID", Visible: showKMS},
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
	cobraCmd.Flags().BoolVarP(&sortID, "sort-id", "i", false, "Sort by descending volume ID.")
	cobraCmd.Flags().BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending volume type.")
	cobraCmd.Flags().BoolVarP(&sortSize, "sort-size", "s", false, "Sort by descending volume size.")
	cobraCmd.Flags().BoolVarP(&showKMS, "show-kms", "k", false, "Show the KMS Key ID column.")
	cobraCmd.Flags().
		BoolVarP(&sortState, "sort-state", "S", false, "Sort by descending volume state.")
	cobraCmd.Flags().
		BoolVarP(&sortAttachTime, "sort-attach-time", "a", false, "Sort by descending attach time.")
	cobraCmd.Flags().
		BoolVarP(&sortCreatedAt, "sort-created-at", "t", false, "Sort by descending creation time.")
	cobraCmd.Flags().
		BoolVarP(&showCreatedAt, "show-created-at", "c", false, "Show the creation time column.")
	cobraCmd.Flags().
		BoolVarP(&showAttachTime, "show-attach-time", "A", false, "Show the attach time column.")
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
