package elb

import (
	"log"

	"github.com/spf13/cobra"
)

var targetGroupCmd = &cobra.Command{
	Use:     "target-group",
	Short:   "Perform target group operations",
	GroupID: "subcommands",
    Run: func(cobraCmd *cobra.Command, args []string) {
        if len(args) == 0 {
            log.Fatalf("Please provide a subcommand: ls")
        }
        
        switch args[0] {
        case "ls":
            targetGroupLsCmd.Run(cobraCmd, args[1:])
        default:
            log.Fatalf("Invalid subcommand: %s", args[0])
        }
    },
}

func init() {
    targetGroupCmd.AddCommand(targetGroupLsCmd)
}
