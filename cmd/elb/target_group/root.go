package target_group

import (
	"log"

	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewTargetGroupRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "target-group",
		Short:   "Perform target group operations",
		Aliases: []string{"tg", "target-groups", "target_group", "target_groups"},
		GroupID: "subcommands",
		Run: func(cobraCmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatalf("Please provide a subcommand: ls")
			}

			switch args[0] {
			case "ls":
				lsCmd.Run(cobraCmd, args[1:])
			default:
				log.Fatalf("Invalid subcommand: %s", args[0])
			}
		},
	}

	// Action commands
	cmd.AddCommand(lsCmd)

	// Add gropus
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
