package elasticache

import (
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// NewElasticacheRootCmd creates and configures the root command for ElastiCache operations
func NewElasticacheRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "elasticache",
		Short:   "Perform ElastiCache operations",
		Long:    "Manage ElastiCache clusters including Redis and Memcached instances",
		Aliases: []string{"redis", "memcached", "cache"},
		GroupID: "service",
	}

	// Add action commands
	cmd.AddCommand(lsCmd)

	// Add command groups for better organization
	cmd.AddGroup(cmdutil.ActionGroups()...)

	return cmd
}
