package policy

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
	showAll     bool
	showARNs    bool
)

func init() {
	NewLsFlags(lsCmd)
}

func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Policy Name", Category: "Policy Details", Visible: true, DefaultSort: true},
		{Name: "ARN", Category: "Policy Details", Visible: showARNs},
		{Name: "Attachment Count", Category: "Policy Details", Visible: true},
		{Name: "Default Version ID", Category: "Policy Details", Visible: false},
		{Name: "Created Date", Category: "Policy Details", Visible: false},
		{Name: "Updated Date", Category: "Policy Details", Visible: false},
	}
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List IAM policies",
	Long:    "List IAM policies. Shows customer-managed policies by default. Use --all to include AWS-managed policies.",
	Aliases: []string{"list"},
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListPolicies(cmd, args))
	},
}

func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Output in list format.")
	cobraCmd.Flags().BoolVar(&showARNs, "arn", false, "Show ARNs for each policy.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cobraCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Include AWS-managed policies.")
	cmdutil.AddTagFlag(cobraCmd)
}

func ListPolicies(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, iam.NewIAMService)
	if err != nil {
		return fmt.Errorf("create iam service: %w", err)
	}

	policies, err := svc.GetPolicies(cmd.Context(), showAll)
	if err != nil {
		return fmt.Errorf("get policies: %w", err)
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         "Policies",
		PlainStyle:    list,
		Fields:        getListFields(),
		Tags:          cmdutil.Tags,
		Data:          utils.SlicesToAny(policies),
		GetFieldValue: iam.GetFieldValue,
		GetTagValue:   iam.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}
