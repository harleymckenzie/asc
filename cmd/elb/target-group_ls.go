package elb

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/elb"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/elb/types"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

var targetGroupLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List target groups for the specified ELB",
	Run: func(cobraCmd *cobra.Command, args []string) {
        lsTargetGroups(cobraCmd, args)
    },
}

func lsTargetGroups(cobraCmd *cobra.Command, args []string) {
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

    ListTargetGroups(svc, input)
}

// ListELBTargetGroups lists all target groups for a given ELB
func ListTargetGroups(svc *elb.ELBService, input *ascTypes.ListTargetGroupsInput) {
    ctx := context.TODO()
    targetGroups, err := svc.GetTargetGroups(ctx, &ascTypes.GetTargetGroupsInput{
        ListTargetGroupsInput: *input,
    })
    if err != nil {
        log.Fatalf("Failed to get target groups: %v", err)
    }

    // Define columns for target groups
    columns := []tableformat.Column{
        {ID: "Name", Visible: true, Sort: false},
        {ID: "Health Check Enabled", Visible: true, Sort: false},
        {ID: "Health Check Path", Visible: true, Sort: false},
        {ID: "Health Check Port", Visible: true, Sort: false},
    }
    selectedColumns, sortBy := tableformat.BuildColumns(columns)

    opts := tableformat.RenderOptions{
        SortBy: sortBy,
        List:   list,
    }

    tableformat.Render(&elb.ELBTargetGroupTable{
        TargetGroups: targetGroups,
        SelectedColumns: selectedColumns,
    }, opts)
}
