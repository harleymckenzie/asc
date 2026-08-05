// show.go displays detailed information about a launch template.
package launch_template

import (
	"fmt"
	"strings"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

// Variables
var (
	version string
)

// Init function
func init() {
	newShowFlags(showCmd)
}

func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Launch Template Name", Category: "Launch Template Details", Visible: true},
		{Name: "Launch Template ID", Category: "Launch Template Details", Visible: true},
		{Name: "Version", Category: "Launch Template Details", Visible: true},
		{Name: "Default Version", Category: "Launch Template Details", Visible: true},
		{Name: "Version Description", Category: "Launch Template Details", Visible: true},
		{Name: "Created Time", Category: "Launch Template Details", Visible: true},
		{Name: "Created By", Category: "Launch Template Details", Visible: false},

		{Name: "Instance Type", Category: "Instance Configuration", Visible: true},
		{Name: "AMI ID", Category: "Instance Configuration", Visible: true},
		{Name: "Key Name", Category: "Instance Configuration", Visible: true},
		{Name: "IAM Instance Profile", Category: "Instance Configuration", Visible: true},
		{Name: "Monitoring", Category: "Instance Configuration", Visible: true},
		{Name: "EBS Optimized", Category: "Instance Configuration", Visible: true},
		{Name: "Security Group IDs", Category: "Instance Configuration", Visible: true},
		{Name: "Security Groups", Category: "Instance Configuration", Visible: false},
	}
}

// showCmd is the cobra command for showing launch template details.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about a launch template",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowLaunchTemplate(cmd, args[0]))
	},
}

// newShowFlags adds flags for the show subcommand.
func newShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "horizontal")
	cobraCmd.Flags().StringVar(&version, "version", "$Default", "The launch template version to show ($Default, $Latest, or a version number).")
}

// ShowLaunchTemplate displays detailed information for a launch template.
// The identifier may be a launch template ID (lt-...) or name.
func ShowLaunchTemplate(cmd *cobra.Command, arg string) error {
	svc, err := cmdutil.CreateService(cmd, ec2.NewEC2Service)
	if err != nil {
		return fmt.Errorf("create ec2 service: %w", err)
	}

	input := &ascTypes.GetLaunchTemplateVersionsInput{
		Versions: []string{version},
	}
	if strings.HasPrefix(arg, "lt-") {
		input.LaunchTemplateID = arg
	} else {
		input.LaunchTemplateName = arg
	}

	versions, err := svc.GetLaunchTemplateVersions(cmd.Context(), input)
	if err != nil {
		return fmt.Errorf("get launch template versions: %w", err)
	}
	if len(versions) == 0 {
		return fmt.Errorf("no launch template version found for %s", arg)
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:          "Launch Template Details\n(" + arg + ")",
		Columns:        3,
		MaxColumnWidth: 90,
	})
	fields, err := tablewriter.PopulateFieldValues(versions[0], getShowFields(), ec2.GetFieldValue)
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

	table.Render()
	return nil
}
