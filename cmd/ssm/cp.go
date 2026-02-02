package ssm

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/ssm"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// Variables
var (
	recursive bool
	overwrite bool
)

// Init function
func init() {
	newCpFlags(cpCmd)
}

var cpCmd = &cobra.Command{
	Use:     "cp <source> <destination>",
	Short:   "Copy SSM parameters",
	Aliases: []string{"copy"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(2),
	Example: "  asc ssm cp /myapp/prod/key /myapp/staging/key\n" +
		"  asc ssm cp /myapp/prod/ /myapp/staging/ --recursive\n" +
		"  asc ssm cp /source /dest --overwrite",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(CopySSMParameter(cmd, args[0], args[1]))
	},
}

// newCpFlags configures the flags for the cp command.
func newCpFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Copy all parameters under the source path.")
	cmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing parameters at destination.")
}

// CopySSMParameter copies a parameter or parameters recursively.
func CopySSMParameter(cmd *cobra.Command, source, dest string) error {
	ctx := context.TODO()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	if recursive {
		// Recursive copy
		count, err := svc.CopyParametersRecursive(ctx, source, dest, overwrite)
		if err != nil {
			return fmt.Errorf("copy parameters recursive: %w", err)
		}
		if count == 0 {
			fmt.Printf("No parameters found under path: %s\n", source)
		} else {
			fmt.Printf("Copied %d parameter(s) from %s to %s\n", count, source, dest)
		}
	} else {
		// Single parameter copy
		err = svc.CopyParameter(ctx, &ascTypes.CopyParameterInput{
			Source:    source,
			Dest:      dest,
			Overwrite: overwrite,
		})
		if err != nil {
			return fmt.Errorf("copy parameter: %w", err)
		}
		fmt.Printf("Copied %s to %s\n", source, dest)
	}

	return nil
}
