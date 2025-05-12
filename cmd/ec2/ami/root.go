package ami

import (
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
	"github.com/spf13/cobra"
)

// Root command
func NewAMIRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ami",
		Short:   "Perform AMI operations",
		Aliases: []string{"amis", "images", "image"},
		GroupID: "subcommands",
	}

	// Add the subcommands to the command
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
