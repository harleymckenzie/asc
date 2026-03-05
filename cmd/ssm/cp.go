package ssm

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
		"  asc ssm cp /myapp/prod/key /myapp/staging/\n" +
		"  asc ssm cp /myapp/prod/ /myapp/staging/ -r\n" +
		"  asc ssm cp /myapp/config:1 /myapp/config --overwrite",
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
	ctx := cmd.Context()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	// Check if source looks like a directory (ends with "/")
	sourceIsDir := strings.HasSuffix(source, "/")

	// Warn if trying to copy a directory without --recursive
	if sourceIsDir && !recursive {
		fmt.Fprintf(os.Stderr, "asc: -r not specified; omitting path '%s'\n", source)
		return nil
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
		// Handle directory-style destination (trailing slash)
		if strings.HasSuffix(dest, "/") {
			basename := paramBasename(source)
			dest = dest + basename
		}

		// Check if destination exists
		shouldOverwrite := overwrite
		existingParam, err := svc.GetParameter(ctx, &ascTypes.GetParameterInput{
			Name:    dest,
			Decrypt: false,
		})
		if err == nil && existingParam != nil {
			if !overwrite {
				fmt.Printf("Parameter %s already exists (type: %s).\n", dest, existingParam.Type)
				fmt.Print("Overwrite? [y/N]: ")

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
			shouldOverwrite = true
		}

		// Single parameter copy
		err = svc.CopyParameter(ctx, &ascTypes.CopyParameterInput{
			Source:    source,
			Dest:      dest,
			Overwrite: shouldOverwrite,
		})
		if err != nil {
			return fmt.Errorf("copy parameter: %w", err)
		}
		fmt.Printf("Copied %s to %s\n", source, dest)
	}

	return nil
}

// paramBasename extracts the parameter name from a path, stripping any version suffix.
// e.g., "/myapp/prod/config:1" -> "config"
func paramBasename(path string) string {
	// Strip version suffix if present (e.g., ":1" or ":my-label")
	if idx := strings.LastIndex(path, ":"); idx != -1 {
		// Make sure it's a version suffix, not part of the path
		slashIdx := strings.LastIndex(path, "/")
		if idx > slashIdx {
			path = path[:idx]
		}
	}

	// Get basename (everything after last slash)
	if idx := strings.LastIndex(path, "/"); idx != -1 {
		return path[idx+1:]
	}
	return path
}
