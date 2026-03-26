package efs

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/efs"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	reverseSort bool
)

func init() {
	newLsFlags(lsCmd)
}

func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "File System Details", Visible: true, DefaultSort: true},
		{Name: "File System ID", Category: "File System Details", Visible: true},
		{Name: "State", Category: "File System Details", Visible: true},
		{Name: "Size (Bytes)", Category: "File System Details", Visible: true},
		{Name: "Mount Targets", Category: "File System Details", Visible: true},
		{Name: "Performance Mode", Category: "File System Details", Visible: true},
		{Name: "Throughput Mode", Category: "File System Details", Visible: true},
		{Name: "Encrypted", Category: "File System Details", Visible: false},
		{Name: "Availability Zone", Category: "File System Details", Visible: false},
		{Name: "Creation Time", Category: "File System Details", Visible: false},
	}
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List EFS file systems",
	Aliases: []string{"list"},
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListFileSystems(cmd, args))
	},
}

func newLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs file systems in list format.")
	cmdutil.AddTagFlag(cobraCmd)
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
}

func ListFileSystems(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, efs.NewEFSService)
	if err != nil {
		return fmt.Errorf("create efs service: %w", err)
	}

	fileSystems, err := svc.GetFileSystems(cmd.Context())
	if err != nil {
		return fmt.Errorf("get file systems: %w", err)
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         "File Systems",
		PlainStyle:    list,
		Fields:        getListFields(),
		Tags:          cmdutil.Tags,
		Data:          utils.SlicesToAny(fileSystems),
		GetFieldValue: efs.GetFieldValue,
		GetTagValue:   efs.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}
