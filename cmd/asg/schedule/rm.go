// rm.go defines the 'rm' subcommand for schedule operations.
package schedule

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/asg"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/asg/types"
	"github.com/spf13/cobra"
)

// rmCmd defines the 'rm' subcommand for schedule operations.
var rmCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Remove a scheduled action from an Auto Scaling Group",
	Example: "asc asg rm schedule my-schedule --asg-name my-asg",
	Run: func(cobraCmd *cobra.Command, args []string) {
		RemoveSchedule(cobraCmd, args)
	},
}

// RemoveSchedule is the handler for the rm schedule subcommand.
func RemoveSchedule(cobraCmd *cobra.Command, args []string) {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	svc, err := asg.NewAutoScalingService(ctx, profile, region)
	if err != nil {
		log.Fatalf("Failed to initialize Auto Scaling Group service: %v", err)
	}

	err = svc.RemoveAutoScalingGroupSchedule(ctx, &ascTypes.RemoveAutoScalingGroupScheduleInput{
		AutoScalingGroupName: asgName,
		ScheduledActionName:  args[0],
	})
	if err != nil {
		log.Fatalf("Failed to remove schedule: %v", err)
	}
}

// addRmFlags adds flags for the rm subcommand.
func addRmFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().StringVarP(&asgName, "asg-name", "a", "", "The name of the Auto Scaling Group")
	cobraCmd.MarkFlagRequired("asg-name")
}

func init() {
	addRmFlags(rmCmd)
}
