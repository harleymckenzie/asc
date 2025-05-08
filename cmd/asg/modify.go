// The modify command allows updating min, max, or desired capacity for an Auto Scaling Group.

package asg

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/asg"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/asg/types"
	"github.com/harleymckenzie/asc/pkg/shared/utils"
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
	Example: "modify my-asg --min 3         # Set the minimum capacity to 3\n" +
		"modify my-asg --max -6        # Decrease the maximum capacity by 6\n" +
		"modify my-asg --desired +5    # Increase the desired capacity by 5",
	Run: func(cmd *cobra.Command, args []string) {
		ModifyAutoScalingGroup(cmd, args)
	},
}

// Flag function
func addModifyFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().StringVarP(&minSizeStr, "min", "m", "", "The minimum capacity (absolute or relative, e.g. 3, +1, -2)")
	cobraCmd.Flags().StringVarP(&maxSizeStr, "max", "M", "", "The maximum capacity (absolute or relative, e.g. 3, +3, -3)")
	cobraCmd.Flags().StringVarP(&desiredCapacityStr, "desired", "d", "", "The desired capacity (absolute or relative, e.g. 3, +1, -2)")
}

// Command functions
func ModifyAutoScalingGroup(cobraCmd *cobra.Command, args []string) {
	ctx := context.Background()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	svc, err := asg.NewAutoScalingService(ctx, profile, region)
	if err != nil {
		log.Fatalf("Error creating Auto Scaling Service: %v", err)
	}

	// Get current information about the Auto Scaling Group
	getInput := &ascTypes.GetAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{args[0]},
	}
	asgOutput, err := svc.GetAutoScalingGroups(ctx, getInput)
	if err != nil {
		log.Fatalf("Error getting Auto Scaling Groups: %v", err)
	}

	// Create a ModifyAutoScalingGroupInput struct to be updated with the new information
	input := &ascTypes.ModifyAutoScalingGroupInput{
		AutoScalingGroupName: args[0],
	}

	// Apply the relative or absolute values to the ModifyAutoScalingGroupInput struct
	if minSizeStr != "" {
		minSizeInt32, err := utils.ApplyRelativeOrAbsolute(minSizeStr, *asgOutput[0].MinSize)
		if err != nil {
			log.Fatalf("Error applying relative or absolute value: %v", err)
		}
		input.MinSize = &minSizeInt32
	}
	if maxSizeStr != "" {
		maxSizeInt32, err := utils.ApplyRelativeOrAbsolute(maxSizeStr, *asgOutput[0].MaxSize)
		if err != nil {
			log.Fatalf("Error applying relative or absolute value: %v", err)
		}
		input.MaxSize = &maxSizeInt32
	}
	if desiredCapacityStr != "" {
		desiredCapacityInt32, err := utils.ApplyRelativeOrAbsolute(desiredCapacityStr, *asgOutput[0].DesiredCapacity)
		if err != nil {
			log.Fatalf("Error applying relative or absolute value: %v", err)
		}
		input.DesiredCapacity = &desiredCapacityInt32
	}

	// Modify the Auto Scaling Group
	err = svc.ModifyAutoScalingGroup(ctx, input)
	if err != nil {
		log.Fatalf("Error modifying Auto Scaling Group: %v", err)
	}
}
