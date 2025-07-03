package iam

import (
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List IAM resources",
	Aliases: []string{"list"},
	GroupID: "actions",
	Example: "  asc iam ls roles              # List all IAM roles",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
