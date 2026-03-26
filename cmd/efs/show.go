package efs

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	efsTypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
	"github.com/harleymckenzie/asc/internal/service/efs"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	newShowFlags(showCmd)
}

func getShowFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "File System ID", Category: "General", Visible: true},
		{Name: "State", Category: "General", Visible: true},
		{Name: "Creation Time", Category: "General", Visible: true},
		{Name: "Owner ID", Category: "General", Visible: true},
		{Name: "ARN", Category: "General", Visible: true},

		{Name: "Size (Bytes)", Category: "Storage", Visible: true},
		{Name: "Performance Mode", Category: "Storage", Visible: true},
		{Name: "Throughput Mode", Category: "Storage", Visible: true},
		{Name: "Provisioned Throughput", Category: "Storage", Visible: true},

		{Name: "Mount Targets", Category: "Network", Visible: true},
		{Name: "Availability Zone", Category: "Network", Visible: true},

		{Name: "Encrypted", Category: "Encryption", Visible: true},
		{Name: "KMS Key ID", Category: "Encryption", Visible: true},
	}
}

var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an EFS file system",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowFileSystem(cmd, args[0]))
	},
}

func newShowFlags(cmd *cobra.Command) {
	cmdutil.AddShowFlags(cmd, "vertical")
}

func ShowFileSystem(cmd *cobra.Command, fileSystemID string) error {
	svc, err := cmdutil.CreateService(cmd, efs.NewEFSService)
	if err != nil {
		return fmt.Errorf("create efs service: %w", err)
	}

	fs, err := svc.GetFileSystem(cmd.Context(), fileSystemID)
	if err != nil {
		return fmt.Errorf("get file system: %w", err)
	}

	title := fmt.Sprintf("File System Details\n(%s)", fileSystemID)
	if fs.Name != nil {
		title = fmt.Sprintf("File System Details\n(%s)", aws.ToString(fs.Name))
	}

	table := tablewriter.NewDetailTable(tablewriter.AscTableRenderOptions{
		Title:   title,
		Columns: 3,
	})

	fields, err := tablewriter.PopulateFieldValues(fs, getShowFields(), efs.GetFieldValue)
	if err != nil {
		return fmt.Errorf("populate field values: %w", err)
	}

	layout := tablewriter.Horizontal
	if cmdutil.GetLayout(cmd) == "grid" {
		layout = tablewriter.Grid
	}

	table.AddSections(tablewriter.BuildSections(fields, layout))
	tags := populateEFSTagFields(fs.Tags)
	table.AddSection(tablewriter.BuildSection("Tags", tags, tablewriter.Horizontal))

	table.Render()
	return nil
}

func populateEFSTagFields(tags []efsTypes.Tag) []tablewriter.Field {
	var fields []tablewriter.Field
	for _, tag := range tags {
		if tag.Key != nil && tag.Value != nil {
			fields = append(fields, tablewriter.Field{
				Category: "Tag",
				Name:     aws.ToString(tag.Key),
				Value:    aws.ToString(tag.Value),
			})
		}
	}
	return fields
}
