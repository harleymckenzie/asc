package rds

import (
	"github.com/spf13/cobra"
)

func NewRDSRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rds",
		Short:   "Perform RDS operations",
		GroupID: "service",
	}

	cmd.AddCommand(lsCmd)

	return cmd
}
