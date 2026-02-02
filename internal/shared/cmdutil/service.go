package cmdutil

import (
	"context"

	"github.com/spf13/cobra"
)

type ServiceCreator[T any] func(ctx context.Context, profile string, region string) (T, error)

func CreateService[T any](cmd *cobra.Command, createService ServiceCreator[T]) (T, error) {
	ctx := cmd.Context()
	profile, region := GetPersistentFlags(cmd)
	return createService(ctx, profile, region)
}
