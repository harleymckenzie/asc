package elasticache

import (
	"github.com/spf13/cobra"
)


func NewElasticacheRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "elasticache",
		Short:   "Perform Elasticache operations",
		GroupID: "service",
	}

	cmd.AddCommand(lsCmd)

	return cmd
}
