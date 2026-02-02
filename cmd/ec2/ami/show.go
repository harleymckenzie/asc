// show.go displays detailed information about an AMI.
package ami

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

// Variables
var ()

// Init function
func init() {
	NewShowFlags(showCmd)
}

func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "AMI Name", Category: "Image Details", Visible: true},
		{Name: "Image Type", Category: "Image Details", Visible: true},
		{Name: "Platform", Category: "Image Details", Visible: false},
		{Name: "Root Device Type", Category: "Image Details", Visible: false},
		{Name: "Owner", Category: "Image Details", Visible: true},
		{Name: "Architecture", Category: "Image Details", Visible: true},
		{Name: "Usage Operation", Category: "Image Details", Visible: true},
		{Name: "Root Device Name", Category: "Image Details", Visible: true},
		{Name: "Status", Category: "Image Details", Visible: true},
		{Name: "Source", Category: "Image Details", Visible: false},
		{Name: "Virtualization", Category: "Image Details", Visible: true},
		{Name: "Boot Mode", Category: "Image Details", Visible: true},
		{Name: "State Reason", Category: "Image Details", Visible: true},
		{Name: "Creation Date", Category: "Image Details", Visible: true},
		{Name: "Kernel ID", Category: "Image Details", Visible: true},
		{Name: "Description", Category: "Image Details", Visible: true},
		{Name: "Product Codes", Category: "Image Details", Visible: true},
		{Name: "RAM Disk ID", Category: "Image Details", Visible: true},
		{Name: "Deprecation Time", Category: "Image Details", Visible: true},
		{Name: "Block Devices", Category: "Image Details", Visible: false},
		{Name: "Deregistration Protection", Category: "Image Details", Visible: true},
		{Name: "Allowed Image", Category: "Image Details", Visible: true},
		{Name: "Source AMI ID", Category: "Image Details", Visible: true},
		{Name: "Source AMI Region", Category: "Image Details", Visible: true},
		{Name: "Visibility", Category: "Image Details", Visible: false},
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

func ShowEC2AMI(cmd *cobra.Command, arg string) error {
	svc, err := cmdutil.CreateService(cmd, ec2.NewEC2Service)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	image, err := getImages(cmd.Context(), svc, &ascTypes.GetImagesInput{
		ImageIds: []string{arg},
	})
	if err != nil {
		return fmt.Errorf("get instances: %w", err)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:          "AMI Details\n(" + *image[0].ImageId + ")",
		Columns:        3,
		MaxColumnWidth: 90,
	})
	fields, err := tablewriter.PopulateFieldValues(image[0], getShowFields(), ec2.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	// Layout = Horizontal or Grid
	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
		table.Options.MaxColumnWidth = 50
	}
	table.AddSections(tablewriter.BuildSections(fields, layout))
	tags, err := awsutil.PopulateTagFields(image[0].Tags)
	if err != nil {
		return fmt.Errorf("unable to retrieve tags from instance: %w", err)
	}
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}
