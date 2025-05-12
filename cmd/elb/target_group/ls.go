package target_group

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/pkg/service/elb"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/elb/types"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/harleymckenzie/asc/pkg/shared/utils"
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
	"github.com/spf13/cobra"
)

var (
	list                   bool
	showARNs               bool
	showHealthCheckEnabled bool
	showHealthCheckPath    bool
	showHealthCheckPort    bool
	reverseSort            bool
)

func init() {
	NewLsFlags(lsCmd)
}

func targetGroupFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Name", Visible: true, Sort: false, DefaultSort: true},
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
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List target groups",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListTargetGroups(cobraCmd, args))
	},
}

func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs target groups in list format.")
	cobraCmd.Flags().BoolVarP(&showARNs, "arn", "a", false, "Show ARNs for each target group.")
	cobraCmd.Flags().
		BoolVarP(&showHealthCheckEnabled, "health-check-enabled", "e", false, "Show health check enabled for each target group.")
	cobraCmd.Flags().
		BoolVarP(&showHealthCheckPath, "health-check-path", "c", false, "Show health check path for each target group.")
	cobraCmd.Flags().
		BoolVarP(&showHealthCheckPort, "health-check-port", "P", false, "Show health check port for each target group.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
}

// ListELBTargetGroups lists all target groups for a given ELB
func ListTargetGroups(cobraCmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
	region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

	svc, err := elb.NewELBService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new ELB service: %w", err)
	}

	input := &ascTypes.ListTargetGroupsInput{}
	if len(args) > 0 {
		input.Names = []string{args[0]}
	}

	targetGroups, err := svc.GetTargetGroups(ctx, &ascTypes.GetTargetGroupsInput{
		ListTargetGroupsInput: *input,
	})
	if err != nil {
		return fmt.Errorf("get target groups: %w", err)
	}

	fields := targetGroupFields()

	opts := tableformat.RenderOptions{
		Title:  "Target Groups",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(targetGroups),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) string {
			return elb.GetTargetGroupAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
