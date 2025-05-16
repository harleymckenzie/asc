package ec2

import (
	"github.com/harleymckenzie/asc/cmd/ec2/ami"
	"github.com/harleymckenzie/asc/cmd/ec2/security_group"
	"github.com/harleymckenzie/asc/cmd/ec2/snapshot"
	"github.com/harleymckenzie/asc/cmd/ec2/volume"

	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	// Add subcommands
	lsCmd.AddCommand(amiLsCmd)
	lsCmd.AddCommand(securityGroupLsCmd)
	lsCmd.AddCommand(snapshotLsCmd)
	lsCmd.AddCommand(volumeLsCmd)
	
	showCmd.AddCommand(amiShowCmd)
	showCmd.AddCommand(securityGroupShowCmd)
	showCmd.AddCommand(snapshotShowCmd)
	showCmd.AddCommand(volumeShowCmd)

	// Add flags
	ami.NewLsFlags(amiLsCmd)
	security_group.NewLsFlags(securityGroupLsCmd)
	snapshot.NewLsFlags(snapshotLsCmd)
	volume.NewLsFlags(volumeLsCmd)
	
	ami.NewShowFlags(amiShowCmd)
	security_group.NewShowFlags(securityGroupShowCmd)
	snapshot.NewShowFlags(snapshotShowCmd)
	volume.NewShowFlags(volumeShowCmd)

	// Add groups
	lsCmd.AddGroup(cmdutil.SubcommandGroups()...)
	showCmd.AddGroup(cmdutil.SubcommandGroups()...)
}

// Subcommand variables
var amiLsCmd = &cobra.Command{
	Use:     "amis",
	Short:   "List AMIs",
	Aliases: ami.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ami.ListAMIs(cmd, args))
	},
}

var securityGroupLsCmd = &cobra.Command{
	Use:     "security-groups",
	Short:   "List all security groups",
	Aliases: security_group.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(security_group.ListSecurityGroups(cmd, args))
	},
}

var snapshotLsCmd = &cobra.Command{
	Use:     "snapshots",
	Short:   "List all snapshots",
	Aliases: snapshot.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(snapshot.ListSnapshots(cmd, args))
	},
}

var volumeLsCmd = &cobra.Command{
	Use:     "volumes",
	Short:   "List all volumes",
	Aliases: volume.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(volume.ListVolumes(cmd, args))
	},
}

var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show detailed information about an EC2 instance",
	Aliases: []string{"describe"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowEC2Resource(cmd, args[0]))
	},
}

var amiShowCmd = &cobra.Command{
	Use:     "amis",
	Short:   "Show detailed information about an AMI",
	Aliases: ami.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ami.ShowEC2AMI(cmd, args[0]))
	},
}

var securityGroupShowCmd = &cobra.Command{
	Use:     "security-groups",
	Short:   "Show detailed information about a security group",
	Aliases: security_group.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(security_group.ShowSecurityGroup(cmd, args[0]))
	},
}

var snapshotShowCmd = &cobra.Command{
	Use:     "snapshots",
	Short:   "Show detailed information about a snapshot",
	Aliases: snapshot.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(snapshot.ShowEC2Snapshot(cmd, args[0]))
	},
}

var volumeShowCmd = &cobra.Command{
	Use:     "volumes",
	Short:   "Show detailed information about an EBS volume",
	Aliases: volume.CmdAliases,
	GroupID: "subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(volume.ShowEC2Volume(cmd, args[0]))
	},
}
