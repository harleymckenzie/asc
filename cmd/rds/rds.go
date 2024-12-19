package rds

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc-go/cmd"
	"github.com/harleymckenzie/asc-go/pkg/service/rds"
	"github.com/spf13/cobra"
)

// RDSCmd represents the ec2 command
var RDSCmd = &cobra.Command{
	Use:   "rds",
	Short: "Perform RDS operations",
}

// rds subcommands
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all RDS clusters and instances",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()

		svc, err := rds.NewRDSService(ctx, cmd.Profile)
		if err != nil {
			log.Fatalf("Failed to initialize RDS service: %v", err)
		}

		err = svc.ListInstances(ctx)
		if err != nil {
			log.Fatalf("Error describing running instances: %v", err)
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(RDSCmd)
	RDSCmd.AddCommand(lsCmd)
}
