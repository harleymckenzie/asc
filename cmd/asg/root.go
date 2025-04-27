package asg

import (
	"github.com/spf13/cobra"
)

func NewASGRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asg",
		Short: "Perform Auto Scaling Group operations",
	}

	cmd.AddCommand(lsCmd)

	return cmd
}