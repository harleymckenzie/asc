package volume

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var CmdAliases = []string{"volumes", "volume", "disk", "disks", "ebs"}

// Root command
func NewVolumeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Short:   "Perform volume operations",
		Aliases: CmdAliases,
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
