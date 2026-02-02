package nacl

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	reverseSort bool
	sortId      bool
)

func init() {
	NewLsFlags(lsCmd)
}

// getListFields returns the fields for the NACL list table.
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Network ACL ID", Category: "NACL", Visible: true, DefaultSort: true, SortBy: sortId, SortDirection: tablewriter.Asc},
		{Name: "Associated with", Category: "NACL", Visible: true},
		{Name: "Default", Category: "NACL", Visible: true},
		{Name: "VPC ID", Category: "NACL", Visible: true},
		{Name: "Inbound Rules", Category: "NACL", Visible: true},
		{Name: "Outbound Rules", Category: "NACL", Visible: true},
		{Name: "Owner", Category: "NACL", Visible: true},
	}
}

func getRulesFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Rule number", Category: "Rule", Visible: true},
		{Name: "Type", Category: "Rule", Visible: true},
		{Name: "Protocol", Category: "Rule", Visible: true},
		{Name: "Port Range", Category: "Rule", Visible: true},
		{Name: "Source", Category: "Rule", Visible: true},
		{Name: "Allow/Deny", Category: "Rule", Visible: true},
	}
}

// lsCmd is the cobra command for listing NACLs.
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all Network ACLs",
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListNACLs(cmd, args))
	},
}

// NewLsFlags adds flags for the ls subcommand.
func NewLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs NACLs in list format.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Reverse the sort order")
	cobraCmd.Flags().BoolVarP(&sortId, "sort-id", "i", false, "Sort by descending network ACL ID.")
}

// ListNACLs is the handler for the ls subcommand.
func ListNACLs(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	if len(args) > 0 {
		return ListNACLRules(cmd, args)
	} else {
		nacls, err := svc.GetNACLs(cmd.Context(), &ascTypes.GetNACLsInput{})
		if err != nil {
			return fmt.Errorf("get network acls: %w", err)
		}

		tablewriter.RenderList(tablewriter.RenderListOptions{
			Title:         "Network ACLs",
			PlainStyle:    list,
			Fields:        getListFields(),
			Tags:          cmdutil.Tags,
			Data:          utils.SlicesToAny(nacls),
			GetFieldValue: vpc.GetFieldValue,
			GetTagValue:   vpc.GetTagValue,
			ReverseSort:   reverseSort,
		})
		return nil
	}
}

// ListNACLRules is the handler for the ls subcommand.
func ListNACLRules(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, vpc.NewVPCService)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	nacl, err := svc.GetNACLs(cmd.Context(), &ascTypes.GetNACLsInput{
		NetworkAclIds: []string{args[0]},
	})
	if err != nil {
		return fmt.Errorf("get network acl rules: %w", err)
	}

	rules := nacl[0].Entries
	ingressRules := svc.FilterNACLRules(rules, true)
	egressRules := svc.FilterNACLRules(rules, false)

	fields := getRulesFields()

	// Print inbound rules table
	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         fmt.Sprintf("%s - Inbound Rules", args[0]),
		PlainStyle:    list,
		Fields:        fields,
		Data:          utils.SlicesToAny(ingressRules),
		GetFieldValue: vpc.GetFieldValue,
		GetTagValue:   vpc.GetTagValue,
		ReverseSort:   reverseSort,
	})

	// Print outbound rules table
	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         fmt.Sprintf("%s - Outbound Rules", args[0]),
		PlainStyle:    list,
		Fields:        fields,
		Data:          utils.SlicesToAny(egressRules),
		GetFieldValue: vpc.GetFieldValue,
		GetTagValue:   vpc.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}
