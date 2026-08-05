// create.go creates a snapshot of an RDS instance or cluster.

package snapshot

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
	createWait    bool
	createCluster bool
)

// Init function
func init() {
	newCreateFlags(createCmd)
}

// Flag function
func newCreateFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().SortFlags = false
	cobraCmd.Flags().BoolVarP(&createWait, "wait", "w", false, "Wait for the snapshot to complete")
	cobraCmd.Flags().BoolVarP(&createCluster, "cluster", "c", false, "Snapshot a cluster instead of an instance")
}

// Command variable
var createCmd = &cobra.Command{
	Use:     "create <identifier> <snapshot-name>",
	Short:   "Create a snapshot of an RDS instance or cluster",
	Long:    "Create a manual snapshot of an RDS instance or cluster. Use --cluster for cluster snapshots.",
	Aliases: []string{"add"},
	Args:    cobra.ExactArgs(2),
	GroupID: "actions",
	Example: `  asc rds snapshot create my-instance my-snapshot            # Snapshot an instance
  asc rds snapshot create my-cluster my-snapshot --cluster   # Snapshot a cluster
  asc rds snapshot create my-instance my-snapshot --wait     # Snapshot and wait for completion`,
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
		IsCluster:          createCluster,
	}

	err = svc.CreateSnapshot(cmd.Context(), input)
	if err != nil {
		return fmt.Errorf("create snapshot: %w", err)
	}

	resourceType := "instance"
	if createCluster {
		resourceType = "cluster"
	}
	fmt.Printf("Snapshot %s created for %s %s\n", args[1], resourceType, args[0])

	if createWait {
		fmt.Printf("Waiting for snapshot to become available...\n")
		err = svc.WaitForSnapshot(cmd.Context(), input, 30*time.Minute)
		if err != nil {
			return fmt.Errorf("wait for snapshot: %w", err)
		}
		fmt.Printf("Snapshot %s is now available\n", args[1])
	}

	return nil
}
