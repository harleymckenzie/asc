package prefix_list

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"prefix-lists", "prefix-list", "prefixlist", "prefixlists"}

// NewPrefixListRootCmd returns the root command for Prefix List operations.
func NewPrefixListRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "prefix-list",
		Short:   "Perform Prefix List operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	// cmd.AddCommand(lsCmd) // Disabled - ls.go.disabled
	cmd.AddCommand(showCmd)

	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)

	return cmd
}
