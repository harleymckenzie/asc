package asg

import (
	schedule "github.com/harleymckenzie/asc/cmd/asg/schedule"
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewASGRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "asg",
		Short:   "Perform Auto Scaling Group operations",
		GroupID: "service",
	}

	// Action commands
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(addCmd)
	cmd.AddCommand(rmCmd)
	cmd.AddCommand(modifyCmd)

	// Subcommands
	cmd.AddCommand(schedule.NewScheduleRootCmd())

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)

	return cmd
}
