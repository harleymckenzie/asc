package ssm

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ssm"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	newRevertFlags(revertCmd)
}

var revertCmd = &cobra.Command{
	Use:   "revert <parameter-name> <version>",
	Short: "Revert a parameter to a previous version",
	Long: `Revert a parameter to a previous version.

This creates a new version of the parameter with the value from the specified
version. The old versions are preserved in history.

You can specify a version number or a label.`,
	GroupID: "actions",
	Args:    cobra.ExactArgs(2),
	Example: `  asc ssm revert /myapp/config 3
  asc ssm revert /myapp/config prod`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(RevertParameter(cmd, args))
	},
}

func newRevertFlags(cmd *cobra.Command) {
	// No flags needed
}

func RevertParameter(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	paramName := args[0]
	versionOrLabel := args[1]

	// Construct the source with version/label
	source := fmt.Sprintf("%s:%s", paramName, versionOrLabel)

	// Copy the old version to itself (creates new current version)
	err = svc.CopyParameter(ctx, &ascTypes.CopyParameterInput{
		Source:    source,
		Dest:      paramName,
		Overwrite: true,
	})
	if err != nil {
		return fmt.Errorf("revert parameter: %w", err)
	}

	fmt.Printf("Reverted %s to version %s\n", paramName, versionOrLabel)
	return nil
}
