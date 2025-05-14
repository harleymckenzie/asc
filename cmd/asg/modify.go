// The modify command allows updating min, max, or desired capacity for an Auto Scaling Group.

package asg

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/asg"
	ascTypes "github.com/harleymckenzie/asc/internal/service/asg/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

// Variables
var (
	minSizeStr         string
	maxSizeStr         string
	desiredCapacityStr string
)

// Init function
func init() {
	addModifyFlags(modifyCmd)
}

// Command variable
var modifyCmd = &cobra.Command{
	Use:     "modify",
	Short:   "Modify an Auto Scaling Group min, max, or desired capacity",
	Long:    "Modify an Auto Scaling Group min, max, or desired capacity",
	Args:    cobra.ExactArgs(1),
	GroupID: "actions",
	Aliases: []string{"edit", "update"},
	Example: "  asc asg modify my-asg --min 3         # Set the minimum capacity to 3\n" +
		"  asc asg modify my-asg --max -6        # Decrease the maximum capacity by 6\n" +
		"  asc asg modify my-asg --desired +5    # Increase the desired capacity by 5",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ModifyAutoScalingGroup(cmd, args))
	},
}

// Flag function
func addModifyFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().SortFlags = false
	cobraCmd.Flags().
		StringVarP(&minSizeStr, "min", "m", "", "The minimum capacity (absolute or relative, e.g. 3, +1, -2)")
	cobraCmd.Flags().
		StringVarP(&maxSizeStr, "max", "M", "", "The maximum capacity (absolute or relative, e.g. 3, +3, -3)")
	cobraCmd.Flags().
		StringVarP(&desiredCapacityStr, "desired", "d", "", "The desired capacity (absolute or relative, e.g. 3, +1, -2)")
}

// Command functions
func ModifyAutoScalingGroup(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	profile, region := cmdutil.GetPersistentFlags(cmd)

	svc, err := asg.NewAutoScalingService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new Auto Scaling Service: %w", err)
	}

	// Get current information about the Auto Scaling Group
	getInput := &ascTypes.GetAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{args[0]},
	}
	asgOutput, err := svc.GetAutoScalingGroups(ctx, getInput)
	if err != nil {
		return fmt.Errorf("get Auto Scaling Groups: %w", err)
	}

	// Create a ModifyAutoScalingGroupInput struct to be updated with the new information
	input := &ascTypes.ModifyAutoScalingGroupInput{
		AutoScalingGroupName: args[0],
	}

	// Apply the relative or absolute values to the ModifyAutoScalingGroupInput struct
	if minSizeStr != "" {
		minSizeInt32, err := utils.ApplyRelativeOrAbsolute(minSizeStr, *asgOutput[0].MinSize)
		if err != nil {
			return fmt.Errorf("apply relative or absolute min size: %w", err)
		}
		input.MinSize = &minSizeInt32
	}
	if maxSizeStr != "" {
		maxSizeInt32, err := utils.ApplyRelativeOrAbsolute(maxSizeStr, *asgOutput[0].MaxSize)
		if err != nil {
			return fmt.Errorf("apply relative or absolute max size: %w", err)
		}
		input.MaxSize = &maxSizeInt32
	}
	if desiredCapacityStr != "" {
		desiredCapacityInt32, err := utils.ApplyRelativeOrAbsolute(
			desiredCapacityStr,
			*asgOutput[0].DesiredCapacity,
		)
		if err != nil {
			return fmt.Errorf("apply relative or absolute desired capacity: %w", err)
		}
		input.DesiredCapacity = &desiredCapacityInt32
	}

	// Modify the Auto Scaling Group
	err = svc.ModifyAutoScalingGroup(ctx, input)
	if err != nil {
		return fmt.Errorf("modify Auto Scaling Group: %w", err)
	}
	return nil
}
