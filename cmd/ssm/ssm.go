package ssm

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/ssm"
	"github.com/spf13/cobra"
)

var (
	selectedColumns []string
	showValues      bool
)

func NewSSMCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssm",
		Short: "Perform SSM operations",
	}

	// ls sub command
	lsCmd := &cobra.Command{
		Use:   "ls",
		Short: "List parameters in Systems Manager Parameter Store",
		PreRun: func(cobraCmd *cobra.Command, args []string) {

			// Set default columns
			selectedColumns = []string{
				"type",
				"name",
			}
			if showValues {
				selectedColumns = append(selectedColumns, "value")
			}
		},
		Run: func(cobraCmd *cobra.Command, args []string) {
			ctx := context.TODO()
			profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")

			svc, err := ssm.NewSSMService(ctx, profile)
			if err != nil {
				log.Fatalf("Failed to initialize SSM service: %v", err)
			}

			if err := svc.ListParameters(ctx, selectedColumns, showValues); err != nil {
				log.Fatalf("Failed to list parameters: %v", err)
			}
		},
	}
	cmd.AddCommand(lsCmd)

	// Add flags - Output
	lsCmd.Flags().BoolVarP(&showValues, "output", "o", false, "Output the values of the parameters.")

	return cmd
}
