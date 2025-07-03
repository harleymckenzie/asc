package iam

import (
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show IAM resources",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Example: "  asc iam show role my-role      # Show details for the role named 'my-role'",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}