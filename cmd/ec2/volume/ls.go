// ls.go defines the 'ls' subcommand for volume operations.
package volume

import (
	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
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

// Define columns for volumes
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Volume ID", Category: "Volume Details", Visible: true, SortBy: true, SortDirection: tablewriter.Asc},
		{Name: "Type", Category: "Volume Details", Visible: true, SortBy: sortType, SortDirection: tablewriter.Asc},
		{Name: "Size", Category: "Volume Details", Visible: true},
		{Name: "Size Raw", Category: "Volume Details", Visible: false, SortBy: sortSize, SortDirection: tablewriter.Desc},
		{Name: "IOPS", Category: "Volume Details", Visible: true},
		{Name: "Throughput", Category: "Volume Details", Visible: true},
		{Name: "Snapshot ID", Category: "Volume Details", Visible: true},
		{Name: "State", Category: "Volume Details", Visible: true, SortBy: sortState, SortDirection: tablewriter.Asc},
		{Name: "Created", Category: "Volume Details", Visible: showCreatedAt, DefaultSort: true, SortBy: sortCreatedAt, SortDirection: tablewriter.Desc},
		{Name: "Attach Time", Category: "Volume Details", Visible: showAttachTime, SortBy: sortAttachTime, SortDirection: tablewriter.Desc},
		{Name: "Availability Zone", Category: "Volume Details", Visible: false},
		{Name: "Encryption", Category: "Volume Details", Visible: true},
		{Name: "Fast Snapshot Restored", Category: "Volume Details", Visible: false},
		{Name: "Multi-Attach Enabled", Category: "Volume Details", Visible: false},
		{Name: "KMS Key ID", Category: "Volume Details", Visible: showKMS},
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
	svc, err := cmdutil.CreateService(cmd, ec2.NewEC2Service)
	if err != nil {
		return err
	}

	volumes, err := svc.GetVolumes(cmd.Context(), &ascTypes.GetVolumesInput{})
	if err != nil {
		return err
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Volumes",
	})
	if list {
		table.SetRenderStyle("plain")
	}

	fields := getListFields()
	fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(volumes))

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(volumes), fields, ec2.GetFieldValue, ec2.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)
	table.Render()
	return nil
}
