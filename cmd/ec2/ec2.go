package ec2

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc-go/cmd"
	"github.com/harleymckenzie/asc-go/pkg/service/ec2"
	"github.com/spf13/cobra"
)

// EC2Cmd represents the ec2 command
var EC2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Perform EC2 operations",
}

// ec2 subcommands
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all EC2 instances",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()

		svc, err := ec2.NewEC2Service(ctx, cmd.Profile)
		if err != nil {
			log.Fatalf("Failed to initialize EC2 service: %v", err)
		}

		err = svc.ListInstances(ctx)
		if err != nil {
			log.Fatalf("Error describing running instances: %v", err)
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(EC2Cmd)
	EC2Cmd.AddCommand(lsCmd)
}
