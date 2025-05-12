package ec2

import (
	"github.com/harleymckenzie/asc/cmd/ec2/ami"
	"github.com/harleymckenzie/asc/cmd/ec2/security_group"
	"github.com/harleymckenzie/asc/cmd/ec2/snapshot"
	"github.com/harleymckenzie/asc/pkg/shared/cmdutil"
	"github.com/spf13/cobra"
)

func NewEC2RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ec2",
		Short:   "Perform EC2 operations",
		GroupID: "service",
	}

	// Add commands
	cmd.AddCommand(lsCmd)
	cmd.AddCommand(showCmd)
	cmd.AddCommand(restartCmd)
	cmd.AddCommand(startCmd)
	cmd.AddCommand(stopCmd)
	cmd.AddCommand(terminateCmd)

	// Subcommands
	cmd.AddCommand(ami.NewAMIRootCmd())
	cmd.AddCommand(snapshot.NewSnapshotRootCmd())
	cmd.AddCommand(security_group.NewSecurityGroupRootCmd())

	// Add groups
	cmd.AddGroup(cmdutil.ActionGroups()...)
	cmd.AddGroup(cmdutil.SubcommandGroups()...)
	
	return cmd
}
