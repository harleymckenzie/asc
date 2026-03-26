package service

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"services", "svc"}

// NewServiceRootCmd creates and configures the root command for ECS service operations
func NewServiceRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service",
		Short:   "Perform ECS service operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)
	cmd.AddCommand(waitCmd)

	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
