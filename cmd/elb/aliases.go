package elb

import (
	tg "github.com/harleymckenzie/asc/cmd/elb/target_group"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	// Add subcommands
	lsCmd.AddCommand(tg.NewTargetGroupRootCmd())

	// Add flags
	tg.NewLsFlags(lsTargetGroupCmd)

	// Add groups
	lsCmd.AddGroup(cmdutil.SubcommandGroups()...)
}

// Subcommand variable
var lsTargetGroupCmd = &cobra.Command{
	Use:   "target-groups",
	Short: "List target groups",
	Run: func(cobraCmd *cobra.Command, args []string) {
		tg.ListTargetGroups(cobraCmd, args)
	},
}