package nacl

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
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

// naclListFields returns the fields for the NACL list table.
func naclListFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Network ACL ID", Display: true, DefaultSort: true},
		{ID: "Associated with", Display: true},
		{ID: "Default", Display: true},
		{ID: "VPC ID", Display: true},
		{ID: "Inbound Rules", Display: true},
		{ID: "Outbound Rules", Display: true},
		{ID: "Owner", Display: true},
	}
}

func naclRulesFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Rule number", Display: true},
		{ID: "Type", Display: true},
		{ID: "Protocol", Display: true},
		{ID: "Port range", Display: true},
		{ID: "Source", Display: true},
		{ID: "Destination", Display: true},
		{ID: "Allow/Deny", Display: true},
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
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	if len(args) > 0 {
		return ListNACLRules(cmd, args)
	} else {
		nacls, err := svc.GetNACLs(ctx, &ascTypes.GetNACLsInput{})
		if err != nil {
			return fmt.Errorf("get network acls: %w", err)
		}

		fields := naclListFields()
		opts := tableformat.RenderOptions{
			Title:  "Network ACLs",
			Style:  "rounded",
			SortBy: tableformat.GetSortByField(fields, reverseSort),
		}

		if list {
			opts.Style = "list"
		}

		tableformat.RenderTableList(&tableformat.ListTable{
			Instances: utils.SlicesToAny(nacls),
			Fields:    fields,
			GetAttribute: func(fieldID string, instance any) (string, error) {
				return vpc.GetNetworkAclAttributeValue(fieldID, instance)
			},
		}, opts)
		return nil
		}
}

// ListNACLRules is the handler for the ls subcommand.
func ListNACLRules(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}
	
	nacl, err := svc.GetNACLs(ctx, &ascTypes.GetNACLsInput{
		NetworkAclIds: []string{args[0]},
	})
	if err != nil {
		return fmt.Errorf("get network acl rules: %w", err)
	}

	rules := nacl[0].Entries
	ingressRules := svc.FilterNACLRules(rules, true)
	egressRules := svc.FilterNACLRules(rules, false)

	fields := naclRulesFields()
	ingressOpts := tableformat.RenderOptions{
		Title:  fmt.Sprintf("%s - Inbound Rules", args[0]),
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}
	egressOpts := tableformat.RenderOptions{
		Title:  fmt.Sprintf("%s - Outbound Rules", args[0]),
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	// Print inbound and outbound rules in separate tables
	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(ingressRules),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return vpc.GetNACLRuleAttributeValue(fieldID, instance)
		},
	}, ingressOpts)
	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(egressRules),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return vpc.GetNACLRuleAttributeValue(fieldID, instance)
		},
	}, egressOpts)
	return nil
}
