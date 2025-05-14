package snapshot

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// Root command
func NewSnapshotRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "snapshot",
		Short:   "Perform snapshot operations",
		Aliases: []string{"snapshots"},
		GroupID: "subcommands",
	}

	// Add the subcommands to the command
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)

	return cmd
}
