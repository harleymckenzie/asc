package asg

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/asg"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

var lsSchedulesCmd = &cobra.Command{
	Use:   "schedules",
	Short: "List all schedules for an Auto Scaling Group",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()
		profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
		region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

		svc, err := asg.NewAutoScalingService(ctx, profile, region)
		if err != nil {
			log.Fatalf("Failed to initialize Auto Scaling Group service: %v", err)
		}

        // If an argument is provided, use it as the Auto Scaling Group name,
        // otherwise dont provide a name to GetSchedules
        var asgName string
        if len(args) > 0 {
            asgName = args[0]
        }
        schedules, err := svc.GetSchedules(ctx, &asg.GetSchedulesInput{
            AutoScalingGroupName: asgName,
        })
        if err != nil {
            log.Fatalf("Failed to get schedules for Auto Scaling Group %s: %v", asgName, err)
        }

        // Define columns for schedules
        columns := []tableformat.Column{
            {ID: "Name", Visible: true, Sort: sortName},
            {ID: "Recurrence", Visible: true, Sort: false},
            {ID: "Start Time", Visible: true, Sort: false},
            {ID: "End Time", Visible: true, Sort: false},
            {ID: "Desired Capacity", Visible: true, Sort: false},
            {ID: "Min", Visible: true, Sort: false},
            {ID: "Max", Visible: true, Sort: false},
        }
        selectedColumns, sortBy := tableformat.BuildColumns(columns)

        tableformat.Render(&asg.AutoScalingSchedulesTable{
            Schedules: schedules,
            SelectedColumns: selectedColumns,
        }, sortBy, list)
	},
}

func init() {
	lsSchedulesCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Auto-Scaling Groups in list format.")
    lsSchedulesCmd.Flags().SortFlags = false
}
