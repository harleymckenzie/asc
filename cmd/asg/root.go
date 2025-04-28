package asg

import (
	"github.com/spf13/cobra"
)

func NewASGRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asg",
		Short: "Perform Auto Scaling Group operations",
	}

	// Action commands
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(addCmd)
	cmd.AddCommand(rmCmd)

	// Subcommands
	cmd.AddCommand(scheduleCmd)

	// Groups
	cmd.AddGroup(
		&cobra.Group{
			ID: "actions",
			Title: "Auto Scaling Group Action Commands",
		},
		&cobra.Group{
			ID: "subcommands",
			Title: "Auto Scaling Group Subcommands",
		},
	)

	return cmd
}