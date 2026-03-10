package ecs

import (
	"github.com/harleymckenzie/asc/cmd/ecs/cluster"
	"github.com/harleymckenzie/asc/cmd/ecs/service"
	"github.com/harleymckenzie/asc/cmd/ecs/task"
	"github.com/harleymckenzie/asc/cmd/ecs/taskdefinition"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// NewECSRootCmd creates and configures the root command for ECS operations
func NewECSRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ecs",
		Short:   "Perform ECS operations",
		Long:    "Manage Amazon ECS clusters, services, tasks, and task definitions",
		GroupID: "service",
	}

	// Subcommands
	cmd.AddCommand(cluster.NewClusterRootCmd())
	cmd.AddCommand(service.NewServiceRootCmd())
	cmd.AddCommand(task.NewTaskRootCmd())
	cmd.AddCommand(taskdefinition.NewTaskDefinitionRootCmd())

	// Add command groups for better organization
	cmd.AddGroup(cmdutil.SubcommandGroups()...)

	return cmd
}
