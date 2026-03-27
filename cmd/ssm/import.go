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

var (
	importType   string
	importYes    bool
	importDryRun bool
)

func init() {
	newImportFlags(importCmd)
}

var importCmd = &cobra.Command{
	Use:     "import <file> <ssm-path>",
	Short:   "Import parameters from a .env file into SSM",
	GroupID: "actions",
	Args:    cobra.ExactArgs(2),
	Example: "  asc ssm import prod.env /myapp/prod/\n" +
		"  asc ssm import secrets.env /myapp/prod/ --type SecureString\n" +
		"  asc ssm import vars.env /myapp/prod/ --yes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ImportSSMParameters(cmd, args))
	},
}

func newImportFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&importType, "type", "t", "String", "Parameter type: String or SecureString")
	cmd.Flags().BoolVarP(&importYes, "yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolVarP(&importDryRun, "dry-run", "n", false, "Show what would be imported without making changes")
}

func ImportSSMParameters(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()

	envFile := args[0]
	ssmPath := strings.TrimSuffix(args[1], "/")

	validTypes := map[string]bool{"String": true, "SecureString": true}
	if !validTypes[importType] {
		return fmt.Errorf("invalid type %q: must be String or SecureString", importType)
	}

	entries, err := parseEnvFile(envFile)
	if err != nil {
		return fmt.Errorf("parse env file: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No parameters found in file.")
		return nil
	}

	if importDryRun {
		fmt.Printf("Dry run: %d parameter(s) from %s to %s as %s\n\n",
			len(entries), envFile, ssmPath, importType)
	} else if !importYes {
		fmt.Printf("Import %d parameter(s) from %s to %s as %s. Continue? [y/N]: ",
			len(entries), envFile, ssmPath, importType)
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

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	imported, skipped := 0, 0
	for _, entry := range entries {
		paramName := ssmPath + "/" + entry[0]
		paramValue := entry[1]

		existing, getErr := svc.GetParameter(ctx, &ascTypes.GetParameterInput{
			Name:    paramName,
			Decrypt: true,
		})

		if getErr == nil && existing != nil {
			if string(existing.Type) != importType {
				fmt.Printf("Would skip %s: type mismatch (existing: %s, requested: %s)\n",
					paramName, existing.Type, importType)
				skipped++
				continue
			}
			if aws.ToString(existing.Value) == paramValue {
				if importDryRun {
					fmt.Printf("Would skip %s: value unchanged\n", paramName)
				} else {
					fmt.Printf("Skipping %s: value unchanged\n", paramName)
				}
				skipped++
				continue
			}
			if importDryRun {
				fmt.Printf("Would update: %s\n", paramName)
			} else {
				if err := svc.PutParameter(ctx, &ascTypes.PutParameterInput{
					Name:      paramName,
					Value:     paramValue,
					Type:      importType,
					Overwrite: true,
				}); err != nil {
					return fmt.Errorf("update %s: %w", paramName, err)
				}
				fmt.Printf("Updated: %s\n", paramName)
			}
		} else {
			if importDryRun {
				fmt.Printf("Would create: %s\n", paramName)
			} else {
				if err := svc.PutParameter(ctx, &ascTypes.PutParameterInput{
					Name:      paramName,
					Value:     paramValue,
					Type:      importType,
					Overwrite: false,
				}); err != nil {
					return fmt.Errorf("create %s: %w", paramName, err)
				}
				fmt.Printf("Created: %s\n", paramName)
			}
		}
		imported++
	}

	if importDryRun {
		fmt.Printf("\nDry run: %d would be imported, %d would be skipped.\n", imported, skipped)
	} else {
		fmt.Printf("\nDone: %d imported, %d skipped.\n", imported, skipped)
	}
	return nil
}

// parseEnvFile reads a .env file and returns [key, value] pairs.
// Blank lines and lines starting with # are ignored.
func parseEnvFile(path string) ([][2]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries [][2]string
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("line %d: missing '=' separator", lineNum)
		}
		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}
		entries = append(entries, [2]string{key, value})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}
