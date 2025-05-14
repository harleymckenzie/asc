// rm.go defines the 'rm' subcommand for schedule operations.
package schedule

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/asg"
	ascTypes "github.com/harleymckenzie/asc/internal/service/asg/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// rmCmd defines the 'rm' subcommand for schedule operations.
var rmCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Remove a scheduled action from an Auto Scaling Group",
	Example: "asc asg rm schedule my-schedule --asg-name my-asg",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(RemoveSchedule(cobraCmd, args))
	},
}

// NewRmFlags adds flags for the rm subcommand.
func NewRmFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().StringVarP(&asgName, "asg-name", "a", "", "The name of the Auto Scaling Group")
	cobraCmd.MarkFlagRequired("asg-name")
}

func init() {
	NewRmFlags(rmCmd)
}

//
// Command functions
//

// RemoveSchedule is the handler for the rm schedule subcommand.
func RemoveSchedule(cobraCmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cobraCmd)

	svc, err := asg.NewAutoScalingService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new Auto Scaling Group service: %w", err)
	}

	err = svc.RemoveAutoScalingGroupSchedule(ctx, &ascTypes.RemoveAutoScalingGroupScheduleInput{
		AutoScalingGroupName: asgName,
		ScheduledActionName:  args[0],
	})
	if err != nil {
		return fmt.Errorf("remove schedule: %w", err)
	}
	return nil
}
