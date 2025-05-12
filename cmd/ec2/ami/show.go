// show.go displays detailed information about an AMI.
package ami

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/harleymckenzie/asc/pkg/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

// newShowFlags adds flags for the show subcommand.
func newShowFlags(cobraCmd *cobra.Command) {}

// ec2AMIShowFields returns the fields for the AMI detail table.
func ec2AMIShowFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "AMI Name", Visible: true},
		{ID: "Image Type", Visible: true},
		{ID: "Platform", Visible: false},
		{ID: "Root Device Type", Visible: false},
		{ID: "Owner", Visible: true},
		{ID: "Architecture", Visible: true},
		{ID: "Usage Operation", Visible: true},
		{ID: "Root Device Name", Visible: true},
		{ID: "Status", Visible: true},
		{ID: "Source", Visible: false},
		{ID: "Virtualization", Visible: true},
		{ID: "Boot Mode", Visible: true},
		{ID: "State Reason", Visible: true},
		{ID: "Creation Date", Visible: true},
		{ID: "Kernel ID", Visible: true},
		{ID: "Description", Visible: true},
		{ID: "Product Codes", Visible: true},
		{ID: "RAM Disk ID", Visible: true},
		{ID: "Deprecation Time", Visible: true},
		{ID: "Block Devices", Visible: false},
		{ID: "Deregistration Protection", Visible: true},
		{ID: "Allowed Image", Visible: true},
		{ID: "Source AMI ID", Visible: true},
		{ID: "Source AMI Region", Visible: true},
		{ID: "Visibility", Visible: false},
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

// ShowEC2AMI displays detailed information for a specified AMI.
func ShowEC2AMI(cmd *cobra.Command, arg string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cmd)
	svc, err := ec2.NewEC2Service(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new EC2 service: %w", err)
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
		Title:  fmt.Sprintf("AMI Details\n(%s)", aws.ToString(images[0].ImageId)),
		Style:  "rounded",
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
