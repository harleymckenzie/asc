package cloudformation

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/cloudformation"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

var list bool
var sortName bool
var sortStatus bool

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all CloudFormation stacks",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()
		profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
		region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

		svc, err := cloudformation.NewCloudFormationService(ctx, profile, region)
		if err != nil {
			log.Fatalf("Failed to initialize CloudFormation service: %v", err)
		}

		stacks, err := svc.GetStacks(ctx)
		if err != nil {
			log.Fatalf("Failed to list CloudFormation stacks: %v", err)
		}

		columns := []tableformat.Column{
			{ID: "Stack Name", Visible: true, Sort: sortName},
			{ID: "Status", Visible: true, Sort: sortStatus},
		}
		selectedColumns, sortBy := tableformat.BuildColumns(columns)

		opts := tableformat.RenderOptions{
			List:  list,
			Title: "CloudFormation Stacks",
		}

		tableformat.Render(&cloudformation.CloudFormationTable{
			Stacks:          stacks,
			SelectedColumns: selectedColumns,
			SortBy:          sortBy,
		}, opts)
	},
}

func addLsFlags(lsCmd *cobra.Command) {
	lsCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs CloudFormation stacks in list format.")
	lsCmd.Flags().BoolVarP(&sortName, "sort-name", "n", true, "Sort by descending CloudFormation stack name.")
	lsCmd.Flags().BoolVarP(&sortStatus, "sort-status", "s", false, "Sort by descending CloudFormation stack status.")
}

func init() {
	addLsFlags(lsCmd)
}
