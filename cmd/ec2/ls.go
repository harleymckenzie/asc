package ec2

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/ec2"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
)

type ListInstancesInput struct {
	GetInstancesInput *ascTypes.GetInstancesInput
	SelectedColumns   []string
	TableOpts         tableformat.RenderOptions
}

var (
	list           bool
	showAMI        bool
	showLaunchTime bool
	showPrivateIP  bool

	sortName       bool
	sortID         bool
	sortType       bool
	sortLaunchTime bool
)

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all EC2 instances",
	Run: func(cobraCmd *cobra.Command, args []string) {
		// Define available columns and associated flags
		columns := []tableformat.Column{
			{ID: "Name", Visible: true, Sort: sortName},
			{ID: "Instance ID", Visible: true, Sort: sortID},
			{ID: "State", Visible: true, Sort: false},
			{ID: "Instance Type", Visible: true, Sort: sortType},
			{ID: "Public IP", Visible: true, Sort: false},
			{ID: "AMI ID", Visible: showAMI, Sort: false},
			{ID: "Launch Time", Visible: showLaunchTime, Sort: sortLaunchTime},
			{ID: "Private IP", Visible: showPrivateIP, Sort: false},
		}
		selectedColumns, sortBy := tableformat.BuildColumns(columns)

		opts := tableformat.RenderOptions{
			SortBy: sortBy,
			List:   list,
			Title:  "EC2 Instances",
		}

		ListEC2Instances(cobraCmd, ListInstancesInput{
			GetInstancesInput: &ascTypes.GetInstancesInput{},
			SelectedColumns:   selectedColumns,
			TableOpts:         opts,
		})
	},
}

func ListEC2Instances(cobraCmd *cobra.Command, input ListInstancesInput) {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		log.Fatalf("Failed to initialize EC2 service: %v", err)
	}

	instances, err := svc.GetInstances(ctx, input.GetInstancesInput)
	if err != nil {
		log.Fatalf("Failed to list EC2 instances: %v", err)
	}

	tableformat.Render(&ec2.EC2Table{
		Instances:       instances,
		SelectedColumns: input.SelectedColumns,
	}, input.TableOpts)
}

func addLsFlags(lsCmd *cobra.Command) {
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

func init() {
	addLsFlags(lsCmd)
}
