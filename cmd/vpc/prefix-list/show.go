package prefix_list

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

// prefixListShowFields returns the fields for the Prefix List detail table.
func prefixListShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Prefix List ID", Display: true},
		{ID: "Prefix List ARN", Display: true},
		{ID: "Name", Display: true},
		{ID: "State", Display: true},
		{ID: "Version", Display: true},
		{ID: "Max entries", Display: true},
		{ID: "Owner", Display: true},
		{ID: "CIDRs", Display: true},
		{ID: "Entries", Display: true},
	}
}

// showCmd is the cobra command for showing Prefix List details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a Prefix List",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowPrefixList(cmd, args[0]))
	},
}

// NewShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {}

// ShowPrefixList displays detailed information for a specified Prefix List.
func ShowPrefixList(cmd *cobra.Command, id string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	pls, err := svc.GetPrefixLists(ctx, &ascTypes.GetPrefixListsInput{PrefixListIDs: []string{id}})
	if err != nil {
		return fmt.Errorf("get prefix lists: %w", err)
	}
	if len(pls) == 0 {
		return fmt.Errorf("Prefix List not found: %s", id)
	}
	pl := pls[0]

	fields := prefixListShowFields()
	opts := tableformat.RenderOptions{
		Title: fmt.Sprintf("Prefix List Details\n(%s)", id),
		Style: "rounded",
	}

	return tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: pl,
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			// Placeholder logic for extra fields
			switch fieldID {
			case "Prefix List ARN":
				return "-", nil // TODO: Lookup ARN
			case "State":
				return "-", nil // TODO: Lookup state
			case "Version":
				return "-", nil // TODO: Lookup version
			case "Max entries":
				return "-", nil // TODO: Lookup max entries
			case "Owner":
				return "-", nil // TODO: Lookup owner
			case "Entries":
				return "-", nil // TODO: Format entries (CIDR, Description)
			}
			return vpc.GetPrefixListAttributeValue(fieldID, instance)
		},
	}, opts)
}
