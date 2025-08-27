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
)

func init() {
	NewLsFlags(lsCmd)
}

// getListFields returns the fields for the NACL list table.
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Network ACL ID", Category: "NACL", Visible: true, SortBy: true, SortDirection: tablewriter.Asc},
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

		table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
			Title: "Network ACLs",
		})
		if list {
			table.SetRenderStyle("plain")
		}

		fields := getListFields()
		fields = cmdutil.AppendTagFields(fields, cmdutil.Tags, utils.SlicesToAny(nacls))

		headerRow := cmdutil.BuildHeaderRow(fields)
		table.AppendHeader(headerRow)
		table.AppendRows(cmdutil.BuildRows(utils.SlicesToAny(nacls), fields, vpc.GetFieldValue, vpc.GetTagValue))
		table.SortBy(fields, reverseSort)
		table.Render()
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
	ingressTable := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: fmt.Sprintf("%s - Inbound Rules", args[0]),
	})
	if list {
		ingressTable.SetRenderStyle("plain")
	}

	headerRow := cmdutil.BuildHeaderRow(fields)
	ingressTable.AppendHeader(headerRow)
	ingressTable.AppendRows(cmdutil.BuildRows(utils.SlicesToAny(ingressRules), fields, vpc.GetFieldValue, vpc.GetTagValue))
	ingressTable.SortBy(fields, reverseSort)
	ingressTable.Render()

	// Print outbound rules table
	egressTable := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: fmt.Sprintf("%s - Outbound Rules", args[0]),
	})
	if list {
		egressTable.SetRenderStyle("plain")
	}

	egressTable.AppendHeader(headerRow)
	egressTable.AppendRows(cmdutil.BuildRows(utils.SlicesToAny(egressRules), fields, vpc.GetFieldValue, vpc.GetTagValue))
	egressTable.SortBy(fields, reverseSort)
	egressTable.Render()
	return nil
}
