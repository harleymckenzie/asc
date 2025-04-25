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
	Default bool
	Flag    *bool
}

var (
	list            bool
	sortOrder       []string

	showAMI         bool
	showLaunchTime  bool
	showPrivateIP   bool
)

func NewEC2Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ec2",
		Short: "Perform EC2 operations",
	}

	// ls sub command
	lsCmd := &cobra.Command{
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
				{ID: "name", Default: true},
				{ID: "instance_id", Default: true},
				{ID: "state", Default: true},
				{ID: "instance_type", Default: true},
				{ID: "public_ip", Default: true},
				{ID: "ami_id", Flag: &showAMI},
				{ID: "launch_time", Flag: &showLaunchTime},
				{ID: "private_ip", Flag: &showPrivateIP},
			}

			selectedColumns := make([]string, 0, len(columns))

			// Dynamically build the list of columns
			for _, col := range columns {
				if col.Default || (col.Flag != nil && *col.Flag) {
					selectedColumns = append(selectedColumns, col.ID)
				}
			}

			tableformat.Render(&ec2.EC2Table{
				Instances:       instances,
				SelectedColumns: selectedColumns,
				SortOrder:       sortOrder,
			}, list)
		},
	}
	cmd.AddCommand(lsCmd)

	// Add flags - Output
	lsCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs EC2 instances in list format.")
	lsCmd.Flags().BoolVarP(&showAMI, "ami", "A", false, "Show the AMI ID of the instance.")
	lsCmd.Flags().BoolVarP(&showLaunchTime, "launch-time", "L", false, "Show the launch time of the instance.")
	lsCmd.Flags().BoolVarP(&showPrivateIP, "private-ip", "P", false, "Show the private IP address of the instance.")

	// Add flags - Sorting
	lsCmd.Flags().BoolP("sort-name", "n", true, "Sort by descending EC2 instance name.")
	lsCmd.Flags().BoolP("sort-id", "i", false, "Sort by descending EC2 instance Id.")
	lsCmd.Flags().BoolP("sort-type", "T", false, "Sort by descending EC2 instance type.")
	lsCmd.Flags().BoolP("sort-launch-time", "t", false, "Sort by descending launch time (most recently launched first).")
	lsCmd.Flags().SortFlags = false

	return cmd
}
