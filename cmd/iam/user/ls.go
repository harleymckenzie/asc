package user

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
		{Name: "User Name", Category: "User Details", Visible: true, DefaultSort: true},
		{Name: "ARN", Category: "User Details", Visible: showARNs},
		{Name: "Path", Category: "User Details", Visible: false},
		{Name: "Created Date", Category: "User Details", Visible: true},
		{Name: "Password Last Used", Category: "User Details", Visible: true},
	}
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List IAM users",
	Aliases: []string{"list"},
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListUsers(cmd, args))
	},
}

func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Output in list format.")
	cobraCmd.Flags().BoolVarP(&showARNs, "arn", "a", false, "Show ARNs for each user.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cmdutil.AddTagFlag(cobraCmd)
}

func ListUsers(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, iam.NewIAMService)
	if err != nil {
		return fmt.Errorf("create iam service: %w", err)
	}

	users, err := svc.GetUsers(cmd.Context())
	if err != nil {
		return fmt.Errorf("get users: %w", err)
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         "Users",
		PlainStyle:    list,
		Fields:        getListFields(),
		Tags:          cmdutil.Tags,
		Data:          utils.SlicesToAny(users),
		GetFieldValue: iam.GetFieldValue,
		GetTagValue:   iam.GetTagValue,
		ReverseSort:   reverseSort,
		HideEmpty:     true,
	})
	return nil
}
