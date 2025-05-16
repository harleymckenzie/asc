// ls.go defines the 'ls' subcommand for AMI operations.
package ami

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	sortID      bool
	sortName    bool
	sortState   bool
	showDesc    bool
	reverseSort bool

	scope      string // Combined owner/visibility flag
	nameFilter string
	limit      int
)

// Init function
func init() {
	NewLsFlags(lsCmd)
}

// ec2AMIListFields returns the fields for the AMI list table.
func ec2AMIListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "AMI Name", Display: true, Sort: sortName},
		{ID: "AMI ID", Display: true, Sort: sortID},
		{ID: "Source", Display: false},
		{ID: "Owner", Display: true},
		{ID: "Visibility", Display: false},
		{ID: "Status", Display: true, Sort: sortState},
		{ID: "Creation Date", Display: true, DefaultSort: true, SortDirection: "desc"},
		{ID: "Platform", Display: false},
		{ID: "Root Device Type", Display: false},
		{ID: "Block Devices", Display: false},
		{ID: "Virtualization", Display: false},
	}
}

// lsCmd is the cobra command for listing AMIs.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List AMIs",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListAMIs(cmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs AMIs in list format.")
	cobraCmd.Flags().BoolVarP(&sortID, "sort-id", "i", false, "Sort by descending image ID.")
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending image name.")
	cobraCmd.Flags().
		BoolVarP(&sortState, "sort-state", "s", false, "Sort by descending image state.")
	cobraCmd.Flags().
		BoolVarP(&showDesc, "show-description", "d", false, "Show the AMI description column.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
	cobraCmd.Flags().
		StringVar(&scope, "scope", "self", "Scope of AMIs to list: self (your private AMIs), private (all private AMIs you can access), public, amazon, all, or AWS account ID.")
	cobraCmd.Flags().StringVar(&nameFilter, "name", "", "Substring to match in AMI name.")
	cobraCmd.Flags().IntVar(&limit, "limit", 0, "Limit the number of AMIs displayed.")
}

// ListAMIs is the handler for the ls subcommand.
func ListAMIs(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	filters := []types.Filter{}
	owners := []string{}

	switch scope {
	case "self":
		owners = append(owners, "self")
		filters = append(filters, types.Filter{
			Name:   aws.String("is-public"),
			Values: []string{"false"},
		})
	case "private":
		filters = append(filters, types.Filter{
			Name:   aws.String("is-public"),
			Values: []string{"false"},
		})
	case "public":
		filters = append(filters, types.Filter{
			Name:   aws.String("is-public"),
			Values: []string{"true"},
		})
	case "amazon":
		owners = append(owners, "amazon")
	case "all":
		// No owner or visibility filter
	default:
		// If it's a valid AWS account ID (all digits, 12 chars), treat as owner
		if len(scope) == 12 {
			owners = append(owners, scope)
		} else {
			return fmt.Errorf("invalid scope: %s", scope)
		}
	}

	if nameFilter != "" {
		filters = append(filters, types.Filter{
			Name:   aws.String("name"),
			Values: []string{"*" + nameFilter + "*"},
		})
	}

	input := &ascTypes.GetImagesInput{}
	images, err := svc.GetImagesWithFilters(ctx, input, filters, owners)
	if err != nil {
		return fmt.Errorf("get images: %w", err)
	}

	if limit > 0 && len(images) > limit {
		images = images[:limit]
	}

	fields := ec2AMIListFields()
	opts := tableformat.RenderOptions{
		Title:  "Amazon Machine Images (AMIs)",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(images),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return ec2.GetImageAttributeValue(fieldID, instance)
		},
	}, opts)
	return nil
}
