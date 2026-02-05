package ssm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/harleymckenzie/asc/internal/service/ssm"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	newLabelFlags(labelCmd)
}

var labelCmd = &cobra.Command{
	Use:     "label <parameter-name:version> <label> [label...]",
	Short:   "Add labels to a parameter version",
	Long: `Add labels to a specific parameter version.

Labels allow you to reference parameter versions by name instead of number.
For example, after labeling version 5 as "prod", you can access it as:
  /myapp/config:prod

A parameter version can have up to 10 labels.
Labels must start with a letter or number, and can contain letters, numbers,
periods (.), hyphens (-), and underscores (_).`,
	GroupID: "actions",
	Args:    cobra.MinimumNArgs(2),
	Example: `  asc ssm label /myapp/config:5 prod
  asc ssm label /myapp/config:5 prod stable
  asc ssm label /myapp/config:3 staging`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(LabelParameterVersion(cmd, args))
	},
}

func newLabelFlags(cmd *cobra.Command) {
	// No flags needed
}

func LabelParameterVersion(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	// Parse parameter name and version
	paramWithVersion := args[0]
	labels := args[1:]

	paramName, version, err := parseParamVersion(paramWithVersion)
	if err != nil {
		return err
	}

	invalidLabels, err := svc.LabelParameterVersion(ctx, &ascTypes.LabelParameterVersionInput{
		Name:    paramName,
		Version: version,
		Labels:  labels,
	})
	if err != nil {
		return fmt.Errorf("label parameter version: %w", err)
	}

	if len(invalidLabels) > 0 {
		fmt.Printf("Warning: invalid labels: %v\n", invalidLabels)
	}

	// Report success for valid labels
	validCount := len(labels) - len(invalidLabels)
	if validCount > 0 {
		if version > 0 {
			fmt.Printf("Added %d label(s) to %s version %d\n", validCount, paramName, version)
		} else {
			fmt.Printf("Added %d label(s) to %s (latest version)\n", validCount, paramName)
		}
	}

	return nil
}

// parseParamVersion parses a parameter name with optional version suffix.
// e.g., "/myapp/config:5" -> ("/myapp/config", 5, nil)
// e.g., "/myapp/config" -> ("/myapp/config", 0, nil) (0 means latest)
func parseParamVersion(input string) (string, int64, error) {
	// Find the last colon that's after the last slash (to handle paths with colons)
	lastSlash := strings.LastIndex(input, "/")
	colonIdx := strings.LastIndex(input, ":")

	if colonIdx > lastSlash && colonIdx < len(input)-1 {
		paramName := input[:colonIdx]
		versionStr := input[colonIdx+1:]

		version, err := strconv.ParseInt(versionStr, 10, 64)
		if err != nil {
			return "", 0, fmt.Errorf("invalid version number %q: must be a positive integer", versionStr)
		}
		if version < 1 {
			return "", 0, fmt.Errorf("invalid version number %d: must be a positive integer", version)
		}

		return paramName, version, nil
	}

	return input, 0, nil
}
