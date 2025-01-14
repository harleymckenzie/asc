package elasticache

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/elasticache"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var showEndpoint bool

// elasticache subcommands
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all Elasticache clusters",
	PreRun: func(cobraCmd *cobra.Command, args []string) {
		// Clear any existing sort order
		sortOrder = []string{}

		// Visit flags in the order they appear in the command line
		cobraCmd.Flags().Visit(func(f *pflag.Flag) {
			switch f.Name {
			case "sort-name":
				sortOrder = append(sortOrder, "name")
			case "sort-type":
				sortOrder = append(sortOrder, "instance_type")
			}
		})
	},
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()

		svc, err := elasticache.NewElasticacheService(ctx, "akoovadev")
		if err != nil {
			log.Fatalf("Failed to initialize Elasticache service: %v", err)
		}

		err = svc.ListInstances(ctx, sortOrder, list, showEndpoint)
		if err != nil {
			log.Fatalf("Error describing clusters: %v", err)
		}
	},
}

var (
	sortOrder    []string
	list         bool
)

func NewElasticacheCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "elasticache",
		Short: "Perform Elasticache operations",
	}

	// ls sub command
	lsCmd := &cobra.Command{
		Use:   "ls",
		Short: "List all Elasticache clusters",
		PreRun: func(cobraCmd *cobra.Command, args []string) {
			// Clear any existing sort order
			sortOrder = []string{}

			// Visit flags in the order they appear in the command line
			cobraCmd.Flags().Visit(func(f *pflag.Flag) {
				switch f.Name {
				case "sort-name":
					sortOrder = append(sortOrder, "name")
				case "sort-type":
					sortOrder = append(sortOrder, "instance_type")
				}
			})
		},
		Run: func(cobraCmd *cobra.Command, args []string) {
			ctx := context.TODO()

			svc, err := elasticache.NewElasticacheService(ctx, "akoovadev")
			if err != nil {
				log.Fatalf("Failed to initialize Elasticache service: %v", err)
			}

			err = svc.ListInstances(ctx, sortOrder, list, showEndpoint)
			if err != nil {
				log.Fatalf("Error describing clusters: %v", err)
			}
		},
	}

	cmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolVarP(&list, "list", "l", false,
		"Outputs Elasticache clusters in list format.")
	lsCmd.Flags().BoolVarP(&showEndpoint, "endpoint", "e", false,
		"Show the endpoint of the cluster")
	lsCmd.Flags().BoolP("sort-name", "n", true,
		"Sort by descending Elasticache cluster name.")
	lsCmd.Flags().BoolP("sort-type", "T", false,
		"Sort by descending Elasticache cluster type.")
	lsCmd.Flags().SortFlags = false

	return cmd
}
