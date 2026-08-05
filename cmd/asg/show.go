// The show command displays detailed information about an Auto Scaling Group,
// including details of the launch template it is configured with.

package asg

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/harleymckenzie/asc/internal/service/asg"
	ascTypes "github.com/harleymckenzie/asc/internal/service/asg/types"
	"github.com/harleymckenzie/asc/internal/service/ec2"
	ec2Types "github.com/harleymckenzie/asc/internal/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

// Init function
func init() {
	newShowFlags(showCmd)
}

// getShowFields returns the Auto Scaling Group detail fields.
func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Group Details", Visible: true},
		{Name: "ARN", Category: "Group Details", Visible: true},
		{Name: "Status", Category: "Group Details", Visible: true},
		{Name: "Created Time", Category: "Group Details", Visible: true},

		{Name: "Desired", Category: "Capacity", Visible: true},
		{Name: "Min", Category: "Capacity", Visible: true},
		{Name: "Max", Category: "Capacity", Visible: true},
		{Name: "Instances", Category: "Capacity", Visible: true},
		{Name: "Default Cooldown", Category: "Capacity", Visible: true},

		{Name: "Health Check Type", Category: "Health Check", Visible: true},
		{Name: "Health Check Grace Period", Category: "Health Check", Visible: true},

		{Name: "Availability Zones", Category: "Network", Visible: true},
		{Name: "Subnets", Category: "Network", Visible: true},

		{Name: "Launch Template", Category: "Launch Configuration", Visible: true},
		{Name: "Launch Configuration", Category: "Launch Configuration", Visible: true},
	}
}

// getLaunchTemplateFields returns the fields shown for the associated launch template version.
func getLaunchTemplateFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Launch Template Name", Category: "Launch Template", Visible: true},
		{Name: "Launch Template ID", Category: "Launch Template", Visible: true},
		{Name: "Version", Category: "Launch Template", Visible: true},
		{Name: "Default Version", Category: "Launch Template", Visible: true},
		{Name: "Instance Type", Category: "Launch Template", Visible: true},
		{Name: "AMI ID", Category: "Launch Template", Visible: true},
		{Name: "Key Name", Category: "Launch Template", Visible: true},
		{Name: "IAM Instance Profile", Category: "Launch Template", Visible: true},
		{Name: "Security Group IDs", Category: "Launch Template", Visible: true},
	}
}

// showCmd is the command for showing detailed information about an Auto Scaling Group.
var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an Auto Scaling Group",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowAutoScalingGroup(cobraCmd, args[0]))
	},
}

// newShowFlags is the function for adding flags to the show command.
func newShowFlags(cobraCmd *cobra.Command) {
	cmdutil.AddShowFlags(cobraCmd, "vertical")
}

// ShowAutoScalingGroup displays detailed information for an Auto Scaling Group,
// including details of its associated launch template.
func ShowAutoScalingGroup(cmd *cobra.Command, name string) error {
	svc, err := cmdutil.CreateService(cmd, asg.NewAutoScalingService)
	if err != nil {
		return fmt.Errorf("create new Auto Scaling Group service: %w", err)
	}

	groups, err := svc.GetAutoScalingGroups(cmd.Context(), &ascTypes.GetAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{name},
	})
	if err != nil {
		return fmt.Errorf("get Auto Scaling Group: %w", err)
	}
	if len(groups) == 0 {
		return fmt.Errorf("no Auto Scaling Group found for %s", name)
	}
	group := groups[0]

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   fmt.Sprintf("Auto Scaling Group Details\n(%s)", name),
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(group, getShowFields(), asg.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}
	table.AddSections(tablewriter.BuildSections(fields, layout))

	// Add a section with the details of the associated launch template, if the group uses one.
	if spec := asg.GetLaunchTemplateSpecification(group); spec != nil {
		ltFields, err := launchTemplateFields(cmd, spec)
		if err != nil {
			return err
		}
		if ltFields != nil {
			table.AddSection(tablewriter.BuildSection("Launch Template", ltFields, layout))
		}
	}

	table.Render()
	return nil
}

// launchTemplateFields resolves the launch template version referenced by the Auto Scaling
// Group and returns its populated fields. It returns nil fields if the version cannot be found.
func launchTemplateFields(cmd *cobra.Command, spec *autoscalingtypes.LaunchTemplateSpecification) ([]tablewriter.Field, error) {
	ec2Svc, err := cmdutil.CreateService(cmd, ec2.NewEC2Service)
	if err != nil {
		return nil, fmt.Errorf("create ec2 service: %w", err)
	}

	ltVersion := aws.ToString(spec.Version)
	if ltVersion == "" {
		ltVersion = "$Default"
	}

	versions, err := ec2Svc.GetLaunchTemplateVersions(cmd.Context(), &ec2Types.GetLaunchTemplateVersionsInput{
		LaunchTemplateID:   aws.ToString(spec.LaunchTemplateId),
		LaunchTemplateName: aws.ToString(spec.LaunchTemplateName),
		Versions:           []string{ltVersion},
	})
	if err != nil {
		return nil, fmt.Errorf("get launch template versions: %w", err)
	}
	if len(versions) == 0 {
		return nil, nil
	}

	fields, err := tablewriter.PopulateFieldValues(versions[0], getLaunchTemplateFields(), ec2.GetFieldValue)
	if err != nil {
		return nil, fmt.Errorf("populate launch template field values: %w", err)
	}
	return fields, nil
}
