// The show command displays detailed information about an VPC internet gateway.

package igw

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/vpc"
	ascTypes "github.com/harleymckenzie/asc/internal/service/vpc/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {
	NewShowFlags(showCmd)
}

// Column functions
func vpcIGWShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Internet Gateway ID", Display: true},
		{ID: "VPC ID", Display: true},
		{ID: "Owner", Display: true},
		{ID: "State", Display: true},
		{ID: "VPC Attachments", Display: true},
	}
}

// Command variable
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an internet gateway",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowVPCIGW(cmd, args[0]))
	},
}

// Flag function
func NewShowFlags(cobraCmd *cobra.Command) {}

// ShowVPCIGW is the function for showing VPC internet gateways
func ShowVPCIGW(cmd *cobra.Command, id string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := vpc.NewVPCService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new VPC service: %w", err)
	}

	igws, err := svc.GetIGWs(ctx, &ascTypes.GetIGWsInput{IGWIDs: []string{id}})
	if err != nil {
		return fmt.Errorf("get internet gateways: %w", err)
	}
	if len(igws) == 0 {
		return fmt.Errorf("Internet Gateway not found: %s", id)
	}
	igw := igws[0]

	fields := vpcIGWShowFields()
	opts := tableformat.RenderOptions{
		Title: fmt.Sprintf("Internet Gateway Details\n(%s)", id),
		Style: "rounded",
	}

	return tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: igw,
		Fields:   fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			// Placeholder logic for extra fields
			switch fieldID {
			case "VPC ID":
				return "-", nil // TODO: Show attached VPC ID
			case "Owner":
				return "-", nil // TODO: Lookup Owner
			case "State":
				if len(igw.Attachments) > 0 {
					return string(igw.Attachments[0].State), nil
				}
				return "-", nil
			}
			return vpc.GetIGWAttributeValue(fieldID, instance)
		},
	}, opts)
}
