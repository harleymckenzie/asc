// The cancel-pending-modifications command allows cancelling pending modifications for an RDS instance.

package rds

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/rds"
	ascTypes "github.com/harleymckenzie/asc/internal/service/rds/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {
	newCancelPendingModificationsFlags(cancelPendingModificationsCmd)
}

// Flag function
func newCancelPendingModificationsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().SortFlags = false
}

// Command variable
var cancelPendingModificationsCmd = &cobra.Command{
	Use:     "cancel-pending-modifications",
	Short:   "Cancel pending modifications for an RDS instance",
	Long:    "Cancel pending modifications for an RDS instance",
	Args:    cobra.ExactArgs(1),
	GroupID: "actions",
	Aliases: []string{"cancel-pending", "cancel-pending-changes"},
	Example: "  asc rds cancel-pending-modifications my-instance  # Cancel pending modifications for the instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(CancelPendingModifications(cmd, args))
	},
}

// Command functions
func CancelPendingModifications(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	profile, region := cmdutil.GetPersistentFlags(cmd)

	svc, err := rds.NewRDSService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new RDS Service: %w", err)
	}

	// Get current information about the RDS instance
	getInput := &ascTypes.GetInstancesInput{
		InstanceIdentifier: args[0],
	}
	instance, err := svc.GetInstances(ctx, getInput)
	if err != nil {
		return fmt.Errorf("get instance: %w", err)
	}

	// Check if the RDS instance exists
	if len(instance) == 0 {
		return fmt.Errorf("RDS instance not found: %s", args[0])
	}

	// Check if the instance has pending modifications
	if instance[0].PendingModifiedValues == nil {
		return fmt.Errorf("no pending modifications found")
	}

	// Get existing instance class
	instanceClass := instance[0].DBInstanceClass
	if instanceClass == nil {
		return fmt.Errorf("instance class not found")
	}

	// Cancel pending modifications
	err = svc.ModifyInstance(ctx, &ascTypes.ModifyInstanceInput{
		DBInstanceIdentifier: &args[0],
		DBInstanceClass: instanceClass,
		ApplyImmediately: &[]bool{false}[0],
	})
	if err != nil {
		return fmt.Errorf("cancel pending modifications: %w", err)
	}

	return nil
}