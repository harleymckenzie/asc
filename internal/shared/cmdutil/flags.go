package cmdutil

import "github.com/spf13/cobra"

// GetPersistentFlags returns the profile and region from the command line flags.
func GetPersistentFlags(cmd *cobra.Command) (string, string) {
	profile, _ := cmd.Root().PersistentFlags().GetString("profile")
	region, _ := cmd.Root().PersistentFlags().GetString("region")
	return profile, region
}
