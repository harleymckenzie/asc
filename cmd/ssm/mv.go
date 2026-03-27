package ssm

import (
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

// Variables for mv command
var (
	mvRecursive bool
	mvDryRun    bool
)

// Init function
func init() {
	newMvFlags(mvCmd)
}

var mvCmd = &cobra.Command{
	Use:     "mv <source> <destination>",
	Short:   "Move/rename SSM parameters",
	Aliases: []string{"move", "rename"},
	GroupID: "actions",
	Args:    cobra.ExactArgs(2),
	Example: "  asc ssm mv /myapp/old-key /myapp/new-key\n" +
		"  asc ssm mv /myapp/old-path/ /myapp/new-path/ --recursive",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(MoveSSMParameter(cmd, args[0], args[1]))
	},
}

// newMvFlags configures the flags for the mv command.
func newMvFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&mvRecursive, "recursive", "r", false, "Move all parameters under the source path.")
	cmd.Flags().BoolVarP(&mvDryRun, "dry-run", "n", false, "Show what would be moved without making changes.")
}

// MoveSSMParameter moves a parameter or parameters recursively.
func MoveSSMParameter(cmd *cobra.Command, source, dest string) error {
	ctx := context.TODO()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	// Warn if source looks like a path but --recursive wasn't specified
	if strings.HasSuffix(source, "/") && !mvRecursive {
		fmt.Fprintf(os.Stderr, "asc: -r not specified; omitting path '%s'\n", source)
		return nil
	}

	if mvRecursive {
		if mvDryRun {
			params, err := svc.GetParametersByPath(ctx, &ascTypes.GetParametersByPathInput{
				Path:      source,
				Recursive: true,
				Decrypt:   false,
			})
			if err != nil {
				return fmt.Errorf("get parameters by path: %w", err)
			}
			if len(params) == 0 {
				fmt.Printf("No parameters found under path: %s\n", source)
				return nil
			}
			fmt.Printf("Dry run: %d parameter(s) would be moved:\n", len(params))
			for _, p := range params {
				srcName := aws.ToString(p.Name)
				destName := strings.TrimSuffix(dest, "/") + strings.TrimPrefix(srcName, strings.TrimSuffix(source, "/"))
				fmt.Printf("  %s -> %s\n", srcName, destName)
			}
			return nil
		}
		// Recursive move
		count, err := svc.MoveParametersRecursive(ctx, source, dest)
		if err != nil {
			return fmt.Errorf("move parameters recursive: %w", err)
		}
		if count == 0 {
			fmt.Printf("No parameters found under path: %s\n", source)
		} else {
			fmt.Printf("Moved %d parameter(s) from %s to %s\n", count, source, dest)
		}
	} else {
		if mvDryRun {
			fmt.Printf("Would move %s to %s\n", source, dest)
			return nil
		}
		// Single parameter move
		err = svc.MoveParameter(ctx, &ascTypes.MoveParameterInput{
			Source: source,
			Dest:   dest,
		})
		if err != nil {
			return fmt.Errorf("move parameter: %w", err)
		}
		fmt.Printf("Moved %s to %s\n", source, dest)
	}

	return nil
}
