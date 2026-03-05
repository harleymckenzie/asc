package ssm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/harleymckenzie/asc/internal/service/ssm"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	newEditFlags(editCmd)
}

var editCmd = &cobra.Command{
	Use:     "edit <parameter-name>",
	Short:   "Edit an SSM parameter in your default editor",
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	Example: `  asc ssm edit /myapp/prod/config
  EDITOR=nano asc ssm edit /myapp/prod/config`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(EditSSMParameter(cmd, args))
	},
}

func newEditFlags(cmd *cobra.Command) {
	// No flags needed for now, but keeping the pattern consistent
}

func EditSSMParameter(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	paramName := args[0]

	// Get current parameter (with decryption for SecureString)
	param, err := svc.GetParameter(ctx, &ascTypes.GetParameterInput{
		Name:    paramName,
		Decrypt: true,
	})
	if err != nil {
		return fmt.Errorf("get parameter %s: %w", paramName, err)
	}

	originalValue := aws.ToString(param.Value)

	// Find editor
	editor := findEditor()
	if editor == "" {
		return fmt.Errorf("no editor found: set $EDITOR environment variable")
	}

	// Create temp file with current value
	tmpFile, err := os.CreateTemp("", "ssm-param-*.txt")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.WriteString(originalValue); err != nil {
		tmpFile.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	tmpFile.Close()

	// Open editor
	editorCmd := exec.Command(editor, tmpPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		return fmt.Errorf("run editor: %w", err)
	}

	// Read updated content
	newContent, err := os.ReadFile(tmpPath)
	if err != nil {
		return fmt.Errorf("read temp file: %w", err)
	}

	newValue := strings.TrimSuffix(string(newContent), "\n")

	// Check if value changed
	if newValue == originalValue {
		fmt.Println("No changes made.")
		return nil
	}

	// Update parameter
	err = svc.PutParameter(ctx, &ascTypes.PutParameterInput{
		Name:      paramName,
		Value:     newValue,
		Type:      string(param.Type),
		Overwrite: true,
	})
	if err != nil {
		return fmt.Errorf("update parameter: %w", err)
	}

	fmt.Printf("Updated: %s\n", paramName)
	return nil
}

// findEditor returns the editor to use, checking common environment variables
// and falling back to common editors.
func findEditor() string {
	// Check environment variables
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}

	// Try common editors
	editors := []string{"vim", "vi", "nano", "emacs"}
	for _, editor := range editors {
		if path, err := exec.LookPath(editor); err == nil {
			return path
		}
	}

	return ""
}
