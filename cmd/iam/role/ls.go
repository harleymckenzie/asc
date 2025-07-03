// ls.go defines the 'ls' subcommand for AMI operations.
package role

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/iam"
	ascTypes "github.com/harleymckenzie/asc/internal/service/iam/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	sortArn     bool
	sortName    bool
	sortPath    bool
	reverseSort bool
	maxItems    int32
	pathPrefix  string
)

// Init function
func init() {
	NewLsFlags(lsCmd)
}

// iamRoleListFields returns the fields for the role list table.
func iamRoleListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Name", Display: true, Sort: sortName},
		{ID: "Last Activity", Display: true},
		{ID: "Arn", Display: false, Sort: sortArn},
		{ID: "Creation Time", Display: true},
		{ID: "Description", Display: false},
		{ID: "Path", Display: false, Sort: sortPath},
	}
}

// lsCmd is the cobra command for listing AMIs.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List roles",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListRoles(cmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs AMIs in list format.")
	cobraCmd.Flags().BoolVarP(&sortArn, "sort-arn", "a", false, "Sort by descending ARN.")
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending role name.")
	cobraCmd.Flags().BoolVarP(&sortPath, "sort-path", "P", false, "Sort by descending path.")
	cobraCmd.Flags().Int32Var(&maxItems, "max-items", 100, "Maximum number of items to return.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
	cobraCmd.Flags().StringVar(&pathPrefix, "path-prefix", "", "Path prefix to filter roles by. If not provided, all roles will be returned.")
}

// ListRoles is the handler for the ls subcommand.
func ListRoles(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := iam.NewIAMService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new IAM service: %w", err)
	}

	roles, err := svc.GetRoles(ctx, &ascTypes.GetRolesInput{
		MaxItems:   maxItems,
		PathPrefix: pathPrefix,
	})
	if err != nil {
		return fmt.Errorf("get roles: %w", err)
	}

	fields := iamRoleListFields()
	opts := tableformat.RenderOptions{
		Title:  "Roles",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(roles),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return iam.GetRoleAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
