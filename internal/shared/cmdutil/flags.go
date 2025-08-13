package cmdutil

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var (
	Tags         []string
	ValidLayouts = []string{"horizontal", "vertical"}
)

// GetPersistentFlags returns the profile and region from the command line flags.
func GetPersistentFlags(cmd *cobra.Command) (string, string) {
	profile, _ := cmd.Root().PersistentFlags().GetString("profile")
	region, _ := cmd.Root().PersistentFlags().GetString("region")
	return profile, region
}

// AddTagFlag adds the --tag flag to the command for filtering resources by tags.
func AddTagFlag(cmd *cobra.Command) {
	cmd.Flags().StringSliceVar(&Tags, "tags", nil, "Filter resources by tags (key=value)")
	if err := cmd.RegisterFlagCompletionFunc("tags", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		panic(err)
	}
}

// AddShowFlags adds shared flags for the show command with a configurable default layout.
func AddShowFlags(cmd *cobra.Command, defaultLayout string) {
	cmd.Flags().StringP("output", "o", defaultLayout, fmt.Sprintf("Output format (%s)", strings.Join(ValidLayouts, ", ")))
}

// GetLayout returns the layout value from the command flags.
func GetLayout(cmd *cobra.Command) string {
	layout, _ := cmd.Flags().GetString("output")
	return layout
}

// ValidateFlagChoice validates that the value of a flag is one of the choices.
func ValidateFlagChoice(cmd *cobra.Command, flag string, choices []string) error {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		return fmt.Errorf("invalid %s: %w", flag, err)
	}
	if !slices.Contains(choices, value) {
		return fmt.Errorf("invalid choice for %s flag: %s. Valid options: %s", flag, value, strings.Join(choices, ", "))
	}
	return nil
}
