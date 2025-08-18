// The show command displays detailed information about an ELB.

package elb

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/elb"
	ascTypes "github.com/harleymckenzie/asc/internal/service/elb/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/spf13/cobra"
)

// Variables
var (
	showCmd = &cobra.Command{
		Use:     "show",
		Short:   "Show detailed information about an Elastic Load Balancer",
		Args:    cobra.ExactArgs(1),
		GroupID: "actions",
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmdutil.DefaultErrorHandler(ShowELB(cobraCmd, args))
		},
	}
)

// Init function
func init() {
	newShowFlags(showCmd)
}

// Column functions
func elbShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "ELB Details", Header: true},
		{ID: "ELB Name", Display: true},
		{ID: "Type", Display: true},
		{ID: "State", Display: true},
		{ID: "Scheme", Display: true},
		{ID: "Hosted Zone", Display: true},
		{ID: "VPC ID", Display: true},
		{ID: "Subnets", Display: true},
		{ID: "IP Type", Display: true},
		{ID: "Created Time", Display: true},
		{ID: "ARN", Display: true},
		{ID: "DNS Name", Display: true},

		{ID: "Attributes", Header: true},
		{ID: "Attribute Name", Display: true},
		{ID: "Attribute Value", Display: true},
	}
}

// Flag function
func newShowFlags(cmd *cobra.Command) {
	cmdutil.AddShowFlags(cmd, "horizontal")
}

// ShowELB is the function for showing ELBs
func ShowELB(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := elb.NewELBService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new ELB service: %w", err)
	}

	if cmd.Flags().Changed("output") {
		if err := cmdutil.ValidateFlagChoice(cmd, "output", cmdutil.ValidLayouts); err != nil {
			return err
		}
	}

	elbs, err := svc.GetLoadBalancers(ctx, &ascTypes.GetLoadBalancersInput{
		ListLoadBalancersInput: ascTypes.ListLoadBalancersInput{
			Names: []string{args[0]},
		},
	})
	if err != nil {
		return fmt.Errorf("get ELB: %w", err)
	}

	attributes, err := svc.GetLoadBalancerAttributes(ctx, *elbs[0].LoadBalancerArn)
	if err != nil {
		return fmt.Errorf("get ELB attributes: %w", err)
	}

	opts := tableformat.RenderOptions{
		Title: "ELB summary for " + *elbs[0].LoadBalancerName,
		Style: "rounded",
		Layout: tableformat.DetailTableLayout{
			Type:           cmdutil.GetLayout(cmd),
			ColumnsPerRow:  3,
			ColumnMaxWidth: 50,
		},
	}

	err = tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: elbs[0],
		Fields:   elbShowFields(),
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return elb.GetAttributeValue(fieldID, instance)
		},
	}, opts)
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}
	return nil
}
