package elasticache

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc-go/cmd"
	"github.com/harleymckenzie/asc-go/pkg/service/elasticache"
	"github.com/spf13/cobra"
)

// RDSCmd represents the ec2 command
var RDSCmd = &cobra.Command{
	Use:   "elasticache",
	Short: "Perform Elasticache operations",
}

// rds subcommands
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all Elasticache clusters",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()

		svc, err := elasticache.NewElasticacheService(ctx, cmd.Profile)
		if err != nil {
			log.Fatalf("Failed to initialize Elasticache service: %v", err)
		}

		err = svc.ListInstances(ctx)
		if err != nil {
			log.Fatalf("Error describing clusters: %v", err)
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(RDSCmd)
	RDSCmd.AddCommand(lsCmd)
}
