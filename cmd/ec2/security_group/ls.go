// ls.go defines the 'ls' subcommand for security group operations.
package security_group

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	sortID      bool
	sortName    bool
	sortVPCID   bool
	sortOwnerID bool
	showDesc    bool
	showOwnerID bool
	reverseSort bool
)

func init() {
	NewLsFlags(lsCmd)
}

// lsCmd is the cobra command for listing security groups.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all security groups",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListSecurityGroups(cmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs security groups in list format.")
	cobraCmd.Flags().BoolVarP(&sortID, "sort-id", "i", false, "Sort by descending group ID.")
	cobraCmd.Flags().BoolVarP(&sortName, "sort-name", "n", false, "Sort by descending group name.")
	cobraCmd.Flags().BoolVarP(&showDesc, "show-description", "d", false, "Show the security group description column.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
	cobraCmd.Flags().BoolVarP(&sortVPCID, "sort-vpc-id", "v", false, "Sort by descending VPC ID.")
	cobraCmd.Flags().BoolVarP(&showOwnerID, "show-owner-id", "O", false, "Show the security group owner ID column.")
}

func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Group Name", Category: "Security Group", Visible: true, DefaultSort: true, SortBy: sortName, SortDirection: tablewriter.Asc},
		{Name: "Group ID", Category: "Security Group", Visible: true, SortBy: sortID, SortDirection: tablewriter.Asc},
		{Name: "Description", Category: "Security Group", Visible: showDesc},
		{Name: "VPC ID", Category: "Security Group", Visible: true, SortBy: sortVPCID, SortDirection: tablewriter.Asc},
		{Name: "Owner ID", Category: "Security Group", Visible: showOwnerID, SortBy: sortOwnerID, SortDirection: tablewriter.Asc},
		{Name: "Ingress Count", Category: "Security Group", Visible: true, SortDirection: tablewriter.Desc},
		{Name: "Egress Count", Category: "Security Group", Visible: true, SortDirection: tablewriter.Desc},
		{Name: "Tag Count", Category: "Security Group", Visible: false},
	}
}

// ListSecurityGroups is the handler for the ls subcommand.
// If a security group name is provided, it will list the IP permissions for that security group.
// Otherwise, it will list all security groups.
func ListSecurityGroups(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return ListSecurityGroupRules(cmd, args)
	}

	svc, err := cmdutil.CreateService(cmd, ec2.NewEC2Service)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	groups, err := svc.GetSecurityGroups(cmd.Context(), &ascTypes.GetSecurityGroupsInput{})
	if err != nil {
		return fmt.Errorf("get security groups: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Security Groups",
	})
	if list {
		table.SetRenderStyle("plain")
	}
	fields := getListFields()
	fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(groups))

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(groups), fields, ec2.GetFieldValue, ec2.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)

	table.Render()
	return nil
}

func getListRulesFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Rule ID", Category: "Security Group Rule", Visible: true},
		{Name: "IP Version", Category: "Security Group Rule", Visible: true},
		{Name: "Type", Category: "Security Group Rule", Visible: true},
		{Name: "Protocol", Category: "Security Group Rule", Visible: true},
		{Name: "Port Range", Category: "Security Group Rule", Visible: true},
		{Name: "Source", Category: "Security Group Rule", Visible: true},
		{Name: "Destination", Category: "Security Group Rule", Visible: true},
		{Name: "Description", Category: "Security Group Rule", Visible: true},
	}
}

func ListSecurityGroupRules(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, ec2.NewEC2Service)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	rules, err := svc.GetSecurityGroupRules(cmd.Context(), &ascTypes.GetSecurityGroupRulesInput{
		SecurityGroupID: args[0],
	})
	if err != nil {
		return fmt.Errorf("get security group rules: %w", err)
	}

	ingressRules := ec2.FilterSecurityGroupRules(rules, false)
	egressRules := ec2.FilterSecurityGroupRules(rules, true)

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title:          fmt.Sprintf("%s - Inbound Rules", args[0]),
		MaxColumnWidth: 50,
		Columns:        8,
	})
	if list {
		table.SetRenderStyle("plain")
	}
	fields := getListRulesFields()
	fields = tablewriter.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(ingressRules))

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(ingressRules), fields, ec2.GetFieldValue, ec2.GetTagValue))
	table.AppendTitleRow(fmt.Sprintf("%s - Outbound Rules", args[0]))
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(egressRules), fields, ec2.GetFieldValue, ec2.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)
	table.Render()
	return nil
}
