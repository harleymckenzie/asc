// The snapshot command creates a snapshot of an RDS instance or cluster.

package rds

import (
	"fmt"
	"time"

	"github.com/harleymckenzie/asc/internal/service/rds"
	ascTypes "github.com/harleymckenzie/asc/internal/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// Variables
var (
	snapshotWait    bool
	snapshotCluster bool
)

// Init function
func init() {
	newSnapshotFlags(snapshotCmd)
}

// Flag function
func newSnapshotFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().SortFlags = false
	cobraCmd.Flags().BoolVarP(&snapshotWait, "wait", "w", false, "Wait for the snapshot to complete")
	cobraCmd.Flags().BoolVarP(&snapshotCluster, "cluster", "c", false, "Snapshot a cluster instead of an instance")
}

// Command variable
var snapshotCmd = &cobra.Command{
	Use:     "snapshot <identifier> <snapshot-name>",
	Short:   "Create a snapshot of an RDS instance or cluster",
	Long:    "Create a manual snapshot of an RDS instance or cluster. Use --cluster for cluster snapshots.",
	Args:    cobra.ExactArgs(2),
	GroupID: "actions",
	Example: `  asc rds snapshot my-instance my-snapshot               # Snapshot an instance
  asc rds snapshot my-cluster my-snapshot --cluster       # Snapshot a cluster
  asc rds snapshot my-instance my-snapshot --wait         # Snapshot and wait for completion`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(CreateRDSSnapshot(cmd, args))
	},
}

// CreateRDSSnapshot creates a snapshot of an RDS instance or cluster.
func CreateRDSSnapshot(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, rds.NewRDSService)
	if err != nil {
		return fmt.Errorf("create new RDS service: %w", err)
	}

	input := &ascTypes.CreateSnapshotInput{
		Identifier:         args[0],
		SnapshotIdentifier: args[1],
		IsCluster:          snapshotCluster,
	}

	err = svc.CreateSnapshot(cmd.Context(), input)
	if err != nil {
		return fmt.Errorf("create snapshot: %w", err)
	}

	resourceType := "instance"
	if snapshotCluster {
		resourceType = "cluster"
	}
	fmt.Printf("Snapshot %s created for %s %s\n", args[1], resourceType, args[0])

	if snapshotWait {
		fmt.Printf("Waiting for snapshot to become available...\n")
		err = svc.WaitForSnapshot(cmd.Context(), input, 30*time.Minute)
		if err != nil {
			return fmt.Errorf("wait for snapshot: %w", err)
		}
		fmt.Printf("Snapshot %s is now available\n", args[1])
	}

	return nil
}
