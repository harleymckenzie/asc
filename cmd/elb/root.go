package elb

import (
	"github.com/spf13/cobra"
)

func NewELBRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "elb",
		Short:   "Perform Elastic Load Balancer operations",
		Aliases: []string{"alb"},
		GroupID: "service",
	}
	
	// Action commands
	cmd.AddCommand(lsCmd)
	// cmd.AddCommand(addCmd)
	// cmd.AddCommand(rmCmd)
	// cmd.AddCommand(modifyCmd)
	
	// Subcommands
	cmd.AddCommand(targetGroupCmd)
	
	// Groups
	cmd.AddGroup(
		&cobra.Group{
			ID: "actions",
			Title: "Elastic Load Balancer Action Commands",
		},
		&cobra.Group{
			ID: "subcommands",
			Title: "Elastic Load Balancer Subcommands",
		},
	)

	return cmd
}