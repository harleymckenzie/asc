package ssm

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/harleymckenzie/asc/internal/service/ssm"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

// Variables
var (
	force       bool
	rmRecursive bool
	rmDryRun    bool
)

// Init function
func init() {
	newRmFlags(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:     "rm <parameter-name>",
	Short:   "Delete SSM parameters",
	Aliases: []string{"remove", "delete"},
	GroupID: "actions",
	Args:    cobra.MinimumNArgs(1),
	Example: "  asc ssm rm /myapp/prod/key\n" +
		"  asc ssm rm /myapp/prod/key --force\n" +
		"  asc ssm rm /myapp/test/ --recursive\n" +
		"  asc ssm rm /myapp/test/ --recursive --force",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(DeleteSSMParameter(cmd, args))
	},
}

// newRmFlags configures the flags for the rm command.
func newRmFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt.")
	cmd.Flags().BoolVarP(&rmRecursive, "recursive", "r", false, "Delete all parameters under the path.")
	cmd.Flags().BoolVarP(&rmDryRun, "dry-run", "n", false, "Show what would be deleted without making changes.")
}

// DeleteSSMParameter deletes parameters with optional confirmation.
func DeleteSSMParameter(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	var names []string

	for _, arg := range args {
		if containsGlob(arg) {
			resolved, err := resolveGlob(ctx, svc, arg)
			if err != nil {
				return fmt.Errorf("resolve glob %s: %w", arg, err)
			}
			names = append(names, resolved...)
		} else {
			// Always try path resolution first; fall back to literal only if
			// no children are found (handles both "/path/" and "/path" forms).
			params, err := svc.GetParametersByPath(ctx, &ascTypes.GetParametersByPathInput{
				Path:      arg,
				Recursive: true,
				Decrypt:   false,
			})
			if err != nil {
				return fmt.Errorf("get parameters by path %s: %w", arg, err)
			}
			if len(params) > 0 {
				for _, p := range params {
					names = append(names, aws.ToString(p.Name))
				}
			} else {
				// No children — treat as a literal parameter name
				names = append(names, arg)
			}
		}
	}

	if len(names) == 0 {
		fmt.Println("No parameters to delete.")
		return nil
	}

	// Dry run: list what would be deleted and exit
	if rmDryRun {
		fmt.Printf("Dry run: %d parameter(s) would be deleted:\n", len(names))
		for _, name := range names {
			fmt.Printf("  - %s\n", name)
		}
		return nil
	}

	// Confirmation prompt unless --force
	if !force {
		fmt.Println("The following parameters will be deleted:")
		for _, name := range names {
			fmt.Printf("  - %s\n", name)
		}
		fmt.Printf("\nDelete %d parameter(s)? [y/N]: ", len(names))

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read confirmation: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	// Delete parameters
	if len(names) == 1 {
		// Single delete
		err = svc.DeleteParameter(ctx, &ascTypes.DeleteParameterInput{
			Name: names[0],
		})
		if err != nil {
			return fmt.Errorf("delete parameter: %w", err)
		}
		fmt.Printf("Deleted: %s\n", names[0])
	} else {
		// Batch delete
		failed, err := svc.DeleteParameters(ctx, &ascTypes.DeleteParametersInput{
			Names: names,
		})
		if err != nil {
			return fmt.Errorf("delete parameters: %w", err)
		}

		deleted := len(names) - len(failed)
		fmt.Printf("Deleted %d parameter(s).\n", deleted)

		if len(failed) > 0 {
			fmt.Println("Failed to delete:")
			for _, name := range failed {
				fmt.Printf("  - %s\n", name)
			}
		}
	}

	return nil
}
