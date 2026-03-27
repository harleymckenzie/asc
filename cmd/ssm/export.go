package ssm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/harleymckenzie/asc/internal/service/ssm"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var (
	exportOutput string
	exportSplit  bool
)

func init() {
	newExportFlags(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:     "export <path>",
	Short:   "Export SSM parameters to a .env file",
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	Example: "  asc ssm export /myapp/prod/\n" +
		"  asc ssm export /myapp/prod/ --output prod.env\n" +
		"  asc ssm export /myapp/prod/ --split",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ExportSSMParameters(cmd, args))
	},
}

func newExportFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file (default: derived from SSM path)")
	cmd.Flags().BoolVar(&exportSplit, "split", false, "Write plain and secret parameters to separate files")
}

func ExportSSMParameters(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	ssmPath := args[0]

	params, err := svc.GetParametersByPath(ctx, &ascTypes.GetParametersByPathInput{
		Path:      ssmPath,
		Recursive: true,
		Decrypt:   true,
	})
	if err != nil {
		return fmt.Errorf("get parameters: %w", err)
	}

	if len(params) == 0 {
		fmt.Printf("No parameters found under: %s\n", ssmPath)
		return nil
	}

	if exportSplit {
		baseName := exportOutput
		if baseName == "" {
			baseName = ssmPathToBasename(ssmPath)
		} else {
			baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
		}
		plainFile := baseName + "_plain.env"
		secretsFile := baseName + "_secrets.env"

		plainF, err := os.Create(plainFile)
		if err != nil {
			return fmt.Errorf("create plain file: %w", err)
		}
		defer plainF.Close()

		secretsF, err := os.Create(secretsFile)
		if err != nil {
			return fmt.Errorf("create secrets file: %w", err)
		}
		defer secretsF.Close()

		plainCount, secretCount := 0, 0
		for _, param := range params {
			key := lastSegment(aws.ToString(param.Name))
			line := fmt.Sprintf("%s=%s\n", key, aws.ToString(param.Value))
			if string(param.Type) == "SecureString" {
				if _, err := fmt.Fprint(secretsF, line); err != nil {
					return fmt.Errorf("write secrets file: %w", err)
				}
				secretCount++
			} else {
				if _, err := fmt.Fprint(plainF, line); err != nil {
					return fmt.Errorf("write plain file: %w", err)
				}
				plainCount++
			}
		}

		if plainCount > 0 {
			fmt.Printf("Exported %d plain parameter(s) to %s\n", plainCount, plainFile)
		}
		if secretCount > 0 {
			fmt.Printf("Exported %d secret parameter(s) to %s\n", secretCount, secretsFile)
		}
	} else {
		outFile := exportOutput
		if outFile == "" {
			outFile = ssmPathToBasename(ssmPath) + ".env"
		}

		f, err := os.Create(outFile)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer f.Close()

		for _, param := range params {
			key := lastSegment(aws.ToString(param.Name))
			if _, err := fmt.Fprintf(f, "%s=%s\n", key, aws.ToString(param.Value)); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
		}

		fmt.Printf("Exported %d parameter(s) to %s\n", len(params), outFile)
	}

	return nil
}

// ssmPathToBasename returns the last non-empty segment of an SSM path.
func ssmPathToBasename(ssmPath string) string {
	parts := strings.Split(strings.TrimSuffix(ssmPath, "/"), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			return parts[i]
		}
	}
	return "parameters"
}

// lastSegment returns the last path component (after the final /).
func lastSegment(name string) string {
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		return name[idx+1:]
	}
	return name
}
