package ec2

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/ec2"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

type Column struct {
	ID      string
	Visible bool
	Sort    bool
}

var (
	list   bool
	sortBy string

	showAMI        bool
	showLaunchTime bool
	showPrivateIP  bool

	sortName       bool
	sortID         bool
	sortType       bool
	sortLaunchTime bool
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all EC2 instances",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()
		profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
		region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

		svc, err := ec2.NewEC2Service(ctx, profile, region)
		if err != nil {
			log.Fatalf("Failed to initialize EC2 service: %v", err)
		}

		instances, err := svc.GetInstances(ctx)
		if err != nil {
			log.Fatalf("Failed to list EC2 instances: %v", err)
		}

		// Define available columns and associated flags
		columns := []Column{
			{ID: "Name", Visible: true, Sort: sortName},
			{ID: "Instance ID", Visible: true, Sort: sortID},
			{ID: "State", Visible: true, Sort: false},
			{ID: "Instance Type", Visible: true, Sort: sortType},
			{ID: "Public IP", Visible: true, Sort: false},
			{ID: "AMI ID", Visible: showAMI, Sort: false},
			{ID: "Launch Time", Visible: showLaunchTime, Sort: sortLaunchTime},
			{ID: "Private IP", Visible: showPrivateIP, Sort: false},
		}

		selectedColumns := make([]string, 0, len(columns))

		// Dynamically build the list of columns
		for _, col := range columns {
			if col.Visible {
				selectedColumns = append(selectedColumns, col.ID)
			}
			if col.Sort {
				sortBy = col.ID
			}
		}

		if sortBy == "" {
			sortBy = "Name"
		}

		tableformat.Render(&ec2.EC2Table{
			Instances:       instances,
			SelectedColumns: selectedColumns,
		}, sortBy)
	},
}

func init() {
	// Add flags - Output
	lsCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs EC2 instances in list format.")
	lsCmd.Flags().BoolVarP(&showAMI, "ami", "A", false, "Show the AMI ID of the instance.")
	lsCmd.Flags().BoolVarP(&showLaunchTime, "launch-time", "L", false, "Show the launch time of the instance.")
	lsCmd.Flags().BoolVarP(&showPrivateIP, "private-ip", "P", false, "Show the private IP address of the instance.")

	// Add flags - Sorting
	lsCmd.Flags().BoolVarP(&sortName, "sort-name", "n", true, "Sort by descending EC2 instance name.")
	lsCmd.Flags().BoolVarP(&sortID, "sort-id", "i", false, "Sort by descending EC2 instance Id.")
	lsCmd.Flags().BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending EC2 instance type.")
	lsCmd.Flags().BoolVarP(&sortLaunchTime, "sort-launch-time", "t", false, "Sort by descending launch time (most recently launched first).")
	lsCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-id", "sort-type", "sort-launch-time")

	lsCmd.Flags().SortFlags = false
}
