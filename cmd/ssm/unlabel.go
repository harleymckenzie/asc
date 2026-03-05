package ssm

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ssm"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	newUnlabelFlags(unlabelCmd)
}

var unlabelCmd = &cobra.Command{
	Use:     "unlabel <parameter-name> <label> [label...]",
	Short:   "Remove labels from a parameter",
	Long: `Remove labels from a parameter.

This removes the specified labels from whatever version they are attached to.
You don't need to specify the version number.`,
	GroupID: "actions",
	Args:    cobra.MinimumNArgs(2),
	Example: `  asc ssm unlabel /myapp/config prod
  asc ssm unlabel /myapp/config prod staging`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(UnlabelParameterVersion(cmd, args))
	},
}

func newUnlabelFlags(cmd *cobra.Command) {
	// No flags needed
}

func UnlabelParameterVersion(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	paramName := args[0]
	labels := args[1:]

	invalidLabels, err := svc.UnlabelParameterVersion(ctx, &ascTypes.UnlabelParameterVersionInput{
		Name:   paramName,
		Labels: labels,
	})
	if err != nil {
		return fmt.Errorf("unlabel parameter version: %w", err)
	}

	if len(invalidLabels) > 0 {
		fmt.Printf("Warning: labels not found: %v\n", invalidLabels)
	}

	// Report success for valid labels
	validCount := len(labels) - len(invalidLabels)
	if validCount > 0 {
		fmt.Printf("Removed %d label(s) from %s\n", validCount, paramName)
	}

	return nil
}
