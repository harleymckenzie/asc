package role

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/iam"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	reverseSort bool
	showARNs    bool
)

func init() {
	NewLsFlags(lsCmd)
}

func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Role Name", Category: "Role Details", Visible: true, DefaultSort: true},
		{Name: "ARN", Category: "Role Details", Visible: showARNs},
		{Name: "Path", Category: "Role Details", Visible: false},
		{Name: "Created Date", Category: "Role Details", Visible: true},
		{Name: "Max Session Duration", Category: "Role Details", Visible: false},
		{Name: "Description", Category: "Role Details", Visible: false},
	}
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List IAM roles",
	Aliases: []string{"list"},
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListRoles(cmd, args))
	},
}

func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Output in list format.")
	cobraCmd.Flags().BoolVarP(&showARNs, "arn", "a", false, "Show ARNs for each role.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cmdutil.AddTagFlag(cobraCmd)
}

func ListRoles(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, iam.NewIAMService)
	if err != nil {
		return fmt.Errorf("create iam service: %w", err)
	}

	roles, err := svc.GetRoles(cmd.Context())
	if err != nil {
		return fmt.Errorf("get roles: %w", err)
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         "Roles",
		PlainStyle:    list,
		Fields:        getListFields(),
		Tags:          cmdutil.Tags,
		Data:          utils.SlicesToAny(roles),
		GetFieldValue: iam.GetFieldValue,
		GetTagValue:   iam.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}
