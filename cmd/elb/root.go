package elb

import (
	tg "github.com/harleymckenzie/asc/cmd/elb/target_group"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewELBRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "elb",
		Short:   "Perform Elastic Load Balancer operations",
		Aliases: []string{"alb"},
		GroupID: "service",
	}
	
	// Action commands
	cmd.AddCommand(lsCmd)
	// cmd.AddCommand(addCmd)
	// cmd.AddCommand(rmCmd)
	// cmd.AddCommand(modifyCmd)
	
	// Subcommands
	cmd.AddCommand(tg.NewTargetGroupRootCmd())

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)

	return cmd
}