package target_group

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/elb"
	ascTypes "github.com/harleymckenzie/asc/internal/service/elb/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
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

func getTargetGroupFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Target Group Details", Visible: true},
		{Name: "ARN", Category: "Target Group Details", Visible: showARNs},
		{Name: "Port", Category: "Network", Visible: true},
		{Name: "Protocol", Category: "Network", Visible: true},
		{Name: "Target Type", Category: "Target Group Details", Visible: true},
		{Name: "Load Balancer", Category: "Load Balancer", Visible: true},
		{Name: "VPC ID", Category: "Network", Visible: true},
		{Name: "Health Check Enabled", Category: "Health Check", Visible: showHealthCheckEnabled},
		{Name: "Health Check Path", Category: "Health Check", Visible: showHealthCheckPath},
		{Name: "Health Check Port", Category: "Health Check", Visible: showHealthCheckPort},
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
func ListTargetGroups(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, elb.NewELBService)
	if err != nil {
		return fmt.Errorf("create elb service: %w", err)
	}

	input := &ascTypes.ListTargetGroupsInput{}
	if len(args) > 0 {
		input.Names = []string{args[0]}
	}

	targetGroups, err := svc.GetTargetGroups(cmd.Context(), &ascTypes.GetTargetGroupsInput{
		ListTargetGroupsInput: *input,
	})
	if err != nil {
		return fmt.Errorf("get target groups: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Target Groups",
	})
	if list {
		table.SetStyle("plain")
	}
	fields := getTargetGroupFields()

	headerRow := cmdutil.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(cmdutil.BuildRows(utils.SlicesToAny(targetGroups), fields, elb.GetFieldValue, elb.GetTagValue))
	table.SortBy(fields, reverseSort)

	table.Render()
	return nil
}
