// rm.go deletes RDS snapshots.

package snapshot

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/harleymckenzie/asc/internal/service/rds"
	ascTypes "github.com/harleymckenzie/asc/internal/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// Variables
var (
	rmForce   bool
	rmCluster bool
)

// Init function
func init() {
	newRmFlags(rmCmd)
}

// Command variable
var rmCmd = &cobra.Command{
	Use:     "rm <snapshot-name> [snapshot-name...]",
	Short:   "Delete RDS snapshots",
	Aliases: []string{"remove", "delete"},
	GroupID: "actions",
	Args:    cobra.MinimumNArgs(1),
	Example: `  asc rds snapshot rm my-snapshot                  # Delete an instance snapshot
  asc rds snapshot rm my-snapshot --force          # Delete without confirmation
  asc rds snapshot rm my-cluster-snapshot --cluster # Delete a cluster snapshot`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(DeleteRDSSnapshot(cmd, args))
	},
}

// Flag function
func newRmFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&rmForce, "force", "f", false, "Skip confirmation prompt.")
	cobraCmd.Flags().BoolVarP(&rmCluster, "cluster", "c", false, "Delete cluster snapshots instead of instance snapshots.")
}

// DeleteRDSSnapshot deletes one or more RDS snapshots, with optional confirmation.
func DeleteRDSSnapshot(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, rds.NewRDSService)
	if err != nil {
		return fmt.Errorf("create new RDS service: %w", err)
	}

	snapshotType := "snapshot"
	if rmCluster {
		snapshotType = "cluster snapshot"
	}

	// Confirmation prompt unless --force
	if !rmForce {
		fmt.Printf("The following %s(s) will be deleted:\n", snapshotType)
		for _, name := range args {
			fmt.Printf("  - %s\n", name)
		}
		fmt.Printf("\nDelete %d %s(s)? [y/N]: ", len(args), snapshotType)

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read confirmation: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	for _, name := range args {
		err = svc.DeleteSnapshot(cmd.Context(), &ascTypes.DeleteSnapshotInput{
			SnapshotIdentifier: name,
			IsCluster:          rmCluster,
		})
		if err != nil {
			return fmt.Errorf("delete %s %s: %w", snapshotType, name, err)
		}
		fmt.Printf("Deleted: %s\n", name)
	}

	return nil
}
