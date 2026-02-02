package cmd

import (
	"github.com/harleymckenzie/asc/cmd/asg"
	"github.com/harleymckenzie/asc/cmd/cloudformation"
	"github.com/harleymckenzie/asc/cmd/ec2"
	"github.com/harleymckenzie/asc/cmd/elasticache"
	"github.com/harleymckenzie/asc/cmd/elb"
	"github.com/harleymckenzie/asc/cmd/rds"
	"github.com/harleymckenzie/asc/cmd/ssm"
	"github.com/harleymckenzie/asc/cmd/vpc"

	"github.com/spf13/cobra"
)

// Global configuration variables
var (
	Profile string     // AWS profile to use for authentication
	Region  string     // AWS region to operate in
	Version = "0.0.40" // Current version of the application
)

// NewRootCmd creates and configures the root command for the AWS Simple CLI
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asc",
		Short: "AWS Simple CLI (asc) - A simplified interface for AWS operations",
	}

	// Add persistent flags for AWS configuration
	cmd.PersistentFlags().StringVarP(&Profile, "profile", "p", "", "AWS profile to use for authentication")
	cmd.PersistentFlags().StringVar(&Region, "region", "", "AWS region to operate in")
	cmd.Version = Version

	// Add service commands
	cmd.AddCommand(asg.NewASGRootCmd())
	cmd.AddCommand(cloudformation.NewCloudFormationRootCmd())
	cmd.AddCommand(ec2.NewEC2RootCmd())
	cmd.AddCommand(elasticache.NewElasticacheRootCmd())
	cmd.AddCommand(elb.NewELBRootCmd())
	cmd.AddCommand(rds.NewRDSRootCmd())
	cmd.AddCommand(ssm.NewSSMRootCmd())
	cmd.AddCommand(vpc.NewVPCRootCmd())

	// Add command groups for better organization
	cmd.AddGroup(
		&cobra.Group{
			ID:    "service",
			Title: "Service Commands",
		},
	)

	return cmd
}

// Execute runs the root command and handles any errors
// This is called by main.main() and only needs to happen once
func Execute() error {
	return NewRootCmd().Execute()
}
