package volume

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// Root command
func NewVolumeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Short:   "Perform volume operations",
		Aliases: []string{"volumes", "vol", "volum"},
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
