package asg

import (
	"context"
	"log"
	"github.com/spf13/cobra"
	"github.com/harleymckenzie/asc/pkg/service/asg"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/asg/types"
)

func rmSchedule(cobraCmd *cobra.Command, args []string) {
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

func addScheduleRmFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().StringVarP(&asgName, "asg-name", "a", "", "The name of the Auto Scaling Group")
	cobraCmd.MarkFlagRequired("asg-name")
}
