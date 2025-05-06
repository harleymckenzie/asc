// The ls command list Elastic Load Balancers, as well as an alias for the relevant subcommand.
// It re-uses existing functions and flags from the relevant commands.

package elb

import (
	"github.com/spf13/cobra"
)

var (
	list bool
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List Elastic Load Balancers and target groups",
	Long: "List Elastic Load Balancers and target groups\n" +
		"  ls                      List all Elastic Load Balancers\n" +
		"  ls [elb-name]           List target groups for the specified ELB\n" +
		"  ls target-groups [elb-name] List target groups for the specified ELB",
	GroupID: "actions",
	Run:     func(cobraCmd *cobra.Command, args []string) {},
}

// lsTargetGroupCmd calls the ListELBTargetGroups function, which is also used by the `target-group ls` command.
var lsTargetGroupCmd = &cobra.Command{
	Use:   "target-groups",
	Short: "List target groups for the specified ELB",
	Run: func(cobraCmd *cobra.Command, args []string) {
		lsTargetGroups(cobraCmd, args)
	},
}

func addLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Elastic Load Balancers in list format.")
}

func init() {
	addLsFlags(lsCmd)

	lsCmd.AddCommand(lsTargetGroupCmd)
}
