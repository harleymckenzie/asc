package nacl

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

// naclShowFields returns the fields for the NACL detail table.
func naclShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Network ACL ID", Display: true},
		{ID: "VPC ID", Display: true},
		{ID: "Is Default", Display: true},
		{ID: "Owner", Display: true},
		{ID: "Associated with", Display: true},
		{ID: "Entry Count", Display: true},
		{ID: "Association Count", Display: true},
		{ID: "Inbound rules", Display: true},
		{ID: "Outbound rules", Display: true},
	}
}

// showCmd is the cobra command for showing NACL details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a Network ACL",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowNACL(cmd, args[0]))
	},
}

// NewShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {}

// ShowNACL displays detailed information for a specified NACL.
func ShowNACL(cmd *cobra.Command, id string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	nacls, err := svc.GetNACLs(ctx, &ascTypes.GetNACLsInput{NACLIDs: []string{id}})
	if err != nil {
		return fmt.Errorf("get network acls: %w", err)
	}
	if len(nacls) == 0 {
		return fmt.Errorf("Network ACL not found: %s", id)
	}
	nacl := nacls[0]

	fields := naclShowFields()
	opts := tableformat.RenderOptions{
		Title: fmt.Sprintf("Network ACL Details\n(%s)", id),
		Style: "rounded",
	}

	return tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: nacl,
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			// Placeholder logic for extra fields
			switch fieldID {
			case "Owner":
				return "-", nil // TODO: Lookup Owner
			case "Associated with":
				return "-", nil // TODO: List associated subnets
			case "Inbound rules":
				return "-", nil // TODO: Format inbound rules
			case "Outbound rules":
				return "-", nil // TODO: Format outbound rules
			}
			return vpc.GetNetworkAclAttributeValue(fieldID, instance)
		},
	}, opts)
}
