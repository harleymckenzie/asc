// ls.go defines the 'ls' subcommand for AMI operations.
package ami

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list             bool
	sortID           bool
	sortName         bool
	sortState        bool
	sortCreationDate bool
	showDesc         bool
	reverseSort      bool

	scope      string // Combined owner/visibility flag
	nameFilter string
	limit      int
)

// Init function
func init() {
	NewLsFlags(lsCmd)
}

func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "AMI Name", Category: "AMI Details", Visible: true, SortBy: sortName, SortDirection: tablewriter.Asc},
		{Name: "AMI ID", Category: "AMI Details", Visible: true, SortBy: sortID, SortDirection: tablewriter.Asc},
		{Name: "Source", Category: "AMI Details", Visible: false},
		{Name: "Owner", Category: "AMI Details", Visible: true},
		{Name: "Visibility", Category: "AMI Details", Visible: false},
		{Name: "Status", Category: "AMI Details", Visible: true, SortBy: sortState, SortDirection: tablewriter.Desc},
		{Name: "Creation Date", Category: "AMI Details", DefaultSort: true, Visible: true, SortBy: sortCreationDate, SortDirection: tablewriter.Desc},
		{Name: "Platform", Category: "AMI Details", Visible: false},
		{Name: "Root Device Type", Category: "AMI Details", Visible: false},
		{Name: "Block Devices", Category: "AMI Details", Visible: false},
		{Name: "Virtualization", Category: "AMI Details", Visible: false},
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
	cobraCmd.Flags().BoolVarP(&sortState, "sort-state", "s", false, "Sort by descending image state.")
	cobraCmd.Flags().BoolVarP(&sortCreationDate, "sort-creation-date", "c", false, "Sort by descending image creation date.")
	cobraCmd.Flags().BoolVarP(&showDesc, "show-description", "d", false, "Show the AMI description column.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
	cobraCmd.Flags().StringVar(&scope, "scope", "self", "Scope of AMIs to list: self (your private AMIs), private (all private AMIs you can access), public, amazon, all, or AWS account ID.")
	cobraCmd.Flags().StringVar(&nameFilter, "name", "", "Substring to match in AMI name.")
	cobraCmd.Flags().IntVar(&limit, "limit", 0, "Limit the number of AMIs displayed.")
}

func ListAMIs(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, ec2.NewEC2Service)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	filters, owners, err := parseScope(scope)
	if err != nil {
		return fmt.Errorf("parse scope: %w", err)
	}
	amis, err := getImages(cmd.Context(), svc, &ascTypes.GetImagesInput{
		Filters: filters,
		Owners:  owners,
	})
	if err != nil {
		return fmt.Errorf("get images: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "AMIs",
	})
	if list {
		table.SetRenderStyle("plain")
	}
	fields := getListFields()
	fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(amis))

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(amis), fields, ec2.GetFieldValue, ec2.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)

	table.Render()
	return nil
}

func parseScope(scope string) ([]types.Filter, []string, error) {
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
			return nil, nil, fmt.Errorf("invalid scope: %s", scope)
		}
	}

	if nameFilter != "" {
		filters = append(filters, types.Filter{
			Name:   aws.String("name"),
			Values: []string{"*" + nameFilter + "*"},
		})
	}

	return filters, owners, nil
}
