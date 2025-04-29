package ec2

import "github.com/spf13/cobra"

func NewEC2RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ec2",
		Short:   "Perform EC2 operations",
		GroupID: "service",
	}

	cmd.AddCommand(lsCmd)
	cmd.AddCommand(restartCmd)
	cmd.AddCommand(startCmd)
	cmd.AddCommand(stopCmd)
	cmd.AddCommand(terminateCmd)

	return cmd
}
