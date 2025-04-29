package cmd

import (
	"os"

	"github.com/harleymckenzie/asc/cmd/asg"
	"github.com/harleymckenzie/asc/cmd/cloudformation"
	"github.com/harleymckenzie/asc/cmd/ec2"
	"github.com/harleymckenzie/asc/cmd/elasticache"
	"github.com/harleymckenzie/asc/cmd/rds"
	"github.com/spf13/cobra"
)

// Global variable to store the profile value
var (
	Profile string
	Region  string
	Version = "0.0.17"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asc",
		Short: "AWS Simple CLI (asc) - A simplified interface for AWS operations",
	}
	cmd.PersistentFlags().StringVarP(&Profile, "profile", "p", "", "AWS profile to use")
	cmd.PersistentFlags().StringVarP(&Region, "region", "r", "", "AWS region to use")
	cmd.Version = Version

	// Add commands
	cmd.AddCommand(asg.NewASGRootCmd())
	cmd.AddCommand(cloudformation.NewCloudFormationRootCmd())
	cmd.AddCommand(ec2.NewEC2RootCmd())
	cmd.AddCommand(rds.NewRDSRootCmd())
	cmd.AddCommand(elasticache.NewElasticacheRootCmd())

	// Add groups
	cmd.AddGroup(
		&cobra.Group{
			ID: "service",
			Title: "Service Commands",
		},
	)

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
