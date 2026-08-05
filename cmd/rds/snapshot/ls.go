// ls.go lists RDS snapshots.

package snapshot

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/rds"
	ascTypes "github.com/harleymckenzie/asc/internal/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

// Variables
var (
	lsList    bool
	lsCluster bool

	sortIdentifier bool
	sortStatus     bool
	sortCreated    bool

	reverseSort bool
)

// Init function
func init() {
	newLsFlags(lsCmd)
}

// Column functions
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Snapshot Identifier", Category: "Snapshot", Visible: true, DefaultSort: true, SortBy: sortIdentifier, SortDirection: tablewriter.Asc},
		{Name: "Source Identifier", Category: "Snapshot", Visible: true},
		{Name: "Type", Category: "Snapshot", Visible: true},
		{Name: "Engine", Category: "Snapshot", Visible: true},
		{Name: "Status", Category: "Snapshot", Visible: true, SortBy: sortStatus, SortDirection: tablewriter.Asc},
		{Name: "Storage", Category: "Snapshot", Visible: true},
		{Name: "Encryption", Category: "Snapshot", Visible: true},
		{Name: "Created Time", Category: "Snapshot", Visible: true, SortBy: sortCreated, SortDirection: tablewriter.Desc},
	}
}

// Command variable
var lsCmd = &cobra.Command{
	Use:     "ls [identifier]",
	Short:   "List RDS snapshots",
	Aliases: []string{"list"},
	GroupID: "actions",
	Args:    cobra.MaximumNArgs(1),
	Example: `  asc rds snapshot ls                       # List all DB instance snapshots
  asc rds snapshot ls my-instance           # List snapshots for a specific instance
  asc rds snapshot ls --cluster             # List all DB cluster snapshots`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListRDSSnapshots(cmd, args))
	},
}

// Flag function
func newLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&lsList, "list", "l", false, "Outputs snapshots in list format.")
	cobraCmd.Flags().BoolVarP(&lsCluster, "cluster", "c", false, "List cluster snapshots instead of instance snapshots.")
	cmdutil.AddTagFlag(cobraCmd)

	cobraCmd.Flags().BoolVarP(&sortIdentifier, "sort-name", "n", false, "Sort by descending snapshot identifier.")
	cobraCmd.Flags().BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by descending snapshot status.")
	cobraCmd.Flags().BoolVarP(&sortCreated, "sort-created", "t", false, "Sort by descending snapshot creation time.")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-status", "sort-created")

	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")

	cobraCmd.Flags().SortFlags = false
}

// ListRDSSnapshots lists RDS DB instance or cluster snapshots.
func ListRDSSnapshots(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, rds.NewRDSService)
	if err != nil {
		return fmt.Errorf("create new RDS service: %w", err)
	}

	input := &ascTypes.GetSnapshotsInput{}
	if len(args) > 0 {
		input.Identifier = args[0]
	}

	var data []any
	if lsCluster {
		snapshots, err := svc.GetClusterSnapshots(cmd.Context(), input)
		if err != nil {
			return fmt.Errorf("list RDS cluster snapshots: %w", err)
		}
		data = utils.SlicesToAny(snapshots)
	} else {
		snapshots, err := svc.GetSnapshots(cmd.Context(), input)
		if err != nil {
			return fmt.Errorf("list RDS snapshots: %w", err)
		}
		data = utils.SlicesToAny(snapshots)
	}

	title := "DB Snapshots"
	if lsCluster {
		title = "DB Cluster Snapshots"
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         title,
		Style:         "rounded-separated",
		PlainStyle:    lsList,
		Fields:        getListFields(),
		Tags:          cmdutil.Tags,
		Data:          data,
		GetFieldValue: rds.GetFieldValue,
		GetTagValue:   rds.GetTagValue,
		ReverseSort:   reverseSort,
		HideEmpty:     true,
	})
	return nil
}
