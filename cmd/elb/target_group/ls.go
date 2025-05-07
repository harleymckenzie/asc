package target_group

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/elb"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/elb/types"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

var (
	list                   bool
	showARNs               bool
	showHealthCheckEnabled bool
	showHealthCheckPath    bool
	showHealthCheckPort    bool
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List target groups",
	GroupID: "actions",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ListTargetGroups(cobraCmd, args)
	},
}

func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs target groups in list format.")
	cobraCmd.Flags().BoolVarP(&showARNs, "arn", "a", false, "Show ARNs for each target group.")
	cobraCmd.Flags().BoolVarP(&showHealthCheckEnabled, "health-check-enabled", "e", false, "Show health check enabled for each target group.")
	cobraCmd.Flags().BoolVarP(&showHealthCheckPath, "health-check-path", "c", false, "Show health check path for each target group.")
	cobraCmd.Flags().BoolVarP(&showHealthCheckPort, "health-check-port", "P", false, "Show health check port for each target group.")
}

func init() {
	NewLsFlags(lsCmd)
}

// ListELBTargetGroups lists all target groups for a given ELB
func ListTargetGroups(cobraCmd *cobra.Command, args []string) {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	svc, err := elb.NewELBService(ctx, profile, region)
	if err != nil {
		log.Fatalf("Failed to initialize ELB service: %v", err)
	}

	input := &ascTypes.ListTargetGroupsInput{}
	if len(args) > 0 {
		input.Names = []string{args[0]}
	}

	targetGroups, err := svc.GetTargetGroups(ctx, &ascTypes.GetTargetGroupsInput{
		ListTargetGroupsInput: *input,
	})
	if err != nil {
		log.Fatalf("Failed to get target groups: %v", err)
	}

	// Define columns for target groups
	columns := []tableformat.Column{
		{ID: "Name", Visible: true, Sort: false},
		{ID: "ARN", Visible: showARNs, Sort: false},
		{ID: "Port", Visible: true, Sort: false},
		{ID: "Protocol", Visible: true, Sort: false},
		{ID: "Target Type", Visible: true, Sort: false},
		{ID: "Load Balancer", Visible: true, Sort: false},
		{ID: "VPC ID", Visible: true, Sort: false},
		{ID: "Health Check Enabled", Visible: false, Sort: showHealthCheckEnabled},
		{ID: "Health Check Path", Visible: false, Sort: showHealthCheckPath},
		{ID: "Health Check Port", Visible: false, Sort: showHealthCheckPort},
	}
	selectedColumns, sortBy := tableformat.BuildColumns(columns)

	opts := tableformat.RenderOptions{
		SortBy: sortBy,
		List:   list,
	}

	tableformat.Render(&elb.ELBTargetGroupTable{
		TargetGroups:    targetGroups,
		SelectedColumns: selectedColumns,
	}, opts)
}
