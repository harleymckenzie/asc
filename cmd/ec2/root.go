package ec2

import "github.com/spf13/cobra"

func NewEC2RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ec2",
		Short: "Perform EC2 operations",
	}

	cmd.AddCommand(lsCmd)

	return cmd
}
