package cmd

import (
	"github.com/harleymckenzie/asc/cmd/asg"
	"github.com/harleymckenzie/asc/cmd/cloudformation"
	"github.com/harleymckenzie/asc/cmd/ec2"
	"github.com/harleymckenzie/asc/cmd/elasticache"
	"github.com/harleymckenzie/asc/cmd/elb"
	"github.com/harleymckenzie/asc/cmd/rds"
	"github.com/harleymckenzie/asc/cmd/vpc"

	"github.com/spf13/cobra"
)

// Global variable to store the profile value
var (
	Profile string
	Region  string
	Version = "0.0.30"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asc",
		Short: "AWS Simple CLI (asc) - A simplified interface for AWS operations",
	}
	cmd.PersistentFlags().StringVarP(&Profile, "profile", "p", "", "AWS profile to use")
	cmd.PersistentFlags().StringVar(&Region, "region", "", "AWS region to use")
	cmd.Version = Version

	// Add commands
	cmd.AddCommand(asg.NewASGRootCmd())
	cmd.AddCommand(cloudformation.NewCloudFormationRootCmd())
	cmd.AddCommand(ec2.NewEC2RootCmd())
	cmd.AddCommand(elasticache.NewElasticacheRootCmd())
	cmd.AddCommand(elb.NewELBRootCmd())
	cmd.AddCommand(rds.NewRDSRootCmd())
	cmd.AddCommand(vpc.NewVPCRootCmd())
	// Add groups
	cmd.AddGroup(
		&cobra.Group{
			ID:    "service",
			Title: "Service Commands",
		},
	)

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return NewRootCmd().Execute()
}
