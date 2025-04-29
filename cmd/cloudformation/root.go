package cloudformation

import "github.com/spf13/cobra"

func NewCloudFormationRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cloudformation",
		Short:   "Perform CloudFormation operations",
		GroupID: "service",
	}

	cmd.AddCommand(lsCmd)

	return cmd
}
