// ls.go defines the 'ls' subcommand for snapshot operations.
package snapshot

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	sortID      bool
	sortSize    bool
	showDesc    bool
	reverseSort bool
	owner       string
)

// Init function
func init() {
	NewLsFlags(lsCmd)
}

func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Snapshot ID", Category: "Snapshot Details", Visible: true, SortBy: sortID, SortDirection: tablewriter.Asc},
		{Name: "Volume Size", Category: "Snapshot Details", Visible: true},
		{Name: "Description", Category: "Snapshot Details", Visible: showDesc},
		{Name: "Tier", Category: "Snapshot Details", Visible: true},
		{Name: "State", Category: "Snapshot Details", Visible: true},
		{Name: "Started", Category: "Snapshot Details", Visible: true, DefaultSort: true},
		{Name: "Progress", Category: "Snapshot Details", Visible: true},
		{Name: "Encryption", Category: "Snapshot Details", Visible: true},
		{Name: "Data Transfer Progress", Category: "Snapshot Details", Visible: false},
		{Name: "KMS Key ID", Category: "Snapshot Details", Visible: false},
		{Name: "Owner ID", Category: "Snapshot Details", Visible: true},
	}
}

// lsCmd is the cobra command for listing snapshots.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all snapshots",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListSnapshots(cmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs snapshots in list format.")
	cobraCmd.Flags().BoolVarP(&sortID, "sort-id", "i", false, "Sort by descending snapshot ID.")
	cobraCmd.Flags().BoolVarP(&sortSize, "sort-size", "s", false, "Sort by descending snapshot size.")
	cobraCmd.Flags().BoolVarP(&showDesc, "show-description", "d", false, "Show the snapshot description column.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
	cobraCmd.Flags().StringVar(&owner, "owner", "", "Accepts a single AWS account ID or 'all' to show all snapshots. If not provided, only your own snapshots are shown.")
}

// ListSnapshots is the handler for the ls subcommand.
func ListSnapshots(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, ec2.NewEC2Service)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	ownerIds := getOwnerIds(owner)
	snapshots, err := svc.GetSnapshots(cmd.Context(), &ascTypes.GetSnapshotsInput{
		OwnerIds: ownerIds,
	})
	if err != nil {
		return fmt.Errorf("get snapshots: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Snapshots",
	})
	if list {
		table.SetStyle("plain")
	}
	fields := getListFields()
	fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(snapshots))

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(snapshots), fields, ec2.GetFieldValue, ec2.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)

	table.Render()
	return nil
}

// If 'all' is provided, dont use a filter
// If a specific owner is provided, use the owner-id filter
// Otherwise, use the self filter
func getOwnerIds(owner string) []string {
	ownerIds := []string{}
	if owner == "all" {
		// Do nothing
	} else if owner != "" {
		ownerIds = append(ownerIds, owner)
	} else {
		ownerIds = append(ownerIds, "self")
	}
	return ownerIds
}
