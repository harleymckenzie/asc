// show.go displays detailed information about an AMI.
package ami

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
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

// ec2AMIShowFields returns the fields for the AMI detail table.
func ec2AMIShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "AMI Name", Display: true},
		{ID: "Image Type", Display: true},
		{ID: "Platform", Display: false},
		{ID: "Root Device Type", Display: false},
		{ID: "Owner", Display: true},
		{ID: "Architecture", Display: true},
		{ID: "Usage Operation", Display: true},
		{ID: "Root Device Name", Display: true},
		{ID: "Status", Display: true},
		{ID: "Source", Display: false},
		{ID: "Virtualization", Display: true},
		{ID: "Boot Mode", Display: true},
		{ID: "State Reason", Display: true},
		{ID: "Creation Date", Display: true},
		{ID: "Kernel ID", Display: true},
		{ID: "Description", Display: true},
		{ID: "Product Codes", Display: true},
		{ID: "RAM Disk ID", Display: true},
		{ID: "Deprecation Time", Display: true},
		{ID: "Block Devices", Display: false},
		{ID: "Deregistration Protection", Display: true},
		{ID: "Allowed Image", Display: true},
		{ID: "Source AMI ID", Display: true},
		{ID: "Source AMI Region", Display: true},
		{ID: "Visibility", Display: false},
	}
}

// showCmd is the cobra command for showing AMI details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an AMI",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowEC2AMI(cmd, args[0]))
	},
}

// NewShowFlags adds flags for the show subcommand.
func NewShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "horizontal")
}

// ShowEC2AMI displays detailed information for a specified AMI.
func ShowEC2AMI(cmd *cobra.Command, arg string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
	}

	if cmd.Flags().Changed("output") {
		if err := cmdutil.ValidateFlagChoice(cmd, "output", cmdutil.ValidLayouts); err != nil {
			return err
		}
	}

	images, err := svc.GetImages(ctx, &ascTypes.GetImagesInput{ImageIDs: []string{arg}})
	if err != nil {
		return fmt.Errorf("get images: %w", err)
	}
	if len(images) == 0 {
		return fmt.Errorf("AMI not found: %s", arg)
	}

	fields := ec2AMIShowFields()
	opts := tableformat.RenderOptions{
		Title: fmt.Sprintf("AMI Details\n(%s)", aws.ToString(images[0].ImageId)),
		Style: "rounded",
		Layout: tableformat.DetailTableLayout{
			Type:           cmdutil.GetLayout(cmd),
			ColumnsPerRow:  3,
		},
		SortBy: tableformat.GetSortByField(fields, false),
	}

	return tableformat.RenderTableDetail(&tableformat.DetailTable{
		Instance: images[0],
		Fields:   fields,
		GetAttribute: func(fieldID string, image any) (string, error) {
			return ec2.GetImageAttributeValue(fieldID, image)
		},
	}, opts)
}
