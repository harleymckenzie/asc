package snapshot

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"snapshots"}

// NewSnapshotRootCmd creates the root command for RDS snapshot operations.
func NewSnapshotRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "snapshot",
		Short:   "Perform RDS snapshot operations",
		Aliases: CmdAliases,
		GroupID: "subcommands",
	}

	// Add the subcommands to the command
	cmd.AddCommand(createCmd)
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(rmCmd)
	cmd.AddCommand(showCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
