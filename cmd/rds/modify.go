// The modify command allows updating the configuration of an RDS instance.

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
var (
	applyImmediately bool
	dbInstanceClass  string
	preferredMaintenanceWindow string
)

// Init function
func init() {
	addModifyFlags(modifyCmd)
}

// Flag function
func addModifyFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().SortFlags = false
	cobraCmd.Flags().BoolVar(&applyImmediately, "apply-immediately", false, "Whether to apply the changes immediately")
	cobraCmd.Flags().StringVarP(&dbInstanceClass, "type", "T", "", "The new instance type")
	cobraCmd.Flags().StringVarP(&preferredMaintenanceWindow, "maintenance-window", "m", "", "The preferred maintenance window in the format ddd:hh24:mi-ddd:hh24:mi")
}

// Command variable
var modifyCmd = &cobra.Command{
	Use:     "modify",
	Short:   "Modify an RDS instance class or preferred maintenance window",
	Long:    "Modify an RDS instance class or preferred maintenance window",
	Args:    cobra.ExactArgs(1),
	GroupID: "actions",
	Aliases: []string{"edit", "update"},
	Example: "  asc rds modify my-instance --apply-immediately --type t3.micro --maintenance-window 'mon:00:00-mon:03:00'  # Modify the instance to a t3.micro instance and apply the changes immediately and set the preferred maintenance window to Monday 00:00-03:00",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ModifyRDSInstance(cmd, args))
	},
}

// Command functions
func ModifyRDSInstance(cmd *cobra.Command, args []string) error {
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

	// Create a ModifyInstanceInput struct to be updated with the new information
	input := &ascTypes.ModifyInstanceInput{
		DBInstanceIdentifier: &args[0],
	}

	// Apply the relative or absolute values to the ModifyInstanceInput struct
	if applyImmediately {
		input.ApplyImmediately = &applyImmediately
	}
	if dbInstanceClass != "" {
		input.DBInstanceClass = &dbInstanceClass
	}
	if preferredMaintenanceWindow != "" {
		input.PreferredMaintenanceWindow = &preferredMaintenanceWindow
	}

	// Modify the RDS instance
	err = svc.ModifyInstance(ctx, input)
	if err != nil {
		return fmt.Errorf("modify RDS instance: %w", err)
	}

	return nil
}
