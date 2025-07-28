package cmdutil

import "github.com/spf13/cobra"

var (
	Tags []string
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
